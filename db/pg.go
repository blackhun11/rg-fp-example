package db

import (
	"wb/app/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB
var dsn = "host=localhost user=postgres password=postgres dbname=latihan port=5432 sslmode=disable TimeZone=Asia/Jakarta"

func Connect() error {
	var err error
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func Migrate() error {
	return DBConn.AutoMigrate(&model.Todo{})
}
