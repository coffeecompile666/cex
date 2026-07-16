package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model

	HashRefreshToken string     `gorm:"type:varchar(255);not null;index:idx_session_hash"`
	ExpiredAt        time.Time  `gorm:"type:timestamp;not null"`
	RevokedAt        *time.Time `gorm:"type:timestamp"` // nullable: nil = active
	UserID           uint       `gorm:"index:idx_session_user_id"`
	User             User       `gorm:"foreignkey:UserID"`
}
