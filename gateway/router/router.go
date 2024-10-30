package router

import (
	handler "api-gateway/handlers"
	middleware "api-gateway/middlewares"

	"github.com/labstack/echo/v4"
)




func InitRoutes(e *echo.Echo) {

	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)

	e.GET("/hotel", handler.GetListHotelsHandler)
	e.GET("/hotel/:id", handler.GetHotelsHandler)

	e.GET("/hotel/room", handler.ListRoomsByHotelIdHandler)
	e.GET("/room/:id", handler.GetRoomByIDHandler)
	

	user := e.Group("/api")
	user.Use(middleware.Authentication, middleware.UserAuth) 
	{
		user.GET("/user", handler.GetUserByIDHandler)

		user.POST("/booking", handler.CreateBookingHandler)
		user.GET("/booking", handler.GetListBooking)
		user.GET("/booking/detail/:booking_id", handler.GetDetailBooking)

		user.POST("/payment", handler.CreatePaymentHandler)
		user.POST("/refund/:booking_id", handler.CreateRefundHandler)
	}

	admin := e.Group("/api")
	admin.Use(middleware.Authentication, middleware.AdminAuth)
	{
		admin.POST("/hotel", handler.CreateHotelHandler)
		admin.POST("/room", handler.CreateRoomHandler)
		admin.PUT("/booking/checkin-status", handler.UpdateCheckinStatusHandler)

	}	
}