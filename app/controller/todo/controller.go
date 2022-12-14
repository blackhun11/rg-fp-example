package todo

import (
	"wb/app/model"

	"gorm.io/gorm"
)

type Contract interface {
	Get(userID int64) []model.Todo
	Insert(desc string, userID int64)
	Update(id int, status bool)
	Delete()
}

type Todo struct {
	dbConn *gorm.DB
}

func NewTodo(dbConn *gorm.DB) Contract {
	return &Todo{
		dbConn: dbConn,
	}
}

func (t Todo) Get(userID int64) []model.Todo {
	todo := model.Todo{
		UserID: userID,
	}
	return todo.Get(t.dbConn)
}

func (t Todo) Insert(desc string, userID int64) {
	todo := model.Todo{
		Desc:   desc,
		UserID: userID,
	}
	todo.Insert(t.dbConn)
}

func (t Todo) Update(id int, status bool) {
	todo := model.Todo{
		Model: gorm.Model{
			ID: uint(id),
		},
		Status: status,
	}
	todo.Update(t.dbConn)
}

func (t Todo) Delete() {
	todo := model.Todo{}
	todo.Delete(t.dbConn)
}
