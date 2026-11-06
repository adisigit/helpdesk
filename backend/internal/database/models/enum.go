package models

type Role string

const (
	RoleRequester Role = "requester"
	RoleAgent     Role = "agent"
	RoleAdmin     Role = "admin"
)

func (r Role) Valid() bool {
	return r == RoleRequester || r == RoleAgent || r == RoleAdmin
}

func (r Role) IsStaff() bool { return r == RoleAgent || r == RoleAdmin }

type TicketStatus string

const (
	StatusOpen       TicketStatus = "open"
	StatusInProgress TicketStatus = "in_progress"
	StatusResolved   TicketStatus = "resolved"
	StatusClosed     TicketStatus = "closed"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)
