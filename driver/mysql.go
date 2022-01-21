package driver

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func ConnectToMySQL() *sql.DB {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "secret123",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "customer",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Println(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	log.Println("Connected")

	return db
}
