package models

import "time"

// ความคืบหน้าต่อบทเรียน
type UserLessonProgress struct {
	ID              string  `gorm:"type:char(36);primaryKey"`
	UserID          string  `gorm:"type:char(36);index;not null"`
	LessonID        string  `gorm:"type:char(36);index;not null"`
	ProgressPercent float64 `gorm:"type:decimal(5,2);default:0"`
	CurrentPosition int64   `gorm:"default:0"` // วินาทีวิดีโอ / เลขหน้า
	MaxPosition     int64   `gorm:"default:0"` // ความยาววิดีโอ / จำนวนหน้า
	StartedAt       *time.Time
	CompletedAt     *time.Time
	IsUnlocked      bool `gorm:"default:false"` // cache เร็ว
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `gorm:"index"`
}

func (UserLessonProgress) TableName() string { return "user_lesson_progress" }
