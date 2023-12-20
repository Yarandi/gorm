package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm/pkg/config"
	"log"
	"net/http"
	"time"
)

var db *gorm.DB
var err error

type Users struct {
	ID        uint           `gorm:"primarykey"`
	Name      string         `gorm:"column:name"`
	Year      string         `gorm:"column:years"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func main() {
	db, err = config.InitDB()
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/create", CreateUser)
	http.HandleFunc("/get-users", GetUsers)
	http.HandleFunc("/get-user", GetUser)
	http.HandleFunc("/delete-user", DeleteUser)
	http.HandleFunc("/update-user", UpdateUser)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// CreateUser will Create a user
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

// GetUsers will get all users
func GetUsers(res http.ResponseWriter, req *http.Request) {
	var Users []Users
	err = db.Find(&Users).Error
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(res).Encode(Users)
	fmt.Println("Get all rows")
}

// GetUser will get specific user
func GetUser(res http.ResponseWriter, req *http.Request) {

	//get specific id from urls
	userId := req.URL.Query().Get("id")

	if userId == "" {
		http.Error(res, "User id is required", http.StatusBadRequest) //400
		return
	}

	var user Users
	err := db.First(&user, userId).Error
	if err != nil {
		http.Error(res, "user not found", http.StatusNotFound) //401
		return
	}

	//return result
	json.NewEncoder(res).Encode(user)
	fmt.Println("user retrieved: ", user)
}

// DeleteUser will delete (soft delete a user)
func DeleteUser(res http.ResponseWriter, req *http.Request) {

	deletedId := req.URL.Query().Get("id")
	if deletedId == "" {
		http.Error(res, "Id is required for deletion", http.StatusBadRequest) //400
		return
	}
	err := db.Delete(&Users{}, deletedId).Error
	if err != nil {
		http.Error(res, "Error deleting user", http.StatusInternalServerError) //500
	}

	//return remain of users as a result, after delete on of them
	GetUsers(res, req)
	return
}

// UpdateUser will update a user
func UpdateUser(res http.ResponseWriter, req *http.Request) {
	updateID := req.URL.Query().Get("id")
	if updateID == "" {
		http.Error(res, "User id is required", http.StatusBadRequest)
		return
	}
	var user Users
	err = db.First(&user, updateID).Error
	if err != nil {
		http.Error(res, "User not Found", http.StatusNotFound) //404
		return
	}

	updatedUser := Users{
		Name: "hamed Yarandi",
		Year: "1987",
	}
	user.Name = updatedUser.Name
	user.Year = updatedUser.Year

	err = db.Save(&user).Error
	if err != nil {
		http.Error(res, "Error updating user", http.StatusInternalServerError) //500
		return
	}

	json.NewEncoder(res).Encode(user)
	fmt.Println("User Updated: ", user)
}
