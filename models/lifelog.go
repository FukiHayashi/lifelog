package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LifeLog struct {
	UserId       uuid.UUID      `json:"-" gorm:"not null"`
	ID           uuid.UUID      `json:"-" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         *string        `json:"name" gorm:"not null"`
	LoggedAt     time.Time      `json:"-" gorm:"index"`
	Appointments []Appointment  `json:"appointments"`
	Remarks      Remarks        `json:"remarks"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
