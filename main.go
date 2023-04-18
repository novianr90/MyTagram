package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"final-project-hacktiv8/models"
	"final-project-hacktiv8/routers"
)

var (
	db  *gorm.DB
	err error
)

func init() {

	// var (
	// 	host     = os.Getenv("PGHOST")
	// 	port     = os.Getenv("PGPORT")
	// 	username = os.Getenv("PGDATABASE")
	// 	password = os.Getenv("PGPASSWORD")
	// 	dbName   = os.Getenv("PGUSER")
	// )

	connString := os.Getenv("DATABASE_URL")

	db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database:", err)
	}

	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func main() {

	PORT := os.Getenv("PGPORT")

	routers.StartServer(db).Run(":" + PORT)
}
