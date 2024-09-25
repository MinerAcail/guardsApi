package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HomeArrival represents a log when a student arrives home.
type HomeArrival struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	StudentID string    `json:"student_id"`  // String reference to Student's UUID
	ParentID  string    `json:"parent_id"`   // String reference to Parent's UUID
	Confirmed bool      `gorm:"default:false" json:"confirmed"` // Confirmation by parent
	CreatedAt time.Time `json:"created_at"`                    // Auto-filled on creation, indicating arrival time
	UpdatedAt time.Time `json:"updated_at"`                    // Auto-updated on modification
}

type SchoolArrival struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	StudentID string    `json:"student_id"`  // String reference to Student's UUID
	StaffID  string    `json:"staff_id"`   // String reference to Parent's UUID
	Confirmed bool      `gorm:"default:false" json:"confirmed"` // Confirmation by parent
	CreatedAt time.Time `json:"created_at"`                    // Auto-filled on creation, indicating arrival time
	UpdatedAt time.Time `json:"updated_at"`                    // Auto-updated on modification
}

// BeforeCreate hook to generate a UUID before creating a new HomeArrival log.
func (h *HomeArrival) BeforeCreate(tx *gorm.DB) (err error) {
	h.ID = uuid.New()
	h.CreatedAt = time.Now() // Set CreatedAt to the current time upon creation
	return
}
