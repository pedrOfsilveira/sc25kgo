package main

import "database/sql"

type Database struct {
	conn *sql.DB
}

type User struct {
	ID     int
	Name   string
	Points int
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

type Completion struct {
	ID           int
	UserID       int
	StageID      int
	PhotoURL     string
	PointsEarned string
}
