package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Department struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	IsActive  bool           `gorm:"type:tinyint(1);default:1" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Department) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return nil
}
