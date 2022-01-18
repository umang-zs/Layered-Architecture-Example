package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
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

	defer func() {
		err := rows.Err()
		if err != nil {
			log.Println(err)
			return
		}
	}()

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

func createCustomer(w http.ResponseWriter, r *http.Request) {
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

	_, err = feedData.Exec(c.ID, c.Name, c.PhoneNo, c.Address)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Fprintf(w, "New customer created!")

}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err = db.Exec("DELETE FROM customers WHERE ID = ?", id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Fprintf(w, "Customer with ID = %s was deleted", id)

}

//func updateCustomer(w http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//	id := params["id"]
//	feedData, err := db.Prepare("UPDATE customers SET Name=? , phoneNo=? , Address=? where id=?")
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//	defer feedData.Close()
//
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		log.Println(err.Error())
//
//	}
//
//	var c customer
//	err = json.Unmarshal(body, &c)
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//
//	_, err = feedData.Exec(c.Name, c.PhoneNo, c.Address, id)
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//
//	fmt.Fprintf(w, "Customer with ID = %s Updated!", id)
//
//}

func createQuery(id string, c customer) string {
	query := "UPDATE customers SET "
	var q []string
	if c.Name != "" {
		q = append(q, " name = \""+c.Name+"\"")
	}
	if c.PhoneNo != "" {
		q = append(q, " phoneNo = \""+c.PhoneNo+"\"")
	}
	if c.Address != "" {
		q = append(q, " address = \""+c.Address+"\"")
	}

	if q == nil {
		return ""
	}

	query += strings.Join(q, " , ")

	query += " where ID = " + string(id) + " ; "

	return query

}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var c customer
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := createQuery(id, c)
	if query == "" {
		return
	}

	_, err = db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Customer with ID = %s is updated!", id)

}
