package handler

import (
	"booking-service/config"
	"booking-service/dto"
	model "booking-service/models"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateHotel(c echo.Context) error {
	var req dto.CreateHotelRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	query := `
		INSERT INTO hotels (name, address, city, country, phone_number, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	var hotelID int

	err := config.DB.QueryRow(query, req.Name, req.Address, req.City, req.Country, req.PhoneNumber, req.Email).Scan(&hotelID)
	if err != nil {
		log.Println("Error executing query:", err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to create hotel"})
	}

	return c.JSON(http.StatusCreated, dto.CreateHotelResponse{
		HotelId: hotelID,
		Message:  "Hotel created successfully",
	})
}

func CreateRoom(c echo.Context) error {
	var req dto.CreateRoomRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	hotelCheckQuery := `
		SELECT id FROM hotels WHERE id = $1
	`
	var existingHotelID int

	err := config.DB.QueryRow(hotelCheckQuery, req.HotelID).Scan(&existingHotelID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "hotel not found"})
	}

	roomCheckQuery := `
		SELECT id FROM rooms WHERE hotel_id = $1 AND room_number = $2
	`
	var existingRoomID int
	err = config.DB.QueryRow(roomCheckQuery, req.HotelID, req.RoomNumber).Scan(&existingRoomID)
	if err == nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "room_number already exists"})
	}

	insertRoomQuery := `
		INSERT INTO rooms (hotel_id, room_number, room_type, price_per_night, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	var roomID int

	err = config.DB.QueryRow(insertRoomQuery, req.HotelID, req.RoomNumber, req.RoomType, req.PricePerNight, req.Description, req.Status).Scan(&roomID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to create room"})
	}

	return c.JSON(http.StatusCreated, dto.CreateRoomResponse{
		RoomId:  roomID,
		Message:  "Room created successfully",
	})
}

func CreateBooking(c echo.Context) error {
	var req dto.CreateBookingRequest

	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	checkinDate, err := time.Parse("2006-01-02", req.CheckinDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid checkin_date format"})
	}

	checkoutDate, err := time.Parse("2006-01-02", req.CheckoutDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid checkout_date format"})
	}

	var roomStatus string
	checkRoomQuery := `
		SELECT status FROM rooms WHERE id = $1
	`
	err = config.DB.QueryRow(checkRoomQuery, req.RoomID).Scan(&roomStatus)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Room not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to check room availability"})
	}

	if roomStatus != "available" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Room is not available"})
	}

	insertBookingQuery := `
		INSERT INTO bookings (user_id, room_id, checkin_date, checkout_date, total_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	var bookingID int
	err = config.DB.QueryRow(insertBookingQuery, req.UserID, req.RoomID, checkinDate, checkoutDate, req.TotalPrice, "pending").Scan(&bookingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to create booking"})
	}

	updateRoomStatusQuery := `
		UPDATE rooms SET status = 'booked', updated_at = NOW() WHERE id = $1
	`
	_, err = config.DB.Exec(updateRoomStatusQuery, req.RoomID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update room status"})
	}

	return c.JSON(http.StatusCreated, dto.CreateBookingResponse{
		BookingId: bookingID,
		Message:   "Booking created successfully",
	})
}


func GetAllHotels(c echo.Context) error {
	query := `SELECT id, name, address, city, country, phone_number, email, created_at, updated_at FROM hotels`

	var hotels []model.Hotel

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Print(err, "ERRORNYA")
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve hotels"})
	}
	defer rows.Close()

	for rows.Next() {
		var hotel model.Hotel
		if err := rows.Scan(&hotel.HotelID, &hotel.Name, &hotel.Address, &hotel.City, &hotel.Country, &hotel.PhoneNumber, &hotel.Email, &hotel.CreatedAt, &hotel.UpdatedAt); 
		err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to scan hotel data"})
		}
		hotels = append(hotels, hotel)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Error occurred during hotels retrieval"})
	}

	return c.JSON(http.StatusOK, hotels)
}

func GetHotelByID(c echo.Context) error {
	id := c.Param("id")

	query := `SELECT id, name, address, city, country, phone_number, email, created_at, updated_at FROM hotels WHERE id = $1`

	var hotel model.Hotel

	err := config.DB.QueryRow(query, id).Scan(
		&hotel.HotelID,
		&hotel.Name,
		&hotel.Address,
		&hotel.City,
		&hotel.Country,
		&hotel.PhoneNumber,
		&hotel.Email,
		&hotel.CreatedAt,
		&hotel.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Hotel not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve hotel"})
	}

	return c.JSON(http.StatusOK, hotel)
}

func GetRoomByID(c echo.Context) error {
	id := c.Param("id")

	query := `SELECT id, hotel_id, room_number, room_type, price_per_night, description, status, created_at, updated_at FROM rooms WHERE id = $1`

	var room model.Room

	err := config.DB.QueryRow(query, id).Scan(
		&room.RoomID,
		&room.HotelID,
		&room.RoomNumber,
		&room.RoomType,
		&room.PricePerNight,
		&room.Description,
		&room.Status,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Room not found"})
	} else if err != nil {
		log.Println("Error retrieving room:", err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve room"})
	}

	return c.JSON(http.StatusOK, room)
}

func ListRoomsByHotelId(c echo.Context) error {
	hotelID := c.QueryParam("hotel_id")

	if hotelID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "hotel_id is required"})
	}

	query := `SELECT id, hotel_id, room_number, room_type, price_per_night, description, status, created_at, updated_at FROM rooms WHERE hotel_id = $1 AND status = $2`

	var rooms []model.Room

	rows, err := config.DB.Query(query, hotelID, "available")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve rooms"})
	}
	defer rows.Close()

	for rows.Next() {
		var room model.Room
		if err := rows.Scan(
			&room.RoomID,
			&room.HotelID,
			&room.RoomNumber,
			&room.RoomType,
			&room.PricePerNight,
			&room.Description,
			&room.Status,
			&room.CreatedAt,
			&room.UpdatedAt,
		); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to scan room data"})
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Error occurred during rooms retrieval"})
	}

	return c.JSON(http.StatusOK, rooms)
}

func GetBookingsByUserID(c echo.Context) error {
	userID := c.Param("user_id")

	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "user_id is required"})
	}

	query := `
		SELECT id, user_id, room_id, checkin_date, checkout_date, total_price, status, created_at, updated_at 
		FROM bookings WHERE user_id = $1
	`

	bookings := []model.Booking{}

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve bookings"})
	}
	defer rows.Close()

	for rows.Next() {
		var booking model.Booking
		if err := rows.Scan(
			&booking.BookingID,
			&booking.UserID,
			&booking.RoomID,
			&booking.CheckinDate,
			&booking.CheckoutDate,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to scan booking data"})
		}
		bookings = append(bookings, booking)
	}

	return c.JSON(http.StatusOK, bookings)
}

func GetBookingByID(c echo.Context) error {
	bookingID := c.Param("booking_id")

	if bookingID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "booking_id is required"})
	}

	query := `
		SELECT id, user_id, room_id, checkin_date, checkout_date, total_price, status, created_at, updated_at 
		FROM bookings WHERE id = $1
	`

	var booking model.Booking

	err := config.DB.QueryRow(query, bookingID).Scan(
		&booking.BookingID,
		&booking.UserID,
		&booking.RoomID,
		&booking.CheckinDate,
		&booking.CheckoutDate,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Booking not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve booking"})
	}

	return c.JSON(http.StatusOK, booking)
}

func UpdateBookingStatusHandler(c echo.Context) error {
	var req dto.UpdateBookingStatusRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	query := `UPDATE bookings SET status = $1 WHERE id = $2`
	_, err := config.DB.Exec(query, req.Status, req.BookingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update booking status"})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Booking status updated successfully"})
}

func UpdateBookingStatus(c echo.Context) error {
	var req dto.UpdateBookingRefundStatusRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	if req.BookingID == nil || req.UserID == nil || req.Status == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Booking ID, User ID, and Status are required"})
	}

	query := `
		UPDATE bookings 
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3 AND checkin_status NOT IN ('checked_in', 'checked_out')
	`
	res, err := config.DB.Exec(query, req.Status, req.BookingID, req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update booking status"})
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Booking not found or cannot refund"})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Booking status updated successfully"})
}

func UpdateCheckinStatus(c echo.Context) error {
	var req dto.UpdateCheckinStatusRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	if req.BookingID == nil || req.CheckinStatus == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Booking ID and Check-in Status are required"})
	}

	if req.CheckinStatus != "checked_in" && req.CheckinStatus != "checked_out" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Check-in Status must be either 'checked_in' or 'checked_out'"})
	}

	var currentStatus, currentCheckinStatus string
	var roomID int
	query := `SELECT status, checkin_status, room_id FROM bookings WHERE id = $1`
	err := config.DB.QueryRow(query, req.BookingID).Scan(&currentStatus, &currentCheckinStatus, &roomID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Booking not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve booking status"})
	}

	if req.CheckinStatus == "checked_in" {
		if currentStatus != "confirmed" || currentCheckinStatus != "not_checked_in" {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Booking must be confirmed and not checked in to proceed with check-in"})
		}
	} else if req.CheckinStatus == "checked_out" {
		if currentCheckinStatus != "checked_in" {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Booking must be checked in to proceed with check-out"})
		}
	}

	updateQuery := `
		UPDATE bookings 
		SET checkin_status = $1, updated_at = NOW()
		WHERE id = $2
	`
	res, err := config.DB.Exec(updateQuery, req.CheckinStatus, req.BookingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update check-in status"})
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Booking not found or no change in status"})
	}

	var roomStatus string
	if req.CheckinStatus == "checked_in" {
		roomStatus = "occupied"
	} else if req.CheckinStatus == "checked_out" {
		roomStatus = "available"
	}

	roomUpdateQuery := `
		UPDATE rooms 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = config.DB.Exec(roomUpdateQuery, roomStatus, roomID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update room status"})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Check-in status updated successfully"})
}

