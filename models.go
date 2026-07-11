package main

type User struct {
	ID   int
	Name string
}

type Stage struct {
	ID     int
	Week   int
	Day    int
	Cycles []StageCycle
}

type StageCycle struct {
	Type     string
	Duration int
}
