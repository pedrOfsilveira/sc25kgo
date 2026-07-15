package main

import (
	"database/sql"
	"time"
)

type App struct {
	DB *Database
}

type Database struct {
	conn *sql.DB
}

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Points int    `json:"points"`
}

func NewUser(name string) (User, error) {
	return User{
		Name:   name,
		Points: 0,
	}, nil
}

type Stage struct {
	ID     int          `json:"id"`
	Week   int          `json:"week"`
	Day    int          `json:"day"`
	Name   string       `json:"name"`
	Cycles []StageCycle `json:"cycles"`
}

type StageSummary struct {
	ID   int    `json:"id"`
	Week int    `json:"week"`
	Day  int    `json:"day"`
	Name string `json:"name"`
}

func NewStage(week, day int, name string) (Stage, error) {
	return Stage{
		Week: week,
		Day:  day,
		Name: name,
	}, nil
}

type StageCycle struct {
	ID         int    `json:"id"`
	StageID    int    `json:"stageId"`
	Type       string `json:"type"`
	Duration   int    `json:"duration"`
	CycleOrder int    `json:"cycleOrder"`
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
	ID           int    `json:"id"`
	UserID       int    `json:"userId"`
	StageID      int    `json:"stageId"`
	PhotoURL     string `json:"photoUrl"`
	PointsEarned int    `json:"pointsEarned"`
}

type RunHistoryEntry struct {
	CompletionID int          `json:"completionId"`
	Stage        StageSummary `json:"stage"`
	PointsEarned int          `json:"pointsEarned"`
	PhotoURL     string       `json:"photoUrl"`
	CreatedAt    time.Time    `json:"createdAt"`
}

type ProgressResponse struct {
	User                User          `json:"user"`
	CompletedStageCount int           `json:"completedStageCount"`
	TotalStageCount     int           `json:"totalStageCount"`
	NextStage           *StageSummary `json:"nextStage"`
}

func NewCompletion(userID, stageID, pointsEarned int, photoURL string) (Completion, error) {
	return Completion{
		UserID:       userID,
		StageID:      stageID,
		PointsEarned: pointsEarned,
		PhotoURL:     photoURL,
	}, nil
}
