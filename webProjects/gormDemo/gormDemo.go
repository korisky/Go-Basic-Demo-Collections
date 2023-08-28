package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
)

// User define simple struct -> we can use gorm.Model for convenience
type User struct {
	gorm.Model // similar to extend

	Name    string
	Address string
}

// BeforeCreate can work as 'aop' that it would work before the user be created
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Before creation")
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if strings.Contains(u.Name, "Tom") {
		err = errors.New("Tom could not be updated")
	}
	return
}

func main() {

	// mysql set up for gorm
	dsn := "root:12345678@tcp(127.0.0.1:3306)/go_use?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// auto migrate (create tables & missing stuff, would not delete unused columns)
	_ = db.AutoMigrate(&User{})

	// CRUD
	user := User{
		Name:    "Tom",
		Address: "0x1234fa23",
	}

	// the BeforeCreate would be executed before the insertion
	result := db.Create(&user)
	fmt.Println(result.RowsAffected) // the result is the db execute resul

}
