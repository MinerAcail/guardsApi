package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Staff represents a staff in the database
type Staff struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phoneNumber"`
	DateOfBirth     string    `json:"dateOfBirth"`
	Address         string    `json:"address"`
	Gender          *string   `json:"gender,omitempty"`        // Optional field
	Password        string    `json:"password"`
	Position        string    `json:"position"`                // Can be teacher, admin, or maintenance
	SuperviseGrade  string    `json:"superviseGrade"`           // The grade the staff supervises
	CreatedAt       time.Time `json:"createdAt"`                // Auto-filled on creation
	UpdatedAt       time.Time `json:"updatedAt"`                // Auto-updated on modification
}


// BeforeCreate hook to generate a UUID before `c`reating a new staff
func (staff *Staff) BeforeCreate(tx *gorm.DB) (err error) {
	staff.ID = uuid.New()
	return
}
