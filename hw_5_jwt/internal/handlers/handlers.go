package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"hw_5_jwt/internal/models"
	"hw_5_jwt/internal/postgres"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo   *postgres.Repository
	logger *slog.Logger
}

func NewHandler(repo *postgres.Repository, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return &Handler{repo: repo, logger: logger}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", h.HealthCheck)
	e.POST("/api/auth/register", h.Register)
	e.POST("/api/auth/login", h.Login)

	protected := e.Group("/api")
	protected.Use(h.AuthMiddleware)
	{
		protected.GET("/users/me", h.GetCurrentUser)

		protected.GET("/students", h.GetAllStudents)
		protected.GET("/students/:id", h.GetStudent)
		protected.GET("/schedule", h.GetAllSchedule)
		protected.GET("/schedule/group/:id", h.GetGroupSchedule)
		protected.GET("/groups", h.GetAllGroups)
		protected.GET("/groups/:id", h.GetGroup)
		protected.POST("/attendance/subject", h.CreateAttendance)
		protected.GET("/attendanceBySubjectId/:id", h.GetAttendanceBySubjectID)
		protected.GET("/attendanceByStudentId/:id", h.GetAttendanceByStudentID)
	}
}

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, models.ServerResponse{
				Status:  "error",
				Message: "Authorization header required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, models.ServerResponse{
				Status:  "error",
				Message: "Invalid authorization format. Use Bearer <token>",
			})
		}

		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			h.logger.Warn("невалидный токен", "error", err)
			return c.JSON(http.StatusUnauthorized, models.ServerResponse{
				Status:  "error",
				Message: "Invalid or expired token",
			})
		}

		c.Set("userID", claims.UserID)
		return next(c)
	}
}

func (h *Handler) Register(c echo.Context) error {
	var req models.RegisterRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Warn("ошибка валидации регистрации", "error", err)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверные данные",
			Error:   err.Error(),
		})
	}

	existingUser, err := h.repo.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		h.logger.Error("ошибка при проверке пользователя", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка сервера",
		})
	}

	if existingUser != nil {
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Пользователь с таким email уже существует",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("ошибка хеширования пароля", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка при обработке пароля",
		})
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := h.repo.CreateUser(c.Request().Context(), user)
	if err != nil {
		h.logger.Error("ошибка создания пользователя", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Не удалось создать пользователя",
			Error:   err.Error(),
		})
	}

	token, err := GenerateToken(createdUser.ID)
	if err != nil {
		h.logger.Error("ошибка генерации токена", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Не удалось создать токен",
		})
	}

	createdUser.Password = ""

	return c.JSON(http.StatusCreated, models.ServerResponse{
		Status:  "success",
		Message: "Пользователь успешно зарегистрирован",
		Data: map[string]interface{}{
			"token": token,
			"user":  createdUser,
		},
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req models.LoginRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Warn("ошибка валидации входа", "error", err)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверные данные",
			Error:   err.Error(),
		})
	}

	user, err := h.repo.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		h.logger.Error("ошибка при получении пользователя", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка сервера",
		})
	}

	if user == nil {
		h.logger.Warn("пользователь не найден", "email", req.Email)
		return c.JSON(http.StatusUnauthorized, models.ServerResponse{
			Status:  "error",
			Message: "Неверный email или пароль",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		h.logger.Warn("неверный пароль", "email", req.Email)
		return c.JSON(http.StatusUnauthorized, models.ServerResponse{
			Status:  "error",
			Message: "Неверный email или пароль",
		})
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		h.logger.Error("ошибка генерации токена", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Не удалось создать токен",
		})
	}

	user.Password = ""

	return c.JSON(http.StatusOK, models.ServerResponse{
		Status:  "success",
		Message: "Успешный вход",
		Data: map[string]interface{}{
			"token": token,
			"user":  user,
		},
	})
}

func (h *Handler) GetCurrentUser(c echo.Context) error {

	userID := c.Get("userID")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, models.ServerResponse{
			Status:  "error",
			Message: "Пользователь не аутентифицирован",
		})
	}

	user, err := h.repo.GetUserByID(c.Request().Context(), userID.(int))
	if err != nil {
		h.logger.Error("ошибка получения пользователя", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка сервера",
		})
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, models.ServerResponse{
			Status:  "error",
			Message: "Пользователь не найден",
		})
	}

	user.Password = ""

	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   user,
	})
}

