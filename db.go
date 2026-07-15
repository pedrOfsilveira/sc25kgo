package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func (db *Database) CreateTables() error {
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
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(stage_id) REFERENCES stages(id)
		);

		INSERT INTO users (name, points)
		SELECT 'Runner', 0
		WHERE NOT EXISTS (
			SELECT 1 FROM users
		);
`

	_, err := db.conn.Exec(sqlstmt)
	return err
}

func (db *Database) GetStages() ([]StageSummary, error) {
	sqlstmt := `SELECT id, week, day, name FROM stages;`

	rows, err := db.conn.Query(sqlstmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stages := make([]StageSummary, 0)

	for rows.Next() {
		var stage StageSummary

		err := rows.Scan(
			&stage.ID,
			&stage.Week,
			&stage.Day,
			&stage.Name,
		)
		if err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stages, nil
}

func (db *Database) CreateStage(s Stage) error {
	_, err := db.conn.Exec(`
		INSERT INTO stages(week, day, name)
		VALUES (?, ?, ?);
		`,
		s.Week, s.Day, s.Name,
	)
	return err
}

func (db *Database) GetStage(id int) (Stage, error) {
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
		return Stage{}, err
	}

	cycles, err := db.GetCyclesByStageID(stage.ID)
	if err != nil {
		return Stage{}, err
	}

	stage.Cycles = cycles

	return stage, nil
}

func (db *Database) CompleteStage(ctx context.Context, userID int, stage Stage, photoURL string) (Completion, error) {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return Completion{}, fmt.Errorf("begin completion transaction: %w", err)
	}
	defer tx.Rollback()

	var isRepeat bool

	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM run_completions
			WHERE user_id = ? AND stage_id = ?
		);
	`, userID, stage.ID).Scan(&isRepeat)

	if err != nil {
		return Completion{}, fmt.Errorf("check if stage is repeat: %w", err)
	}

	pointsEarned := calculateCompletionPoints(stage, isRepeat)

	result, err := tx.ExecContext(ctx, `
		INSERT INTO run_completions (
		user_id,
		stage_id,
		photo_url,
		points_earned)
		VALUES (?, ?, ?, ?);`, userID, stage.ID, photoURL, pointsEarned)

	if err != nil {
		return Completion{}, fmt.Errorf("insert completion: %w", err)
	}

	completionID, err := result.LastInsertId()
	if err != nil {
		return Completion{}, fmt.Errorf("get completion id: %w", err)
	}

	result, err = tx.ExecContext(ctx, `
		UPDATE users
		SET points = points + ?
		WHERE id = ?;`, pointsEarned, userID)

	if err != nil {
		return Completion{}, fmt.Errorf("update user points: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Completion{}, fmt.Errorf("check updated user: %w", err)
	}

	if rowsAffected != 1 {
		return Completion{}, fmt.Errorf("update user points: %w", sql.ErrNoRows)
	}

	if err := tx.Commit(); err != nil {
		return Completion{}, fmt.Errorf("commit completion transaction: %w", err)
	}

	return Completion{
		ID:           int(completionID),
		UserID:       userID,
		StageID:      stage.ID,
		PhotoURL:     photoURL,
		PointsEarned: pointsEarned,
	}, nil
}

func (db *Database) GetRunHistory(userID int) ([]RunHistoryEntry, error) {
	rows, err := db.conn.Query(`
		SELECT
			run_completions.id,
			stages.id,
			stages.week,
			stages.day,
			stages.name,
			run_completions.points_earned,
			run_completions.photo_url,
			run_completions.created_at
		FROM run_completions
		JOIN stages ON stages.id = run_completions.stage_id
		WHERE run_completions.user_id = ?
		ORDER BY run_completions.created_at DESC, run_completions.id DESC;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	runs := make([]RunHistoryEntry, 0)

	for rows.Next() {
		var run RunHistoryEntry
		if err := rows.Scan(
			&run.CompletionID,
			&run.Stage.ID,
			&run.Stage.Week,
			&run.Stage.Day,
			&run.Stage.Name,
			&run.PointsEarned,
			&run.PhotoURL,
			&run.CreatedAt,
		); err != nil {
			return nil, err
		}

		runs = append(runs, run)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return runs, nil
}

func (db *Database) GetCompletedStages(userID int) ([]StageSummary, error) {
	rows, err := db.conn.Query(
		`SELECT DISTINCT
		stages.id,
		stages.week,
		stages.day,
		stages.name
	FROM stages
	JOIN run_completions
		ON run_completions.stage_id = stages.id
	WHERE run_completions.user_id = ?
	ORDER BY stages.week, stages.day;`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stages := make([]StageSummary, 0)

	for rows.Next() {
		var stage StageSummary

		if err := rows.Scan(
			&stage.ID,
			&stage.Week,
			&stage.Day,
			&stage.Name,
		); err != nil {
			return nil, err
		}

		stages = append(stages, stage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stages, nil
}

func (db *Database) GetCyclesByStageID(id int) ([]StageCycle, error) {
	rows, err := db.conn.Query(`
	SELECT
		id,
		stage_id,
		type,
		duration,
		cycle_order
	FROM stage_cycles
	WHERE stage_id = ?
	ORDER BY cycle_order;`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cycles := make([]StageCycle, 0)

	for rows.Next() {
		var cycle StageCycle

		err := rows.Scan(
			&cycle.ID,
			&cycle.StageID,
			&cycle.Type,
			&cycle.Duration,
			&cycle.CycleOrder,
		)
		if err != nil {
			return nil, err
		}

		cycles = append(cycles, cycle)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return cycles, nil
}

func (db *Database) CreateUser(user User) error {
	_, err := db.conn.Exec(`
	INSERT INTO users(name, points)
	VALUES (?, ?)`, user.Name, user.Points)
	return err
}

func (db *Database) GetUsers() ([]User, error) {
	rows, err := db.conn.Query(`
	SELECT id, name, points
	FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Points,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *Database) GetActiveUser() (User, error) {
	row := db.conn.QueryRow(`
		SELECT id, name, points
		FROM users
		ORDER BY id
		LIMIT 1;`)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Points); err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) GetUser(id int) (User, error) {
	row := db.conn.QueryRow(`
	SELECT id, name, points
	FROM users
	WHERE id = ?`, id)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Points,
	)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *Database) GetCompletedStageCount(userID int) (int, error) {
	var count int

	err := db.conn.QueryRow(
		`SELECT COUNT(DISTINCT stage_id)
				FROM run_completions
				WHERE user_id = ?;
			`, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *Database) GetTotalStageCount() (int, error) {
	var count int
	if err := db.conn.QueryRow(`SELECT COUNT(*) FROM stages;`).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (db *Database) GetNextStage(userID int) (StageSummary, error) {
	row := db.conn.QueryRow(`
		SELECT stages.id, stages.week, stages.day, stages.name
		FROM stages
		WHERE NOT EXISTS (
			SELECT 1
			FROM run_completions
			WHERE run_completions.user_id = ?
				AND run_completions.stage_id = stages.id
		)
		ORDER BY stages.week, stages.day, stages.id
		LIMIT 1;`, userID)

	var stage StageSummary
	if err := row.Scan(&stage.ID, &stage.Week, &stage.Day, &stage.Name); err != nil {
		return StageSummary{}, err
	}

	return stage, nil
}
