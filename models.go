package main

import "database/sql"

type App struct {
	DB *Database
}

type Database struct {
	conn *sql.DB
}

type User struct {
	ID     int
	Name   string
	Points int
}

func NewUser(name string) (User, error) {
	return User{
		Name:   name,
		Points: 0,
	}, nil
}

type Stage struct {
	ID     int
	Week   int
	Day    int
	Name   string
	Cycles []StageCycle
}

func NewStage(week, day int, name string) (Stage, error) {
	return Stage{
		Week: week,
		Day:  day,
		Name: name,
	}, nil
}

type StageCycle struct {
	ID         int
	StageID    int
	Type       string
	Duration   int
	CycleOrder int
}

func NewStageCycle(stageID int, cycleType string, duration, cycleOrder int) (StageCycle, error) {
	return StageCycle{
		StageID:    stageID,
		Type:       cycleType,
		Duration:   duration,
		CycleOrder: cycleOrder,
	}, nil
}

type Completion struct {
	ID           int
	UserID       int
	StageID      int
	PhotoURL     string
	PointsEarned int
}

func NewCompletion(userID, stageID, pointsEarned int, photoURL string) (Completion, error) {
	return Completion{
		UserID:       userID,
		StageID:      stageID,
		PointsEarned: pointsEarned,
		PhotoURL:     photoURL,
	}, nil
}
