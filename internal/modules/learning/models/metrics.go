package models

import "time"

type LearningMetric struct {
	ID               string   `gorm:"type:char(36);primaryKey"`
	UserID           string   `gorm:"type:char(36);index;not null"`
	CourseID         string   `gorm:"type:char(36);index;not null"`
	AvgScore         float64  `gorm:"type:decimal(5,2);default:0"`
	LastScore        *float64 `gorm:"type:decimal(5,2)"`
	AttemptsCount    int      `gorm:"default:0"`
	PassCount        int      `gorm:"default:0"`
	TotalTimeSeconds int64    `gorm:"default:0"`
	CompletionStatus string   `gorm:"size:32;default:'not_enrolled'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (LearningMetric) TableName() string { return "learning_metrics" }

type CourseOutcome struct {
	CourseID          string  `gorm:"type:char(36);primaryKey"`
	TotalEnrollments  int     `gorm:"default:0"`
	TotalCompleted    int     `gorm:"default:0"`
	AvgScore          float64 `gorm:"type:decimal(5,2);default:0"`
	PassRate          float64 `gorm:"type:decimal(5,2);default:0"`
	MedianTimeSeconds int64   `gorm:"default:0"`
	UpdatedAt         time.Time
}

func (CourseOutcome) TableName() string { return "course_outcomes" }
