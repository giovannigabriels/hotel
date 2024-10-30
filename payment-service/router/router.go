package router

import (
	"net/http"
	handler "payment-service/handlers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! Server PAYMENT is running.")
	})

	e.POST("/payment", handler.CreatePayment)
	e.POST("/payment/callback", handler.PaymentCallbackHandler)
	e.POST("/refund", handler.CreateRefund)
}
