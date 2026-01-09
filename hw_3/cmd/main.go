// cmd/server/main.go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	handler "hw_3/internal/handlers"
	"hw_3/internal/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	logger.Info("запуск университета")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://akbota:4912@localhost:5432/university_test"
	}

	logger.Info("подключение к базе данных", "dsn", connStr)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logger.Error("ошибка подключения к базе данных", "error", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	if err := conn.Ping(context.Background()); err != nil {
		logger.Error("ошибка ping базы данных", "error", err)
		os.Exit(1)
	}

	logger.Info("успешное подключение к базе данных")

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info("запрос",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
				)
			} else {
				logger.Error("ошибка запроса",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"error", v.Error,
				)
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

	repo := postgres.NewRepository(conn)
	h := handler.NewHandler(repo, logger)

	h.RegisterRoutes(e)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("сервер запущен на :8080")
		if err := e.Start("127.0.0.1:8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("ошибка запуска сервера", "error", err)
		}
	}()

	<-quit
	logger.Info("получен сигнал завершения")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("ошибка graceful shutdown", "error", err)
	}

	logger.Info("сервер остановлен")
}
