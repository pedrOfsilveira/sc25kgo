package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func (db *Database) createTables() error {
	sqlstmt := `
	CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			points INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS stages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			week INTEGER NOT NULL,
			day INTEGER NOT NULL,
			name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS stage_cycles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stage_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			duration INTEGER NOT NULL,
			cycle_order INTEGER NOT NULL,
			FOREIGN KEY(stage_id) REFERENCES stages(id)
		);

		CREATE TABLE IF NOT EXISTS run_completions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			stage_id INTEGER NOT NULL,
			photo_url TEXT NOT NULL,
			points_earned INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.conn.Exec(sqlstmt)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (db *Database) getStages() ([]Stage, error) {
	sqlstmt := `SELECT id, week, day, name FROM stages;`

	rows, err := db.conn.Query(sqlstmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var stages []Stage
	for rows.Next() {
		var stage Stage

		err := rows.Scan(
			&stage.ID,
			&stage.Week,
			&stage.Day,
			&stage.Name,
		)
		if err != nil {
			log.Fatal(err)
		}
		stages = append(stages, stage)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return stages, nil
}

func (db *Database) createStage(s Stage) error {
	_, err := db.conn.Exec(`
		INSERT INTO stages(week, day, name)
		VALUES (?, ?, ?);
		`,
		s.Week, s.Day, s.Name,
	)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (db *Database) getStage(id int) (Stage, error) {
	row := db.conn.QueryRow(`
	SELECT id, week, day, name
	FROM stages
	WHERE id = ?;
	`, id)

	var stage Stage
	err := row.Scan(
		&stage.ID,
		&stage.Week,
		&stage.Day,
		&stage.Name,
	)
	if err != nil {
		log.Fatal(err)
	}

	return stage, err
}

func (db *Database) completeStage(c Completion) error {
	_, err := db.conn.Exec(`
	INSERT INTO run_completions(user_id, stage_id, photo_url, points_earned)
	VALUES (?, ?, ?, ?);`,
		c.UserID,
		c.StageID,
		c.PhotoURL,
		c.PointsEarned,
	)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (db *Database) getCompletedStages() ([]Stage, error) {
	rows, err := db.conn.Query(`
	SELECT DISTINCT
	    stages.id,
	    stages.week,
	    stages.day,
	    stages.name
    FROM stages
    JOIN run_completions
    ON run_completions.stage_id = stages.id;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var stages []Stage
	for rows.Next() {
		var stage Stage

		rows.Scan(
			&stage.ID,
			&stage.Week,
			&stage.Day,
			&stage.Name,
		)

		stages = append(stages, stage)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return stages, nil
}

func (db *Database) getCyclesByStageID(id int) ([]StageCycle, error) {
	rows, err := db.conn.Query(`
	SELECT
		id,
		stage_id,
		type,
		duration
	FROM stage_cycles
	WHERE stage_id = ?
	ORDER BY cycle_order;`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cycles []StageCycle
	for rows.Next() {
		var cycle StageCycle

		rows.Scan(
			&cycle.ID,
			&cycle.StageID,
			&cycle.Type,
			&cycle.Duration,
		)

		cycles = append(cycles, cycle)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return cycles, err
}
