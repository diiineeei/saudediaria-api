package database

import (
	"github.com/diiineeei/ApiLojinha/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var InstancePostgres *gorm.DB
var dbError error

func ConnectPostgres(connectionString string) {
	InstancePostgres, dbError = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database Postgres!")
}

func Migrate() {
	err := InstancePostgres.AutoMigrate(&models.Produto{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database Migration Completed!")
}
