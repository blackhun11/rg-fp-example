package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPG() (*gorm.DB, error) {
	var dsn = "host=localhost user=postgres password=postgres dbname=latihan port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	// var dsn = "host=docker.for.mac.localhost user=postgres password=postgres dbname=latihan port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
