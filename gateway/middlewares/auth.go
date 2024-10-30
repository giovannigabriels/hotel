package middleware

import (
	"log"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("rahasia")

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(401, echo.Map{
				"message": "invalid token format",
			})
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}


		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Println(ok, "OKEE 1")
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}

		role, ok := claims["role"].(string)
		if !ok {
			log.Println(ok, "OKEE")
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}

		c.Set("id", userID)
		c.Set("role", role)

		return next(c)
	}
}

func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "admin" {
			return c.JSON(403, echo.Map{
				"message": "forbidden",
			})
		}
		return next(c)
	}
}

func UserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "user" {
			return c.JSON(403, echo.Map{
				"message": "forbidden",
			})
		}
		return next(c)
	}
}
