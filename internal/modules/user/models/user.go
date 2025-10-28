package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser              Role = "user"
	RoleAdmin             Role = "admin"
	RolePlanning          Role = "planning"
	RolePordcutionControl Role = "production_control"
	RoleProduction        Role = "production"
	RoleOutsourcing       Role = "outsourcing"
	RoleProductionADH     Role = "production_adh"
	RoleProductionWeld    Role = "production_weld"
	RoleProductionAssy    Role = "production_assy"
)

func AllRoles() []Role {
	return []Role{
		RoleUser, RoleAdmin, RolePlanning, RolePordcutionControl, RoleProduction, RoleOutsourcing, RoleProductionADH, RoleProductionWeld, RoleProductionAssy,
	}
}

type User struct {
	ID           string         `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeCode string         `gorm:"size:30;uniqueIndex;not null" json:"employee_code"`
	Email        string         `gorm:"size:190;uniqueIndex;not null" json:"email"`
	FirstName    string         `gorm:"size:100;not null" json:"first_name"`
	LastName     string         `gorm:"size:100;not null" json:"last_name"`
	Role         Role           `gorm:"type:varchar(32);not null;default:'user'" json:"role"`
	Phone        *string        `gorm:"size:50" json:"phone"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	IsActive     bool           `gorm:"type:tinyint(1);default:1" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}
