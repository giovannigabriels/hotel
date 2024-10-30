package router

import (
	"net/http"
	handler "user-service/handlers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! Server User is running.")
	})

	e.POST("/register", handler.RegisterUser)
	e.POST("/login", handler.LoginUser)

	e.GET("/user/:id", handler.GetUserByID)
}
