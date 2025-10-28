package models

import "time"

type Lesson struct {
	ID           string  `gorm:"type:char(36);primaryKey"`
	ModuleID     string  `gorm:"type:char(36);index;not null"`
	Title        string  `gorm:"size:255;not null"`
	ContentType  string  `gorm:"type:enum('slide','video','document','quiz');not null"`
	Seq          int     `gorm:"not null"`
	IsMandatory  bool    `gorm:"not null;default:1"`
	AssetID      *string `gorm:"type:char(36)"`
	AssessmentID *string `gorm:"type:char(36)"`
	DurationS    *int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}

func (Lesson) TableName() string { return "lessons" }
