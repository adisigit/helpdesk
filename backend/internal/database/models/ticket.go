package models

import "github.com/google/uuid"

type Ticket struct {
	Base
	Title       string       `gorm:"not null" json:"title"`
	Description string       `gorm:"not null" json:"description"`
	Category    string       `gorm:"not null" json:"category"`
	Priority    Priority     `gorm:"not null;default:medium" json:"priority"`
	Status      TicketStatus `gorm:"not null;default:open" json:"status"`
	RequesterID uuid.UUID    `gorm:"type:uuid;not null" json:"requester_id"`
	AssigneeID  *uuid.UUID   `gorm:"type:uuid" json:"assignee_id"`

	Requester User  `gorm:"foreignKey:RequesterID" json:"requester"`
	Assignee  *User `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
}
