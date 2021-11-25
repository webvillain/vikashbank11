package memdb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Newdb *sql.DB
)

func Connect() {
	d, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	Newdb = d
}

func GetDb() *sql.DB {
	return Newdb

}
