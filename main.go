package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sql.Open("sqlite3", "./sc25k")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	db := &Database{
		conn: conn,
	}

	err = db.createTables()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database is ready")
}
