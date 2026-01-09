package models

import "time"

type ServerResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Schedule struct {
	GroupID   int       `json:"group_id"`
	Subject   string    `json:"subject"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type Student struct {
	StudentID int       `json:"student_id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	GroupID   int       `json:"group_id"`
}

type Group struct {
	GroupID   int    `json:"id"`
	GroupName string `json:"name"`
	Faculty   string `json:"department"`
}
