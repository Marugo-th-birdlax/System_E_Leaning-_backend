package models

import "time"

type Choice struct {
	ID         string `gorm:"type:char(36);primaryKey"`
	QuestionID string `gorm:"type:char(36);index;not null"`
	Label      string `gorm:"type:text;not null"`
	IsCorrect  bool   `gorm:"not null;default:0"`
	Seq        int    `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}

func (Choice) TableName() string { return "assessment_choices" }
