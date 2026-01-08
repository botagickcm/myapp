package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://akbota:4912@localhost:5432/university_test"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	fmt.Println("You connected to database")
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// e.GET("/student/:id", h.GetStudent)
	// e.GET("/all_class_schedule", h.GetAllSchedules)
	// e.GET("/schedule/group/:id", h.GetGroupSchedule)

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

}
