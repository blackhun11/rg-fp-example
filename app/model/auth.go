package model

import "gorm.io/gorm"

type Auth struct {
	gorm.Model
	Username string
	Password string
	Salt     string
	Fullname string
}
