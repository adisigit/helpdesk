package domain

import (
	"helpdesk/backend/internal/database/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	open   = models.StatusOpen
	inProg = models.StatusInProgress
	solved = models.StatusResolved
	closed = models.StatusClosed

	req   = models.RoleRequester
	agent = models.RoleAgent
	admin = models.RoleAdmin
)

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name     string
		from, to models.TicketStatus
		role     models.Role
		owner    bool
		want     error
	}{
		{"agent open->in_progress", open, inProg, agent, false, nil},
		{"admin open->in_progress", open, inProg, admin, false, nil},
		{"agent in_progress->resolved", inProg, solved, agent, false, nil},
		{"admin in_progress->resolved", inProg, solved, admin, false, nil},
		{"owner resolved->closed", solved, closed, req, true, nil},
		{"agent resolved->closed", solved, closed, agent, false, nil},
		{"admin resolved->closed", solved, closed, admin, false, nil},
		{"owner reopen resolved->open", solved, open, req, true, nil},
		{"owner cancel open->closed", open, closed, req, true, nil},
		{"admin open->closed", open, closed, admin, false, nil},

		{"requester cannot open->in_progress", open, inProg, req, true, ErrTransitionForbidden},
		{"requester cannot in_progress->resolved", inProg, solved, req, true, ErrTransitionForbidden},
		{"agent cannot reopen", solved, open, agent, false, ErrTransitionForbidden},
		{"admin cannot reopen", solved, open, admin, false, ErrTransitionForbidden},
		{"agent cannot cancel open->closed", open, closed, agent, false, ErrTransitionForbidden},
		{"non-owner requester cannot close resolved", solved, closed, req, false, ErrTransitionForbidden},
		{"non-owner requester cannot cancel", open, closed, req, false, ErrTransitionForbidden},

		{"open->resolved skips step", open, solved, admin, false, ErrInvalidTransition},
		{"open->open same status", open, open, admin, false, ErrInvalidTransition},
		{"in_progress->open", inProg, open, admin, false, ErrInvalidTransition},
		{"in_progress->closed", inProg, closed, admin, false, ErrInvalidTransition},
		{"resolved->in_progress", solved, inProg, agent, false, ErrInvalidTransition},
		{"closed->open", closed, open, admin, false, ErrInvalidTransition},
		{"closed->in_progress", closed, inProg, admin, false, ErrInvalidTransition},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CanTransition(tc.from, tc.to, tc.role, tc.owner)
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestCanTransition_OnlyDocumentedCombinationsAllowed(t *testing.T) {
	allowed := 0
	for _, from := range allStatuses {
		for _, to := range allStatuses {
			for _, role := range []models.Role{req, agent, admin} {
				if CanTransition(from, to, role, true) == nil {
					allowed++
				}
			}
		}
	}
	assert.Equal(t, 10, allowed)
}

func TestAllowedTargets(t *testing.T) {
	assert.Equal(t, []models.TicketStatus{inProg, closed}, AllowedTargets(open, admin, false))
	assert.Equal(t, []models.TicketStatus{open, closed}, AllowedTargets(solved, req, true))
	assert.Empty(t, AllowedTargets(closed, admin, true))
}
