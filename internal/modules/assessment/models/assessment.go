package models

import "time"

type Assessment struct {
	ID          string `gorm:"type:char(36);primaryKey"`
	OwnerType   string `gorm:"size:20;not null"` // "course","module","lesson"
	OwnerID     string `gorm:"type:char(36);not null"`
	Type        string `gorm:"type:enum('pre','post','quiz');not null"`
	Title       string `gorm:"size:255;not null"`
	PassScore   int    `gorm:"not null;default:80"`
	MaxAttempts *int   `gorm:""`
	TimeLimitS  *int   `gorm:""`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (Assessment) TableName() string { return "assessments" }
