package service

import (
	"context"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/domain"
	"helpdesk/backend/internal/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeTickets struct {
	byID       map[uuid.UUID]*models.Ticket
	lastFilter repository.TicketFilter
}

func newFakeTickets() *fakeTickets {
	return &fakeTickets{byID: map[uuid.UUID]*models.Ticket{}}
}

func (f *fakeTickets) Create(_ context.Context, t *models.Ticket) error {
	t.ID = uuid.New()
	cp := *t
	f.byID[t.ID] = &cp
	return nil
}

func (f *fakeTickets) FindByID(_ context.Context, id uuid.UUID) (*models.Ticket, error) {
	if t, ok := f.byID[id]; ok {
		cp := *t
		return &cp, nil
	}
	return nil, repository.ErrNotFound
}

func (f *fakeTickets) List(_ context.Context, flt repository.TicketFilter) ([]models.Ticket, int64, error) {
	f.lastFilter = flt
	return nil, 0, nil
}

func (f *fakeTickets) Update(_ context.Context, t *models.Ticket, expected models.TicketStatus) error {
	cur, ok := f.byID[t.ID]
	if !ok {
		return repository.ErrNotFound
	}
	if cur.Status != expected {
		return repository.ErrConflict
	}
	cp := *t
	f.byID[t.ID] = &cp
	return nil
}

type ticketEnv struct {
	svc                        *TicketService
	tickets                    *fakeTickets
	owner, other, agent, admin Actor
}

func newTicketEnv(t *testing.T) *ticketEnv {
	t.Helper()
	users := newFakeUsers()
	mk := func(email string, role models.Role) Actor {
		u := &models.User{Name: email, Email: email, Role: role}
		require.NoError(t, users.Create(context.Background(), u))
		return Actor{ID: u.ID, Role: role}
	}
	tickets := newFakeTickets()
	return &ticketEnv{
		svc:     NewTicketService(tickets, users),
		tickets: tickets,
		owner:   mk("owner@x.io", models.RoleRequester),
		other:   mk("other@x.io", models.RoleRequester),
		agent:   mk("agent@x.io", models.RoleAgent),
		admin:   mk("admin@x.io", models.RoleAdmin),
	}
}

func (e *ticketEnv) newTicket(t *testing.T) *models.Ticket {
	t.Helper()
	tk, err := e.svc.Create(context.Background(), e.owner, CreateTicketInput{
		Title: "VPN down", Description: "cannot connect", Category: "network",
	})
	require.NoError(t, err)
	return tk
}

func TestCreate_ForcesOpenStatusAndRequester(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)
	assert.Equal(t, models.StatusOpen, tk.Status)
	assert.Equal(t, e.owner.ID, tk.RequesterID)
	assert.Equal(t, models.PriorityMedium, tk.Priority)
	assert.Nil(t, tk.AssigneeID)
}

func TestCreate_RejectsBlankTitle(t *testing.T) {
	e := newTicketEnv(t)
	_, err := e.svc.Create(context.Background(), e.owner, CreateTicketInput{
		Title: "   ", Description: "isi", Category: "other",
	})
	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestGet_RequesterCannotReadOthersTicket(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)
	ctx := context.Background()

	_, err := e.svc.Get(ctx, e.other, tk.ID)
	assert.ErrorIs(t, err, ErrTicketNotFound, "requester lain tidak boleh melihat")

	for name, a := range map[string]Actor{"owner": e.owner, "agent": e.agent, "admin": e.admin} {
		_, err := e.svc.Get(ctx, a, tk.ID)
		assert.NoError(t, err, name)
	}
}

func TestGet_UnknownIDIsNotFound(t *testing.T) {
	e := newTicketEnv(t)
	_, err := e.svc.Get(context.Background(), e.admin, uuid.New())
	assert.ErrorIs(t, err, ErrTicketNotFound)
}

