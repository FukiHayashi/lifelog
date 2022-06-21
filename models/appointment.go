package models

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Appointment struct {
	LifeLogId uuid.UUID      `json:"-" gorm:"not null"`
	ID        uuid.UUID      `json:"payload" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Start     string         `json:"start"`
	End       string         `json:"end"`
	Title     string         `json:"title"`
	Class     string         `json:"class"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
