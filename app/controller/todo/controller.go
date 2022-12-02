package todo

import (
	"time"
	"wb/app/model"

	"gorm.io/gorm"
)

type Controller struct {
	dbConn *gorm.DB
}

type Contract interface {
	Get() []model.Todo
	Insert(desc string)
	Update(id int, status bool)
	Delete()
}

func NewController(dbConn *gorm.DB) Contract {
	return &Controller{
		dbConn: dbConn,
	}
}

func (c *Controller) Get() []model.Todo {
	todos := make([]model.Todo, 0)
	c.dbConn.Order("status").Order("created_at desc").Find(&todos)
	return todos
}

func (c *Controller) Insert(desc string) {
	data := model.Todo{
		Desc: desc,
	}
	c.dbConn.Create(&data)
}

func (c *Controller) Update(id int, status bool) {
	c.dbConn.Model(&model.Todo{}).Where("id = ?", id).Update("status", status)
}

func (c *Controller) Delete() {
	c.dbConn.Model(&model.Todo{}).Where("status = ?", true).Update("deleted_at", time.Now())
}
