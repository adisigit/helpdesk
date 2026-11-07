package domain

import (
	"errors"
	"helpdesk/backend/internal/database/models"
	"slices"
)

var (
	ErrInvalidTransition   = errors.New("status transition is not allowed")
	ErrTransitionForbidden = errors.New("your role is not allowed to perform this transition")
)

type transition struct{ from, to models.TicketStatus }

var rules = map[transition][]models.Role{
	{models.StatusOpen, models.StatusInProgress}:     {models.RoleAgent, models.RoleAdmin},
	{models.StatusInProgress, models.StatusResolved}: {models.RoleAgent, models.RoleAdmin},
	{models.StatusResolved, models.StatusClosed}:     {models.RoleRequester, models.RoleAgent, models.RoleAdmin},
	{models.StatusResolved, models.StatusOpen}:       {models.RoleRequester},
	{models.StatusOpen, models.StatusClosed}:         {models.RoleRequester, models.RoleAdmin},
}

var allStatuses = []models.TicketStatus{
	models.StatusOpen, models.StatusInProgress, models.StatusResolved, models.StatusClosed,
}

func CanTransition(from, to models.TicketStatus, role models.Role, isOwner bool) error {
	roles, ok := rules[transition{from, to}]
	if !ok {
		return ErrInvalidTransition
	}
	if !slices.Contains(roles, role) {
		return ErrTransitionForbidden
	}
	if role == models.RoleRequester && !isOwner {
		return ErrTransitionForbidden
	}
	return nil
}

func AllowedTargets(from models.TicketStatus, role models.Role, isOwner bool) []models.TicketStatus {
	var out []models.TicketStatus
	for _, to := range allStatuses {
		if CanTransition(from, to, role, isOwner) == nil {
			out = append(out, to)
		}
	}
	return out
}
