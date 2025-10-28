package models

import "time"

type CourseModule struct {
	ID          string `gorm:"type:char(36);primaryKey"`
	CourseID    string `gorm:"type:char(36);index;not null"`
	Title       string `gorm:"size:255;not null"`
	Description *string
	Seq         int  `gorm:"not null"`
	IsMandatory bool `gorm:"not null;default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (CourseModule) TableName() string { return "course_modules" }
