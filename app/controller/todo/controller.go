package todo

import (
	"wb/app/model"

	"gorm.io/gorm"
)

type Contract interface {
	Get() []model.Todo
	Insert(desc string)
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

func (t Todo) Get() []model.Todo {
	todo := model.Todo{}
	return todo.Get(t.dbConn)
}

func (t Todo) Insert(desc string) {
	todo := model.Todo{}
	todo.Insert(t.dbConn, desc)
}

func (t Todo) Update(id int, status bool) {
	todo := model.Todo{}
	todo.Update(t.dbConn, id, status)
}

func (t Todo) Delete() {
	todo := model.Todo{}
	todo.Delete(t.dbConn)
}
