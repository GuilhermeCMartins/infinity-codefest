package db

import (
	"log"
	"myapp/models"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	lock     = &sync.Mutex{}
	instance *gorm.DB
)

func GetInstance() *gorm.DB {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			dbURL := "postgres://goponey:poney@localhost:5432/goponey_db"

			db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
			if err != nil {
				log.Fatalln(err)
			}

			db.AutoMigrate(&models.User{})
			db.AutoMigrate(&models.Transaction{})

			instance = db
		}
	}

	return instance
}

func Close() {
	if instance == nil {
		return
	}

	sqlDB, err := instance.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}
