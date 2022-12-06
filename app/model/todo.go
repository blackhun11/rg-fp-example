package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Desc   string `json:"desc,omitempty"`
	Status bool   `json:"status,omitempty"`
}

func (t Todo) Get(dbConn *gorm.DB) []Todo {
	todos := make([]Todo, 0)
	dbConn.Find(&todos)
	return todos
}

func (t Todo) Insert(dbConn *gorm.DB, desc string) {
	dbConn.Create(&Todo{
		Desc: desc,
	})
}

func (t Todo) Update(dbConn *gorm.DB, id int, status bool) {
	dbConn.Model(&Todo{}).Where("id = ?", id).Update("status", status)
}

func (t Todo) Delete(dbConn *gorm.DB) {
	dbConn.Model(&Todo{}).Where("status = ?", true).Update("deleted_at", time.Now())
}
