package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Parent represents a Parent in the database
type Parent struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName       string         `json:"firstName"`
	LastName        string         `json:"lastName"`
	Email           string         `json:"email"`
	PhoneNumber     string         `json:"phoneNumber"`
	DateOfBirth     string         `json:"dateOfBirth"`
	Address         string         `json:"address"`
	Gender          *string        `json:"gender,omitempty"` // Optional field
	Position        string    `json:"position"`                // Can be teacher, admin, or maintenance
	Password        string    	     `json:"password"`
	Supervise       *pq.StringArray `gorm:"type:text[];column:supervise" json:"supervise"` // Array of children's IDs
	CreatedAt       time.Time      `json:"createdAt"`         // Auto-filled on creation
	UpdatedAt       time.Time      `json:"updatedAt"`         // Auto-updated on modification
}

// BeforeCreate hook to generate a UUID before creating a new Parent
func (p *Parent) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
