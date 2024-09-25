package models

import (
	"time"

	"github.com/google/uuid"
)

type Calendar struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Height    int       `gorm:"type:int;not null" json:"height"`
	Day       string `json:"day"`
	EndDay    string `json:"endDay"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate a UUID before `c`reating a new calendar event