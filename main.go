package main

import (
	"database/sql"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "./sc25k")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
