package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ServerResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
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
