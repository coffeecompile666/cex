package model

import "gorm.io/gorm"

type Market struct {
	gorm.Model

	Name string `gorm:"type:varchar(255);not null;unique"`
	Code string `gorm:"type:varchar(255);not null;unique"`
}
