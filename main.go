package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sql.Open("sqlite3", "file:sc25k.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	var foreignKeysEnabled bool

	if err := conn.QueryRow(
		"PRAGMA foreign_keys;",
	).Scan(&foreignKeysEnabled); err != nil {
		log.Fatal(err)
	}

	if !foreignKeysEnabled {
		log.Fatal("SQLite foreign keys are disabled")
	}

	db := &Database{
		conn: conn,
	}

	app := &App{
		DB: db,
	}

	err = db.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database is ready")

	// stage, err := NewStage(1, 1, "Week 1 Day 1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = db.CreateStage(stage)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stages, err := db.GetStages()
	// log.Println(stages)

	// stage, err := db.GetStage(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(stage)

	// completion, err := newCompletion(1, 1, 100, "path/to/file")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = db.CompleteStage(completion)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(completion)

	// completedStages, err := db.GetCompletedStages()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(completedStages)

	// cycles, err := db.GetCyclesByStageID(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(cycles)

	err = http.ListenAndServe(":8080", app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}
