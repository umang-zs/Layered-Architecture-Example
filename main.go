package main

import (
	"log"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	PhoneNo string `json:"phone"`
	Address string `json:"address"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:secret123@tcp(127.0.0.1:3306)/customer")
	if err != nil {
		log.Println(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/posts", getCustomers).Methods(http.MethodGet)
	router.HandleFunc("/posts/{name}", getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/createPost", createPost).Methods(http.MethodPost)
	http.ListenAndServe(":8000", router)
}
