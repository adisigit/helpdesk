package service

import (
	"context"
	"helpdesk/backend/internal/database/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeComments struct {
	items           []models.Comment
	lastIncludeFlag bool
}

func (f *fakeComments) Create(_ context.Context, c *models.Comment) error {
	c.ID = uuid.New()
	f.items = append(f.items, *c)
	return nil
}

func (f *fakeComments) ListByTicket(_ context.Context, id uuid.UUID, includeInternal bool) ([]models.Comment, error) {
	f.lastIncludeFlag = includeInternal
	var out []models.Comment
	for _, c := range f.items {
		if c.TicketID == id && (includeInternal || !c.IsInternal) {
			out = append(out, c)
		}
	}
	return out, nil
}

func newCommentEnv(t *testing.T) (*ticketEnv, *CommentService, *fakeComments) {
	e := newTicketEnv(t)
	fc := &fakeComments{}
	return e, NewCommentService(e.tickets, fc), fc
}

func TestAddComment_RequesterCannotPostInternalNote(t *testing.T) {
	e, cs, _ := newCommentEnv(t)
	tk := e.newTicket(t)
	ctx := context.Background()

	_, err := cs.Add(ctx, e.owner, tk.ID, "catatan", true)
	assert.ErrorIs(t, err, ErrForbidden)

	c, err := cs.Add(ctx, e.agent, tk.ID, "catatan", true)
	require.NoError(t, err)
	assert.True(t, c.IsInternal)
}

func TestAddComment_OtherRequesterGetsNotFound(t *testing.T) {
	e, cs, _ := newCommentEnv(t)
	tk := e.newTicket(t)
	_, err := cs.Add(context.Background(), e.other, tk.ID, "halo", false)
	assert.ErrorIs(t, err, ErrTicketNotFound)
}

func TestAddComment_BlankBodyRejected(t *testing.T) {
	e, cs, _ := newCommentEnv(t)
	tk := e.newTicket(t)
	_, err := cs.Add(context.Background(), e.owner, tk.ID, "   ", false)
	assert.ErrorIs(t, err, ErrEmptyComment)
}

func TestListFor_RequesterNeverSeesInternalNotes(t *testing.T) {
	e, cs, fc := newCommentEnv(t)
	tk := e.newTicket(t)
	ctx := context.Background()

	_, _ = cs.Add(ctx, e.agent, tk.ID, "internal", true)
	_, _ = cs.Add(ctx, e.agent, tk.ID, "publik", false)
	_, _ = cs.Add(ctx, e.owner, tk.ID, "terima kasih", false)

	got, err := cs.ListFor(ctx, e.owner, tk.ID)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.False(t, fc.lastIncludeFlag, "requester harus meminta tanpa catatan internal")
	for _, c := range got {
		assert.False(t, c.IsInternal)
	}

	got, err = cs.ListFor(ctx, e.agent, tk.ID)
	require.NoError(t, err)
	assert.Len(t, got, 3)
}
