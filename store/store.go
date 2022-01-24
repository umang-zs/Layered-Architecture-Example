package store

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/umang01-hash/layered-architecture/model"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) store {
	return store{db: db}
}

func (s store) GetByID(id int) (model.Customer, error) {
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

	return s.GetByID((c.ID))
}

func (s store) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM customers WHERE id = ?;", id)

	return err
}

func (s store) Update(c model.Customer) error {
	id := (strconv.Itoa(c.ID))
	query, args := createQuery(id, c)

	if query == "" {
		return nil
	}

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func createQuery(id string, c model.Customer) (string, []interface{}) {
	var q []string
	var args []interface{}

	if c.Name != "" {
		q = append(q, " name=?")
		args = append(args, c.Name)
	}

	if c.Address != "" {
		q = append(q, " address=?")
		args = append(args, c.Address)
	}

	if c.PhoneNo != "" {
		q = append(q, " phoneNo=?")
		args = append(args, c.PhoneNo)
	}

	if q == nil {
		return "", args
	}

	args = append(args, id)
	query := "UPDATE customers SET" + strings.Join(q, ",") + " WHERE id = ?;"
	return query, args
}
