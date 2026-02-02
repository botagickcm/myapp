
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    surname VARCHAR(100),
    role VARCHAR(50) DEFAULT 'student',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS groups (
    group_id SERIAL PRIMARY KEY,
    group_name VARCHAR(50) NOT NULL,
    faculty VARCHAR(100)
);


CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name VARCHAR(100),
    surname VARCHAR(100),
    gender VARCHAR(10),
    subject_id INTEGER
);


CREATE TABLE IF NOT EXISTS students (
    student_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    gender VARCHAR(10),
    birthday DATE,
    group_id INTEGER REFERENCES groups(group_id)
);


CREATE TABLE IF NOT EXISTS subjects (
    subject_id SERIAL PRIMARY KEY,
    subject_name VARCHAR(100) NOT NULL
);


CREATE TABLE IF NOT EXISTS schedule (
    schedule_id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups(group_id),
    subject_id INTEGER REFERENCES subjects(subject_id),
    lesson_name VARCHAR(100),
    day_of_week INTEGER,
    start_time TIME,
    end_time TIME
);


CREATE TABLE IF NOT EXISTS attendance (
    attendance_id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES students(student_id),
    schedule_id INTEGER REFERENCES schedule(schedule_id),
    attendance_date DATE NOT NULL,
    is_present BOOLEAN DEFAULT false,
    UNIQUE(student_id, schedule_id, attendance_date)
);


CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_attendance_student ON attendance(student_id);
CREATE INDEX idx_attendance_schedule ON attendance(schedule_id);
CREATE INDEX idx_schedule_group ON schedule(group_id);