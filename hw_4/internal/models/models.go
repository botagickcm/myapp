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
type Subject struct {
	ScheduleID int    `json:"schedule_id"`
	LessonName string `json:"lesson_name"`
}
type Attendance struct {
	AttendanceID   int       `json:"attendance_id,omitempty"`
	StudentID      int       `json:"student_id"`
	ScheduleID     int       `json:"schedule_id"`
	AttendanceDate time.Time `json:"attendance_date"`
	IsPresent      bool      `json:"is_present"`
}
type AttendanceRequest struct {
	ScheduleID int    `json:"schedule_id"`
	VisitDay   string `json:"visit_day"`
	Visited    bool   `json:"visited"`
	StudentID  int    `json:"student_id"`
}
type AttendanceBySubject struct {
	StudentID      int    `json:"student_id"`
	StudentName    string `json:"student_name"`
	StudentSurname string `json:"student_surname"`
	GroupName      string `json:"group_name"`
	VisitDay       string `json:"visit_day"`
	Visited        bool   `json:"visited"`
}

type AttendanceByStudent struct {
	SubjectID   int    `json:"subject_id"`
	SubjectName string `json:"subject_name"`
	VisitDay    string `json:"visit_day"`
	Visited     bool   `json:"visited"`
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
