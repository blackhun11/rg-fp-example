package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Desc   string `json:"desc,omitempty"`
	Status bool   `json:"status,omitempty"`
}
