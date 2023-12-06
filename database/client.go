package database

import (
	_ "database/sql"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/diiineeei/ApiLojinha/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var InstancePostgres *gorm.DB
var dbError error

func ConnectPostgres(connectionString string) {

	user := "postgres"
	password := "l;+t#N8@I:8DC#h5"
	dbHost := "testes-rodinei:us-central1:lojinha"
	databaseName := "lojinha"

	connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, user, databaseName, password)

	InstancePostgres, dbError = gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        connString,
	}))

	//InstancePostgres, dbError = gorm.Open(postgres.Open("postgres@cloudsql(testes-rodinei:us-central1:lojinha)/lojinha?charset=utf8&parseTime=True&loc=UTC"), &gorm.Config{})

	//InstancePostgres, dbError = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	//if dbError != nil {
	//	log.Fatal(dbError)
	//	panic("Cannot connect to DB")
	//}
	log.Println("Connected to Database Postgres!")
}

func Migrate() {
	err := InstancePostgres.AutoMigrate(&models.Produto{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database Migration Completed!")
}