func (h *Handler) CreateAttendance(c echo.Context) error {
	var req models.AttendanceRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Warn("ошибка привязки данных", "error", err)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверный формат данных",
		})
	}

	if req.ScheduleID == 0 {
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "schedule_id обязателен",
		})
	}

	if req.VisitDay == "" {
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "visit_day обязателен",
		})
	}

	if req.StudentID == 0 {
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "student_id обязателен",
		})
	}

	normalizedDate, err := normalizeDate(req.VisitDay)
	if err != nil {
		h.logger.Warn("неверный формат даты", "date", req.VisitDay, "error", err)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверный формат даты. Используйте формат DD.MM.YYYY",
		})
	}

	req.VisitDay = normalizedDate

	h.logger.Info("создание записи посещаемости",
		"schedule_id", req.ScheduleID,
		"student_id", req.StudentID,
		"visit_day", req.VisitDay,
		"visited", req.Visited,
	)

	if err := h.repo.CreateAttendance(c.Request().Context(), req); err != nil {
		h.logger.Error("ошибка создания посещаемости", "error", err.Error()) // ← Добавь .Error()
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка создания записи посещаемости: " + err.Error(), // ← Добавь err.Error()
		})
	}

	h.logger.Info("запись посещаемости успешно создана")
	return c.JSON(http.StatusCreated, models.ServerResponse{
		Status:  "success",
		Message: "Запись посещаемости успешно создана",
	})
}

func normalizeDate(dateStr string) (string, error) {

	if t, err := time.Parse("02.01.2006", dateStr); err == nil {
		return t.Format("02.01.2006"), nil
	}

	formats := []string{
		"2006/01/02",
		"02-01-2006",
		"02/01/2006",
		"01/02/2006",
		"2006.01.02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("02.01.2006"), nil
		}
	}

	return "", fmt.Errorf("неподдерживаемый формат даты")
}

func (h *Handler) GetAttendanceBySubjectID(c echo.Context) error {
	idStr := c.Param("id")

	subjectID, err := strconv.Atoi(idStr)
	if err != nil || subjectID <= 0 {
		h.logger.Warn("неверный формат ID предмета", "id", idStr)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверный формат ID предмета",
		})
	}

	h.logger.Info("получение посещаемости по предмету", "subject_id", subjectID)

	attendances, err := h.repo.GetAttendanceBySubjectID(c.Request().Context(), subjectID)
	if err != nil {
		h.logger.Error("ошибка получения посещаемости", "subject_id", subjectID, "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка получения посещаемости",
		})
	}

	if len(attendances) == 0 {
		h.logger.Info("посещаемость не найдена", "subject_id", subjectID)
		return c.JSON(http.StatusOK, models.ServerResponse{
			Status:  "success",
			Message: "Посещаемость по данному предмету не найдена",
			Data:    []models.AttendanceBySubject{},
		})
	}

	h.logger.Info("посещаемость успешно получена", "subject_id", subjectID, "count", len(attendances))
	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   attendances,
	})
}

func (h *Handler) GetAttendanceByStudentID(c echo.Context) error {
	idStr := c.Param("id")

	studentID, err := strconv.Atoi(idStr)
	if err != nil || studentID <= 0 {
		h.logger.Warn("неверный формат ID студента", "id", idStr)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверный формат ID студента",
		})
	}

	h.logger.Info("получение посещаемости по студенту", "student_id", studentID)

	attendances, err := h.repo.GetAttendanceByStudentID(c.Request().Context(), studentID)
	if err != nil {
		h.logger.Error("ошибка получения посещаемости", "student_id", studentID, "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка получения посещаемости",
		})
	}

	if len(attendances) == 0 {
		h.logger.Info("посещаемость не найдена", "student_id", studentID)
		return c.JSON(http.StatusOK, models.ServerResponse{
			Status:  "success",
			Message: "Посещаемость данного студента не найдена",
			Data:    []models.AttendanceByStudent{},
		})
	}

	h.logger.Info("посещаемость успешно получена", "student_id", studentID, "count", len(attendances))
	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   attendances,
	})
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   "API работает ура",
	})
}

