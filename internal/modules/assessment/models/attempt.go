package models

import "time"

type Attempt struct {
	ID           string    `gorm:"type:char(36);primaryKey"`
	AssessmentID string    `gorm:"type:char(36);index;not null"`
	UserID       string    `gorm:"type:char(36);index;not null"`
	Status       string    `gorm:"type:enum('in_progress','submitted','expired');default:'in_progress'"`
	StartedAt    time.Time `gorm:"not null"`
	SubmittedAt  *time.Time
	TimeLimitS   *int // duplicate จาก assessment ตอนเริ่ม เพื่อกันเปลี่ยนทีหลัง
	ScoreRaw     *int
	ScorePercent *float64
	IsPassed     *bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}

func (Attempt) TableName() string { return "assessment_attempts" }
