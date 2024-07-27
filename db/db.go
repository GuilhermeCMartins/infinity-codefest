package db

import (
	"log"
	"os/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	// TODO: Add env variable for db url
	dbURL := "postgres://goponey:poney@localhost:5432/goponey_db"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&user.User{})

	return db
}

//TODO: Add Close function
