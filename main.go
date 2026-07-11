package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sql.Open("sqlite3", "./sc25k.db")
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

	// completion, err := newCompletion(1, 1, 100, "path/to/file")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = db.completeStage(completion)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(completion)

	// completedStages, err := db.getCompletedStages()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(completedStages)

	// cycles, err := db.getCyclesByStageID(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(cycles)
}
