package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func getCustomers(w http.ResponseWriter, r *http.Request) {
	//telling the client that the response sent is in json format
	w.Header().Set("Content-Type", "application/json")

	var customers []customer

	rows, err := db.Query("select * from customers")
	if err != nil {
		log.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var c customer
		err := rows.Scan(&c.ID, &c.Name, &c.PhoneNo, &c.Address)
		if err != nil {
			log.Println(err)
		}
		customers = append(customers, c)
	}

	customerJSON, _ := json.Marshal(customers)
	w.Write([]byte(customerJSON))
}

func getCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var customers []customer
	params := mux.Vars(r)
	name := params["name"]

	rows, err := db.Query("SELECT ID, Name , phoneNo , Address FROM customers WHERE Name= ?", name)

	if err != nil {
		log.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var c customer
		err := rows.Scan(&c.ID, &c.Name, &c.PhoneNo, &c.Address)
		if err != nil {
			log.Println(err)
		}

		customers = append(customers, c)

	}

	customerJSON, _ := json.Marshal(customers)

	if string(customerJSON) == "null" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No customer with Name = %s found! ", name)

		return
	}

	w.Write(customerJSON)

}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	feedData, err := db.Prepare("INSERT INTO customers (ID,Name,phoneNo,Address) VALUES( ?, ?, ?, ? )")

	if err != nil {
		log.Println(err)
	}

	defer feedData.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var c customer
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(c)
	_, err = feedData.Exec(c.ID, c.Name, c.PhoneNo, c.Address)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Fprintf(w, "New customer created!")

}
