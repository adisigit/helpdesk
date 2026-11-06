package models

import "github.com/google/uuid"

type Comment struct {
	Base
	TicketID   uuid.UUID `gorm:"type:uuid;not null;index" json:"ticket_id"`
	AuthorID   uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	Body       string    `gorm:"not null" json:"body"`
	IsInternal bool      `gorm:"not null;default:false" json:"is_internal"`

	Author User `gorm:"foreignKey:AuthorID" json:"author"`
}
