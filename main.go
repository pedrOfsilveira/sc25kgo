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

	// stage, err := NewStage(1, 1, "Week 1 Day 1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = db.createStage(stage)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stages, err := db.getStages()
	// log.Println(stages)

	// stage, err := db.getStage(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(stage)
}