func UpdateBookingStatusRefund(c echo.Context) error {
	var req dto.UpdateBookingStatusRefundRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	if req.BookingID == nil || req.Status == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Booking ID and Status are required"})
	}

	var roomID int
	var currentStatus, checkinStatus string
	getBookingQuery := `SELECT room_id, status, checkin_status FROM bookings WHERE id = $1`
	err := config.DB.QueryRow(getBookingQuery, req.BookingID).Scan(&roomID, &currentStatus, &checkinStatus)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Booking not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve booking"})
	}


	if currentStatus == "refunded" {
		log.Println(currentStatus, "RED?")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Booking has already been refunded",
		})
	}


	if req.Status == "refunded" {
		if currentStatus != "request_refund" || checkinStatus != "not_checked_in" {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Booking must have status 'request_refund' and check-in status 'not_checked_in' to be refunded",
			})
		}
	}


	updateBookingQuery := `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err = config.DB.Exec(updateBookingQuery, req.Status, req.BookingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update booking status"})
	}


	if req.Status == "refunded" {
		updateRoomQuery := `UPDATE rooms SET status = 'available', updated_at = NOW() WHERE id = $1`
		_, err = config.DB.Exec(updateRoomQuery, roomID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update room status"})
		}
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Booking status updated successfully"})
}

func CancelBooking(c echo.Context) error {
	var req dto.CancelBookingRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	if req.BookingID == nil || req.UserID == nil || *req.BookingID == 0 || *req.UserID == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Booking ID and User ID are required and must be non-zero values"})
	}

	var bookingUserID, roomID int
	var bookingStatus string
	checkBookingQuery := `SELECT user_id, room_id, status FROM bookings WHERE id = $1`
	err := config.DB.QueryRow(checkBookingQuery, *req.BookingID).Scan(&bookingUserID, &roomID, &bookingStatus)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Booking not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve booking"})
	}

	if bookingUserID != *req.UserID {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Message: "Forbidden: Not your booking"})
	}

	if bookingStatus == "canceled" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Booking already canceled"})
	}
	if bookingStatus != "pending" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Booking already confirmed, please request refund"})
	}

	updateBookingQuery := `UPDATE bookings SET status = 'canceled', updated_at = NOW() WHERE id = $1`
	_, err = config.DB.Exec(updateBookingQuery, *req.BookingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update booking status"})
	}

	updateRoomQuery := `UPDATE rooms SET status = 'available', updated_at = NOW() WHERE id = $1`
	_, err = config.DB.Exec(updateRoomQuery, roomID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update room status"})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Booking canceled successfully"})
}
























