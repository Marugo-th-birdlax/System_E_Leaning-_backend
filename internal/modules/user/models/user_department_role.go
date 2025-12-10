package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ระดับ “หัวหน้าในแผนก”
type DepartmentRole string

const (
	DeptRoleManager DepartmentRole = "manager"
	DeptRoleLeader  DepartmentRole = "leader"
	DeptRoleMember  DepartmentRole = "member"
)

func AllDepartmentRoles() []DepartmentRole {
	return []DepartmentRole{DeptRoleManager, DeptRoleLeader, DeptRoleMember}
}

type UserDepartmentRole struct {
	ID           string         `gorm:"type:char(36);primaryKey" json:"id"`
	UserID       string         `gorm:"type:char(36);index;not null" json:"user_id"`
	DepartmentID string         `gorm:"type:char(36);index;not null" json:"department_id"`
	Role         DepartmentRole `gorm:"type:varchar(32);not null" json:"role"`

	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *UserDepartmentRole) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}
