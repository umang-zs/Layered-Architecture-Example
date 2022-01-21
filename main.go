package main

import (
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/umang01-hash/layered-architecture/driver"
	handler "github.com/umang01-hash/layered-architecture/handlers/customer"
	"github.com/umang01-hash/layered-architecture/store"
)

func main() {

	db := driver.ConnectToMySQL()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	stores := store.New(db)
	h := handler.New(stores)

	r := mux.NewRouter()

	r.HandleFunc("/customers/{id}", h.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/customers", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/customers/{id}", h.DeleteByID).Methods(http.MethodDelete)
	r.HandleFunc("/customers/{id}", h.UpdateByID).Methods(http.MethodPut)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatalln(srv.ListenAndServe())
}
