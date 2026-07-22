package model

import "gorm.io/gorm"

type Market struct {
	gorm.Model

	Name           string `gorm:"type:varchar(255);not null;unique"`
	Symbol         string `gorm:"type:varchar(255);not null;unique"`
	Decimals       int32  `gorm:"type:int;not null"`
	Precision      int32  `gorm:"type:int;not null"`
	SmallestUnit   string `gorm:"type:varchar(255);not null"`
	IsBaseCurrency bool   `gorm:"type:bool;not null"`
}
