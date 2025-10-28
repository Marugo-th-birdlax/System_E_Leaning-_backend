package models

import "time"

type Answer struct {
	ID         string `gorm:"type:char(36);primaryKey"`
	AttemptID  string `gorm:"type:char(36);index;not null"`
	QuestionID string `gorm:"type:char(36);index;not null"`
	// สำหรับตัวเลือก: เก็บ choice_ids (หลายค่าเมื่อ multiple_choice)
	SelectedChoiceIDs *string `gorm:"type:text"` // เก็บเป็น CSV/JSON (เลือกแบบใดแบบหนึ่ง) — เราจะใช้ CSV ง่ายๆ
	// สำหรับ short_text:
	TextAnswer *string
	IsCorrect  *bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}

func (Answer) TableName() string { return "assessment_answers" }
