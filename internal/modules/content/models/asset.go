package models

import (
	"time"
)

type Asset struct {
	ID        string  `gorm:"type:char(36);primaryKey"`
	OwnerType *string `gorm:"size:20"` // "lesson" / nil
	OwnerID   *string `gorm:"type:char(36)"`
	Kind      string  `gorm:"size:20;not null"` // "video","slide","image","doc"
	Filename  string  `gorm:"size:255;not null"`
	MimeType  string  `gorm:"size:100;not null"`
	SizeBytes int64   `gorm:"not null"`
	Storage   string  `gorm:"size:20;not null"` // "local","s3"
	URL       string  `gorm:"size:500;not null"`
	Checksum  *string `gorm:"size:80"`
	DurationS *int64  // วินาที (ถ้ารู้—optional)
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (Asset) TableName() string { return "assets" }
