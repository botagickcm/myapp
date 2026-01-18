package postgres

import (
	"context"
	"fmt"
	"hw_5_jwt/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (email, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, email, created_at
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&user.ID, &user.Email, &user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, created_at 
		FROM users 
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, created_at 
		FROM users 
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

func (r *Repository) CreateAttendance(ctx context.Context, req models.AttendanceRequest) error {

	visitDate, err := time.Parse("02.01.2006", req.VisitDay)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %w", err)
	}

	query := `
		INSERT INTO attendance (student_id, attendance_date, is_present, schedule_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (student_id, schedule_id, attendance_date) 
		DO UPDATE SET is_present = EXCLUDED.is_present
	`

	_, err = r.db.Exec(ctx, query, req.StudentID, visitDate, req.Visited, req.ScheduleID)
	if err != nil {
		return fmt.Errorf("ошибка создания записи посещаемости: %w", err)
	}

	return nil
}

func (r *Repository) GetAttendanceBySubjectID(ctx context.Context, subjectID int) ([]models.AttendanceBySubject, error) {

	query := `
		SELECT 
			s.student_id,
			s.name,
			s.surname,
			g.group_name,
			TO_CHAR(a.attendance_date, 'DD.MM.YYYY') as visit_day,
			a.is_present as visited
		FROM attendance a
		JOIN students s ON a.student_id = s.student_id
		JOIN groups g ON s.group_id = g.group_id
		WHERE a.schedule_id = $1
		ORDER BY s.surname, s.name, a.attendance_date DESC
	`

	rows, err := r.db.Query(ctx, query, subjectID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения посещаемости по предмету: %w", err)
	}
	defer rows.Close()

	var attendances []models.AttendanceBySubject
	for rows.Next() {
		var attendance models.AttendanceBySubject
		err := rows.Scan(
			&attendance.StudentID,
			&attendance.StudentName,
			&attendance.StudentSurname,
			&attendance.GroupName,
			&attendance.VisitDay,
			&attendance.Visited,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования посещаемости: %w", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации посещаемости: %w", err)
	}

	return attendances, nil
}

func (r *Repository) GetAttendanceByStudentID(ctx context.Context, studentID int) ([]models.AttendanceByStudent, error) {

	query := `
		SELECT 
			a.schedule_id as subject_id,
			sch.lesson_name as subject_name,
			TO_CHAR(a.attendance_date, 'DD.MM.YYYY') as visit_day,
			a.is_present as visited
		FROM attendance a
		JOIN schedule sch ON a.schedule_id = sch.schedule_id
		WHERE a.student_id = $1
		ORDER BY a.attendance_date DESC, sch.lesson_name
	`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения посещаемости по студенту: %w", err)
	}
	defer rows.Close()

	var attendances []models.AttendanceByStudent
	for rows.Next() {
		var attendance models.AttendanceByStudent
		err := rows.Scan(
			&attendance.SubjectID,
			&attendance.SubjectName,
			&attendance.VisitDay,
			&attendance.Visited,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования посещаемости: %w", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации посещаемости: %w", err)
	}

	return attendances, nil
}

func (r *Repository) GetStudent(ctx context.Context, id int) (*models.Student, error) {
	query := `
		SELECT student_id, name, surname, gender, birthday, group_id 
		FROM students 
		WHERE student_id = $1
	`

	var student models.Student
	err := r.db.QueryRow(ctx, query, id).Scan(
		&student.StudentID,
		&student.Name,
		&student.Surname,
		&student.Gender,
		&student.Birthday,
		&student.GroupID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("студент с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка получения студента: %w", err)
	}

	return &student, nil
}

func (r *Repository) GetAllSchedule(ctx context.Context) ([]models.Schedule, error) {
	query := `
		SELECT group_id, lesson_name, start_time, end_time
		FROM schedule
		ORDER BY day_of_week, start_time
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения расписания: %w", err)
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.GroupID,
			&schedule.Subject,
			&schedule.StartTime,
			&schedule.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования расписания: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации расписания: %w", err)
	}

	return schedules, nil
}

func (r *Repository) GetGroupSchedule(ctx context.Context, groupID int) ([]models.Schedule, error) {
	query := `
		SELECT group_id, lesson_name, start_time, end_time
		FROM schedule
		WHERE group_id = $1
		ORDER BY day_of_week, start_time
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения расписания группы: %w", err)
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.GroupID,
			&schedule.Subject,
			&schedule.StartTime,
			&schedule.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования расписания: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации расписания: %w", err)
	}

	return schedules, nil
}

func (r *Repository) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	query := `
		SELECT student_id, name, surname, gender, birthday, group_id 
		FROM students 
		ORDER BY student_id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения студентов: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.StudentID,
			&student.Name,
			&student.Surname,
			&student.Gender,
			&student.Birthday,
			&student.GroupID,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования студента: %w", err)
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации студентов: %w", err)
	}

	return students, nil
}

func (r *Repository) GetGroups(ctx context.Context) ([]models.Group, error) {
	query := `
		SELECT group_id, group_name, faculty
		FROM groups
		ORDER BY group_id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения групп: %w", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.GroupID,
			&group.GroupName,
			&group.Faculty,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования группы: %w", err)
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации групп: %w", err)
	}

	return groups, nil
}

func (r *Repository) GetGroup(ctx context.Context, id int) (*models.Group, error) {
	query := `
		SELECT group_id, group_name, faculty
		FROM groups
		WHERE group_id = $1
	`

	var group models.Group
	err := r.db.QueryRow(ctx, query, id).Scan(
		&group.GroupID,
		&group.GroupName,
		&group.Faculty,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("группа с ID %d не найдена", id)
		}
		return nil, fmt.Errorf("ошибка получения группы: %w", err)
	}

	return &group, nil
}