func (h *Handler) GetStudent(c echo.Context) error {
	idStr := c.Param("id")

	studentID, err := strconv.Atoi(idStr)
	if err != nil || studentID <= 0 {
		h.logger.Warn("неверный формат ID студента", "id", idStr)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверный формат ID",
		})
	}

	h.logger.Info("получение студента", "id", studentID)

	student, err := h.repo.GetStudent(c.Request().Context(), studentID)
	if err != nil {
		h.logger.Error("ошибка получения студента", "id", studentID, "error", err)

		if err.Error() == "студент с ID %d не найден" {
			return c.JSON(http.StatusNotFound, models.ServerResponse{
				Status:  "error",
				Message: "Студент не найден",
			})
		}

		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Внутренняя ошибка сервера",
		})
	}

	h.logger.Info("студент успешно получен", "id", studentID)
	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   student,
	})
}

func (h *Handler) GetAllStudents(c echo.Context) error {
	h.logger.Info("получение всех студентов")

	students, err := h.repo.GetAllStudents(c.Request().Context())
	if err != nil {
		h.logger.Error("ошибка получения студентов", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка получения студентов",
		})
	}

	h.logger.Info("студенты успешно получены", "count", len(students))
	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   students,
	})
}

func (h *Handler) GetAllSchedule(c echo.Context) error {
	h.logger.Info("получение всего расписания")

	_, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	schedule, err := h.repo.GetAllSchedule(c.Request().Context())
	if err != nil {
		h.logger.Error("ошибка получения расписания",
			"error", err,
			"error_type", fmt.Sprintf("%T", err),
			"stack", string(debug.Stack()),
		)
	}

	h.logger.Info("расписание успешно получено", "count", len(schedule))
	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   schedule,
	})
}

func (h *Handler) GetGroupSchedule(c echo.Context) error {
	idStr := c.Param("id")

	groupID, err := strconv.Atoi(idStr)
	if err != nil || groupID <= 0 {
		h.logger.Warn("неверный формат ID группы", "id", idStr)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверный формат ID группы",
		})
	}

	h.logger.Info("получение расписания группы", "group_id", groupID)

	schedule, err := h.repo.GetGroupSchedule(c.Request().Context(), groupID)
	if err != nil {
		h.logger.Error("ошибка получения расписания группы", "group_id", groupID, "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка получения расписания",
		})
	}

	h.logger.Info("расписание группы успешно получено",
		"group_id", groupID, "count", len(schedule))

	if len(schedule) == 0 {
		return c.JSON(http.StatusOK, models.ServerResponse{
			Status:  "success",
			Message: "Расписание для группы пустое",
			Data:    []models.Schedule{},
		})
	}

	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   schedule,
	})
}

func (h *Handler) GetAllGroups(c echo.Context) error {
	h.logger.Info("получение всех групп")

	groups, err := h.repo.GetGroups(c.Request().Context())
	if err != nil {
		h.logger.Error("ошибка получения групп", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Ошибка получения групп",
		})
	}

	h.logger.Info("группы успешно получены", "count", len(groups))
	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   groups,
	})
}

func (h *Handler) GetGroup(c echo.Context) error {
	idStr := c.Param("id")

	groupID, err := strconv.Atoi(idStr)
	if err != nil || groupID <= 0 {
		h.logger.Warn("неверный формат ID группы", "id", idStr)
		return c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status:  "error",
			Message: "Неверный формат ID",
		})
	}

	h.logger.Info("получение группы", "id", groupID)

	group, err := h.repo.GetGroup(c.Request().Context(), groupID)
	if err != nil {
		h.logger.Error("ошибка получения группы", "id", groupID, "error", err)

		if err.Error() == "группа с ID %d не найдена" {
			return c.JSON(http.StatusNotFound, models.ServerResponse{
				Status:  "error",
				Message: "Группа не найдена",
			})
		}

		return c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status:  "error",
			Message: "Внутренняя ошибка сервера",
		})
	}

	h.logger.Info("группа успешно получена", "id", groupID)
	return c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   group,
	})
}
