package models

import "gorm.io/gorm"

type AdminModel struct {
	gorm.Model
	Name     string `gorm:"size(100)"`
	Email    string `gorm:"unique"`
	Password string `gorm:"not null"`
}
