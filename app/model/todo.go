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

func (t *Todo) Get(dbConn *gorm.DB) []Todo {
	var todos = make([]Todo, 0)
	dbConn.Find(&todos)
	return todos
}

func (t *Todo) Insert(dbConn *gorm.DB) {
	dbConn.Create(&t)
}

func (t *Todo) Update(dbConn *gorm.DB) {
	dbConn.Model(&t).
		Where("id = ?", t.ID).
		Update("status", t.Status)
}

func (t *Todo) Delete(dbConn *gorm.DB) {
	dbConn.Model(&t).
		Where("status = ?", true).
		Update("deleted_at", time.Now())
}
