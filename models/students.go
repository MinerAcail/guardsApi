package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Student represents a student in the database
type Student struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phoneNumber"`
	DateOfBirth   string    `json:"dateOfBirth"`
	Address       string    `json:"address"`
	Gender        *string   `json:"gender"` // Optional field
	Grade         string    `json:"grade"`
	ParentContact string    `json:"parentContact"`
}

// BeforeCreate hook to generate a UUID before `c`reating a new student
func (student *Student) BeforeCreate(tx *gorm.DB) (err error) {
	student.ID = uuid.New()
	return
}
