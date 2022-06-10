package models

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Appointment struct {
	LifeLogId uuid.UUID      `json:"-" gorm:"not null"`
	ID        uuid.UUID      `json:"payload" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Start     string         `json:"start" grom:"default:00:00"`
	End       string         `json:"end" gorm:"default:01:00"`
	Title     string         `json:"title" gorm:"default:title"`
	Class     string         `json:"class" gorm:"default:default"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
