package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Mail         string `gorm:"type:varchar(255);not null;index:idx_mail,unique"`
	HashPassword string `gorm:"type:varchar(255);not null"`
	IsVerified   bool   `gorm:"type:boolean;not null;default:false"`
}
