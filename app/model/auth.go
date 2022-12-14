package model

import "gorm.io/gorm"

type Auth struct {
	gorm.Model
	Username string
	Password string
	Salt     string
	Fullname string
}

func (a *Auth) FindByUsername(dbConn *gorm.DB) {
	dbConn.Where("username = ?", a.Username).Find(&a)
}

func (a *Auth) Create(dbConn *gorm.DB) {
	dbConn.Create(a)
}
