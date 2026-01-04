package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", salem)
	e.GET("/check", checkSayHey)

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

func salem(c echo.Context) error {
	response := Response{
		Message: "Salem alem",
		Status:  "keremet",
	}
	return c.JSON(http.StatusOK, response)
}
func checkSayHey(c echo.Context) error {
	sayhey := c.QueryParam("sayhey")
	if sayhey == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "Empty, you need to say salem",
			Status:  "bad error code: 400?",
		})
	} else if sayhey == "salem" {
		return c.JSON(http.StatusOK, Response{
			Message: "Yeah, you said salem",
			Status:  "keremet, code: 200",
		})
	}
	return c.JSON(http.StatusOK, Response{
		Message: "Hm, you forgot something, say salem",
		Status:  "keremet, code: 200",
	})

}
