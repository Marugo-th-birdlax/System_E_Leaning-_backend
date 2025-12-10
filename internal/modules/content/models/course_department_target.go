package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// คอร์สนี้ target แผนกไหนบ้าง
type CourseDepartmentTarget struct {
	ID           string         `gorm:"type:char(36);primaryKey" json:"id"`
	CourseID     string         `gorm:"type:char(36);index;not null" json:"course_id"`
	DepartmentID string         `gorm:"type:char(36);index;not null" json:"department_id"`
	IsMandatory  bool           `gorm:"type:tinyint(1);default:0" json:"is_mandatory"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *CourseDepartmentTarget) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return nil
}
