
-- Table: public.attendance

-- DROP TABLE IF EXISTS public.attendance;

CREATE TABLE IF NOT EXISTS public.attendance
(
    attendance_id integer NOT NULL DEFAULT nextval('attendance_attendance_id_seq'::regclass),
    student_id integer NOT NULL,
    attendance_date date NOT NULL DEFAULT CURRENT_DATE,
    is_present boolean NOT NULL DEFAULT false,
    schedule_id integer NOT NULL,
    lesson_name text COLLATE pg_catalog."default",
    CONSTRAINT attendance_pkey PRIMARY KEY (attendance_id),
    CONSTRAINT attendance_student_id_schedule_id_attendance_date_key UNIQUE (student_id, schedule_id, attendance_date),
    CONSTRAINT attendance_schedule_id_fkey FOREIGN KEY (schedule_id)
        REFERENCES public.schedule (schedule_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT attendance_student_id_fkey FOREIGN KEY (student_id)
        REFERENCES public.students (student_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.attendance
    OWNER to postgres;


    
CREATE TABLE IF NOT EXISTS public.subjects (
    subject_id SERIAL PRIMARY KEY,
    subject_name TEXT NOT NULL UNIQUE,
);

CREATE TABLE IF NOT EXISTS public.groups
(
    group_id integer NOT NULL DEFAULT nextval('groups_id_seq'::regclass),
    group_name text COLLATE pg_catalog."default" NOT NULL,
    faculty text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT groups_pkey PRIMARY KEY (group_id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.groups
    OWNER to postgres;


CREATE TABLE IF NOT EXISTS public.schedule
(
    schedule_id integer NOT NULL DEFAULT nextval('schedule_schedule_id_seq'::regclass),
    day_of_week character varying(20) COLLATE pg_catalog."default" NOT NULL,
    group_id integer,
    faculty text COLLATE pg_catalog."default",
    start_time time without time zone,
    end_time time without time zone,
    lesson_name text COLLATE pg_catalog."default",
    CONSTRAINT schedule_pkey PRIMARY KEY (schedule_id),
    CONSTRAINT group_id FOREIGN KEY (group_id)
        REFERENCES public.groups (group_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT day_of_week_check CHECK (day_of_week::text = ANY (ARRAY['Понедельник'::character varying, 'Вторник'::character varying, 'Среда'::character varying, 'Четверг'::character varying, 'Пятница'::character varying, 'Суббота'::character varying]::text[]))
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.schedule
    OWNER to postgres;



CREATE TABLE IF NOT EXISTS public.students
(
    name text COLLATE pg_catalog."default",
    student_id integer NOT NULL DEFAULT nextval('students_id_seq'::regclass),
    surname text COLLATE pg_catalog."default",
    gender character varying(1) COLLATE pg_catalog."default",
    birthday date,
    group_id integer,
    CONSTRAINT fk_students_group FOREIGN KEY (group_id)
        REFERENCES public.groups (group_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT gender_check CHECK (gender::text = ANY (ARRAY['М'::character varying, 'Ж'::character varying]::text[]))
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.students
    OWNER to postgres;


INSERT INTO public.subjects(subject_id, subject_name) VALUES
    (1, 'Математика'),
    (2, 'Физика'),
    (3, 'Программирование'),
    (4, 'Базы данных'),
    (5, 'История'),
    (6, 'Английский язык'),
    (7, 'Химия'),
    (8, 'Биология'),
    (9, 'Физкультура'),
    (10, 'Экономика'),
    (11, 'Информатика'),
    (12, 'Алгоритмы'),
    (13, 'Сети'),
    (14, 'Веб-разработка'),
    (15, 'Литература'),
    (16, 'Философия'),
    (17, 'Искусство'),
    (18, 'Психология'),
    (19, 'Социология'),
    (20, 'Политология');


INSERT INTO public.schedule(
    day_of_week, 
    group_id, 
    faculty, 
    start_time, 
    end_time, 
    lesson_name, 
    subject_id
) VALUES 
    ('Понедельник', 1, 'Инженеры', '09:00', '10:30', 'Математика', 1),
    ('Понедельник', 1, 'Инженеры', '10:45', '12:15', 'Физика', 2),
    ('Понедельник', 1, 'Инженеры', '13:00', '14:30', 'Программирование', 3),
    

    ('Вторник', 1, 'Инженеры', '09:00', '10:30', 'Базы данных', 4),
    ('Вторник', 1, 'Инженеры', '10:45', '12:15', 'История', 5),
    ('Вторник', 1, 'Инженеры', '13:00', '14:30', 'Английский язык', 6),

    ('Среда', 1, 'Инженеры', '09:00', '10:30', 'Химия', 7),
    ('Среда', 1, 'Инженеры', '10:45', '12:15', 'Биология', 8),
    ('Среда', 1, 'Инженеры', '13:00', '14:30', 'Физкультура', 9),
    

    ('Четверг', 1, 'Инженеры', '09:00', '10:30', 'Экономика', 10),
    ('Четверг', 1, 'Инженеры', '10:45', '12:15', 'Информатика', 11),
    ('Четверг', 1, 'Инженеры', '13:00', '14:30', 'Алгоритмы', 12),
    
   
    ('Пятница', 1, 'Инженеры', '09:00', '10:30', 'Сети', 13),
    ('Пятница', 1, 'Инженеры', '10:45', '12:15', 'Веб-разработка', 14),

    ('Понедельник', 2, 'Инженеры', '09:00', '10:30', 'Программирование', 3),
    ('Понедельник', 2, 'Инженеры', '10:45', '12:15', 'Базы данных', 4),
    ('Понедельник', 2, 'Инженеры', '13:00', '14:30', 'Математика', 1),
    
    ('Вторник', 2, 'Инженеры', '09:00', '10:30', 'Физика', 2),
    ('Вторник', 2, 'Инженеры', '10:45', '12:15', 'Английский язык', 6),
    ('Вторник', 2, 'Инженеры', '13:00', '14:30', 'История', 5),

    ('Понедельник', 4, 'Гуманитарные науки', '09:00', '10:30', 'История', 5),
    ('Понедельник', 4, 'Гуманитарные науки', '10:45', '12:15', 'Английский язык', 6),
    ('Понедельник', 4, 'Гуманитарные науки', '13:00', '14:30', 'Литература', 15),
    
    
    ('Вторник', 4, 'Гуманитарные науки', '09:00', '10:30', 'Философия', 16),
    ('Вторник', 4, 'Гуманитарные науки', '10:45', '12:15', 'Искусство', 17);


INSERT INTO public.groups(
	group_name, faculty)
	VALUES ('GR11', 'Инженеры'),
	('GR12', 'Инженеры'),
	('GR13', 'Инженеры'),
	('GR21', 'Гуманитарии'),
	('GR22', 'Гуманитарии');


INSERT INTO public.students(
    name, surname, gender, birthday, group_id) VALUES
    ('Анна', 'Иванова', 'Ж', '2003-05-15', 1),
    ('Мария', 'Петрова', 'Ж', '2002-08-22', 1),
    ('Екатерина', 'Сидорова', 'Ж', '2003-01-10', 2),
    ('Ольга', 'Смирнова', 'Ж', '2002-11-30', 2),
    ('Татьяна', 'Кузнецова', 'Ж', '2003-03-18', 3),
    ('Елена', 'Попова', 'Ж', '2002-07-25', 3),
    ('Наталья', 'Васильева', 'Ж', '2003-04-05', 4),
    ('Юлия', 'Федорова', 'Ж', '2002-12-12', 4),
    ('Александра', 'Морозова', 'Ж', '2003-02-28', 5),
    ('Ирина', 'Николаева', 'Ж', '2002-09-14', 5);







SELECT *
FROM students
WHERE gender ='Ж'
ORDER BY birthday DESC;
