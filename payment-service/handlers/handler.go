package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"payment-service/config"
	"payment-service/dto"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)



func CreatePayment(c echo.Context) error {
	var req dto.CreatePaymentRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	if req.BookingID == 0 || req.UserID == 0 || req.Amount <= 0 || req.PaymentMethod == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Missing or invalid payment details"})
	}

	paymentUID := uuid.New().String()

	log.Println(paymentUID, "paymentUID")

	query := `
		INSERT INTO payments (payment_uid, booking_id, user_id, amount, payment_method, payment_status, payment_date, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, payment_date
	`

	var paymentID int
	var paymentDate time.Time
	err := config.DB.QueryRow(query, paymentUID, req.BookingID, req.UserID, req.Amount, req.PaymentMethod, "pending").
		Scan(&paymentID, &paymentDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to create payment"})
	}

	return c.JSON(http.StatusCreated, dto.CreatePaymentResponse{
		PaymentID:     paymentID,
		PaymentUID:    paymentUID,
		BookingID:     req.BookingID,
		UserID:        req.UserID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: "pending",
		PaymentDate:   paymentDate,
		Message:       "Payment created successfully",
	})
}

func PaymentCallbackHandler(c echo.Context) error {
	var req dto.CallbackRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	query := `SELECT id, booking_id FROM payments WHERE payment_uid = $1`
	var paymentID, bookingID int
	err := config.DB.QueryRow(query, req.PaymentUID).Scan(&paymentID, &bookingID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Payment not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve payment"})
	}

	updateQuery := `UPDATE payments SET payment_status = 'success' WHERE id = $1`
	_, err = config.DB.Exec(updateQuery, paymentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update payment status"})
	}


	bookingServiceURL := "http://localhost:5001"
	client := &http.Client{}
	reqBody := dto.UpdateBookingStatusRequest{
		BookingID: bookingID,
		Status:     "confirmed",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to prepare booking update request"})
	}

	resp, err := client.Post(bookingServiceURL+"/booking/status", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update booking status"})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Payment status updated successfully"})
}
