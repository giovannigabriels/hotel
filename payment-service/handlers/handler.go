package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	bookingServiceURL := os.Getenv("BOOKING_SERVICE_URL")
	if bookingServiceURL == "" {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Booking service URL is not configured"})
	}

	reqBody := dto.UpdateBookingStatusRequest{
		BookingID: bookingID,
		Status:    "confirmed",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to prepare booking update request"})
	}

	resp, err := http.Post(bookingServiceURL+"/booking/callback/status", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to update booking status"})
	}
	defer resp.Body.Close()

	return c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Payment status updated successfully"})
}

func CreateRefund(c echo.Context) error {
	var req dto.CreateRefundRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	if req.UserID == nil || req.BookingID == nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "User ID and Booking ID are required"})
	}

	bookingServiceURL := os.Getenv("BOOKING_SERVICE_URL")
	if bookingServiceURL == "" {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Booking service URL is not configured"})
	}

	updateStatusRequest := dto.UpdateBookingStatusRefundRequest{
		UserID:    req.UserID,
		BookingID: req.BookingID,
		Status:    "request_refund",
	}

	jsonUpdateStatusData, err := json.Marshal(updateStatusRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process status update request"})
	}

	url := fmt.Sprintf("%s/booking/refund/status", bookingServiceURL)
	updateStatusResp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonUpdateStatusData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to connect to booking service"})
	}
	defer updateStatusResp.Body.Close()

	if updateStatusResp.StatusCode != http.StatusOK {
		return c.JSON(updateStatusResp.StatusCode, dto.ErrorResponse{Message: "Failed to update booking status"})
	}

	var paymentID int
	var paymentAmount float64
	checkPaymentQuery := `
		SELECT id, amount FROM payments WHERE user_id = $1 AND booking_id = $2 AND payment_status = 'success'
	`
	err = config.DB.QueryRow(checkPaymentQuery, req.UserID, req.BookingID).Scan(&paymentID, &paymentAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "No completed payment found for the provided booking"})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to check payment"})
	}

	var existingRefundID int
	checkRefundQuery := `SELECT id FROM refunds WHERE payment_id = $1`
	err = config.DB.QueryRow(checkRefundQuery, paymentID).Scan(&existingRefundID)
	if err == nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Refund request already exists for this booking"})
	} else if err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to check existing refund"})
	}

	query := `
		INSERT INTO refunds (payment_id, refund_amount, refund_status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	var refundID int
	err = config.DB.QueryRow(query, paymentID, paymentAmount, "requested").Scan(&refundID)
	if err != nil {
		log.Println(err, "error")
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to create refund"})
	}

	return c.JSON(http.StatusCreated, dto.CreateRefundResponse{
		RefundID: refundID,
		Message:  "Refund requested successfully",
	})
}

