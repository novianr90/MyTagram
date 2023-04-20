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

	connString := os.Getenv("DATABASE_URL")

	db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database:", err)
	}

	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{}, models.File{})
}

func main() {

	PORT := os.Getenv("PORT")

	routers.StartServer(db).Run(":" + PORT)
}
