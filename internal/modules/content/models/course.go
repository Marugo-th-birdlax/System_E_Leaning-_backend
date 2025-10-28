package models

import "time"

type Course struct {
	ID               string `gorm:"type:char(36);primaryKey"`
	Code             string `gorm:"size:50;uniqueIndex;not null"`
	Title            string `gorm:"size:255;not null"`
	Description      *string
	IsActive         bool `gorm:"not null;default:1"`
	EstimatedMinutes *int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `gorm:"index"`
}

func (Course) TableName() string { return "courses" }
