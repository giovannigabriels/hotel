package router

import (
	handler "booking-service/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! Server Booking is running.")
	})
	
	e.GET("/hotel", handler.GetAllHotels)
	e.GET("/hotel/room", handler.ListRoomsByHotelId)
	
	e.GET("/hotel/:id", handler.GetHotelByID)
	e.GET("/room/:id", handler.GetRoomByID)
	
	e.POST("/hotel", handler.CreateHotel)
	e.POST(("/room"), handler.CreateRoom)
	e.POST(("/booking"), handler.CreateBooking)

	e.POST("/booking/status", handler.UpdateBookingStatusRefund)
	e.POST("/booking/cancel", handler.CancelBooking)

	e.GET("/booking/:user_id", handler.GetBookingsByUserID)
	e.GET("/booking/detail/:booking_id", handler.GetBookingByID)

	e.POST("/booking/callback/status", handler.UpdateBookingStatusHandler)

	e.POST("/booking/refund/status", handler.UpdateBookingStatus)
	e.PUT("/booking/checkin-status", handler.UpdateCheckinStatus)
}
