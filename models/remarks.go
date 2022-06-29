package models

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Remarks struct {
	LifeLogId uuid.UUID      `json:"-" gorm:"not null"`
	ID        uuid.UUID      `json:"payload" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title     *string        `json:"title" gorm:"not null"`
	Date      *string        `json:"date" gorm:"not null"`
	Class     string         `json:"class"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
