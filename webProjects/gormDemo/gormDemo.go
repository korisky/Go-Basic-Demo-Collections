package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {

	// mysql set up for gorm
	dsn := "root:12345678@tcp(127.0.0.1:3306)/epay_payment?charset=utf8mb4&parseTime=True&loc=Local"
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
}
