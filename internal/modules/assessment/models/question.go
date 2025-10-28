package models

import "time"

type Question struct {
	ID           string `gorm:"type:char(36);primaryKey"`
	AssessmentID string `gorm:"type:char(36);index;not null"`
	Type         string `gorm:"type:enum('single_choice','multiple_choice','true_false','short_text');not null"`
	Stem         string `gorm:"type:text;not null"`
	Explanation  *string
	Points       int `gorm:"not null;default:1"`
	Seq          int `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}

func (Question) TableName() string { return "assessment_questions" }