func TestList_RequesterIsScopedToOwnTickets(t *testing.T) {
	e := newTicketEnv(t)
	ctx := context.Background()

	_, err := e.svc.List(ctx, e.owner, ListTicketsInput{Page: 1, PageSize: 20})
	require.NoError(t, err)
	require.NotNil(t, e.tickets.lastFilter.RequesterID)
	assert.Equal(t, e.owner.ID, *e.tickets.lastFilter.RequesterID)

	_, err = e.svc.List(ctx, e.agent, ListTicketsInput{Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Nil(t, e.tickets.lastFilter.RequesterID, "staff tidak dibatasi")
}

func TestList_AssigneeFilterAndPageClamp(t *testing.T) {
	e := newTicketEnv(t)
	ctx := context.Background()

	_, err := e.svc.List(ctx, e.agent, ListTicketsInput{Assignee: "me", Page: 0, PageSize: 9999})
	require.NoError(t, err)
	assert.Equal(t, e.agent.ID, *e.tickets.lastFilter.AssigneeID)
	assert.Equal(t, 1, e.tickets.lastFilter.Page)
	assert.Equal(t, 100, e.tickets.lastFilter.PageSize)

	_, err = e.svc.List(ctx, e.agent, ListTicketsInput{Assignee: "none"})
	require.NoError(t, err)
	assert.True(t, e.tickets.lastFilter.Unassigned)

	_, err = e.svc.List(ctx, e.agent, ListTicketsInput{Assignee: "bukan-uuid"})
	assert.ErrorIs(t, err, ErrInvalidFilter)
}

func ptr[T any](v T) *T { return &v }

func TestUpdate_InvalidTransitionIsRejectedAndNotPersisted(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)

	_, err := e.svc.Update(context.Background(), e.admin, tk.ID,
		UpdateTicketInput{Status: ptr(models.StatusResolved)})
	assert.ErrorIs(t, err, domain.ErrInvalidTransition)
	assert.Equal(t, models.StatusOpen, e.tickets.byID[tk.ID].Status)
}

func TestUpdate_RequesterCannotStartOrChangePriority(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)
	ctx := context.Background()

	_, err := e.svc.Update(ctx, e.owner, tk.ID, UpdateTicketInput{Status: ptr(models.StatusInProgress)})
	assert.ErrorIs(t, err, domain.ErrTransitionForbidden)

	_, err = e.svc.Update(ctx, e.owner, tk.ID, UpdateTicketInput{Priority: ptr(models.PriorityUrgent)})
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdate_OtherRequesterGetsNotFound(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)
	_, err := e.svc.Update(context.Background(), e.other, tk.ID,
		UpdateTicketInput{Status: ptr(models.StatusClosed)})
	assert.ErrorIs(t, err, ErrTicketNotFound)
}

func TestUpdate_FullLifecycleWithAutoAssign(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)
	ctx := context.Background()

	got, err := e.svc.Update(ctx, e.agent, tk.ID, UpdateTicketInput{Status: ptr(models.StatusInProgress)})
	require.NoError(t, err)
	require.NotNil(t, got.AssigneeID)
	assert.Equal(t, e.agent.ID, *got.AssigneeID, "agent otomatis jadi assignee")

	got, err = e.svc.Update(ctx, e.agent, tk.ID, UpdateTicketInput{Status: ptr(models.StatusResolved)})
	require.NoError(t, err)
	assert.Equal(t, models.StatusResolved, got.Status)

	got, err = e.svc.Update(ctx, e.owner, tk.ID, UpdateTicketInput{Status: ptr(models.StatusOpen)})
	require.NoError(t, err)
	assert.Equal(t, models.StatusOpen, got.Status)
}

func TestUpdate_AssigneeMustBeStaff(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)
	ctx := context.Background()

	_, err := e.svc.Update(ctx, e.admin, tk.ID, UpdateTicketInput{AssigneeID: &e.other.ID})
	assert.ErrorIs(t, err, ErrInvalidAssignee, "requester tidak boleh jadi assignee")

	_, err = e.svc.Update(ctx, e.admin, tk.ID, UpdateTicketInput{AssigneeID: ptr(uuid.New())})
	assert.ErrorIs(t, err, ErrInvalidAssignee, "user tidak ada")

	got, err := e.svc.Update(ctx, e.admin, tk.ID, UpdateTicketInput{AssigneeID: &e.agent.ID})
	require.NoError(t, err)
	assert.Equal(t, e.agent.ID, *got.AssigneeID)
}

func TestUpdate_AgentCanOnlyTakeTicketForSelf(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)
	ctx := context.Background()

	_, err := e.svc.Update(ctx, e.agent, tk.ID, UpdateTicketInput{AssigneeID: &e.admin.ID})
	assert.ErrorIs(t, err, ErrForbidden, "agent tidak boleh menugaskan orang lain")

	got, err := e.svc.Update(ctx, e.agent, tk.ID, UpdateTicketInput{AssigneeID: &e.agent.ID})
	require.NoError(t, err)
	assert.Equal(t, e.agent.ID, *got.AssigneeID)

	got, err = e.svc.Update(ctx, e.admin, tk.ID, UpdateTicketInput{AssigneeID: &e.admin.ID})
	require.NoError(t, err)
	assert.Equal(t, e.admin.ID, *got.AssigneeID)

	_, err = e.svc.Update(ctx, e.agent, tk.ID, UpdateTicketInput{AssigneeID: &e.agent.ID})
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdate_EmptyPatchIsRejected(t *testing.T) {
	e := newTicketEnv(t)
	tk := e.newTicket(t)
	_, err := e.svc.Update(context.Background(), e.agent, tk.ID, UpdateTicketInput{})
	assert.ErrorIs(t, err, ErrNothingToUpdate)
}
