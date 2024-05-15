package db

import (
	"database-example/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=super dbname=gorm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		print(err)
		return nil
	}
	err = db.AutoMigrate(

		&model.User{},
	)

	if err != nil {
		print(err)
		return nil
	}
	DB = db

	return DB
}
