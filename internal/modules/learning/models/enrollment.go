package models

import "time"

// Enrollment ระดับคอร์ส
type Enrollment struct {
	ID              string     `gorm:"type:char(36);primaryKey"`
	UserID          string     `gorm:"type:char(36);index;not null"`
	CourseID        string     `gorm:"type:char(36);index;not null"`
	Status          string     `gorm:"type:enum('enrolled','in_progress','completed','dropped');default:'enrolled'"`
	StartedAt       *time.Time `gorm:""`
	CompletedAt     *time.Time `gorm:""`
	LastAccessedAt  *time.Time `gorm:""`
	ProgressPercent float64    `gorm:"type:decimal(5,2);default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `gorm:"index"`
}

func (Enrollment) TableName() string { return "enrollments" }
