package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var dsn = "root:root@tcp(127.0.0.1:8889)/gorm?parseTime=true"
var db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

type Users struct {
	ID        uint           `gorm:"primarykey"`
	Name      string         `gorm:"column:name"`
	Year      string         `gorm:"column:years"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func main() {
	http.HandleFunc("/create", CreateUser)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func CreateUser(res http.ResponseWriter, req *http.Request) {

	//AutoMigrate will create the "users" table
	err := db.AutoMigrate(&Users{})
	if err != nil {
		log.Fatal(err)
	}

	Users := Users{
		Name: "hamed",
		Year: "2023",
	}

	err = db.Create(&Users).Error
	if err != nil {
		log.Fatalln(err)
	}

	json.NewEncoder(res).Encode(Users)
	fmt.Println("Field added ", Users)

}
