package postgres

import (
	"context"
	"fmt"
	"hw_3/internal/models"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
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
