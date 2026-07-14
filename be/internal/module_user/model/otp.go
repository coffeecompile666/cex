package model

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model

	HashCode     string     `gorm:"type:varchar(255);not null"`
	FailAttempts int        `gorm:"type:int;not null;default:0"`
	ExpiredAt    time.Time  `gorm:"type:timestamp;not null"`
	RevokedAt    *time.Time `gorm:"type:timestamp"` // nullable: nil = not revoked
	UserID       uint       `gorm:"index:idx_otp_user_id"`
	User         User       `gorm:"foreignkey:UserID"`
}
