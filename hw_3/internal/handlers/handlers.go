package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"hw_3/internal/models"
	"hw_3/internal/postgres"

	"github.com/labstack/echo/v4"
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

	e.GET("/students", h.GetAllStudents)
	e.GET("/students/:id", h.GetStudent)

	e.GET("/schedule", h.GetAllSchedule)
	e.GET("/schedule/group/:id", h.GetGroupSchedule)

	e.GET("/groups", h.GetAllGroups)
	e.GET("/groups/:id", h.GetGroup)

	e.GET("/health", h.HealthCheck)
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

	// Если расписание пустое
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

// GetAllGroups возвращает все группы
// GET /groups
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
