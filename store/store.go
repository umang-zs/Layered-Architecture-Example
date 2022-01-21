package store

import (
	"database/sql"
	"strings"
	"log"
	"strconv"

	"github.com/umang01-hash/layered-architecture/model"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) store {
	return store{db: db}
}

func (s store) Get(id int) (model.Customer, error) {
	var c model.Customer

	err := s.db.QueryRow("SELECT * FROM customers WHERE ID = ?", id).
		Scan(&c.ID, &c.Name, &c.Address, &c.PhoneNo)


	return c, err
}

func (s store) Create(c model.Customer) (model.Customer, error) {
	_, err := s.db.Exec("INSERT INTO customers (ID,Name,phoneNo,Address) VALUES( ?, ?, ?, ? )", &c.ID, &c.Name, &c.PhoneNo, &c.Address)

	if err != nil {
		return model.Customer{}, err
	}


	return s.Get((c.ID))
}

func (s store) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM customers WHERE id = ?;", id)

	return err
}

func (s store) Update(c model.Customer) error {
	id := (strconv.Itoa(c.ID))
	query := createQuery(id, c)
	log.Println(query)
	if query == "" {
		return nil
	}

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func createQuery(id string, c model.Customer) string {
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

