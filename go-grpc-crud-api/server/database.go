package server

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func Init() {
	DatabaseConnection()
}

// DB is gorm -> go ORM struct
var DB *gorm.DB
var err error

type Movie struct {
	ID        string `gorm:"primarykey"`
	Title     string
	Genre     string
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"autoCreateTime:false"`
}

func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "postgres"
	dbUser := "root"
	pwd := "12345678"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		pwd,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	_ = DB.AutoMigrate(Movie{})
	if err != nil {
		log.Fatal("Error on connecting to the db...", err)
	}
	fmt.Println("Database connection successful...")
}
