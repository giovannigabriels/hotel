package dto

import (
	"time"
)

type CreatePaymentRequest struct {
	BookingID     int     `json:"booking_id"`
	UserID        int     `json:"user_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

type CreatePaymentResponse struct {
	PaymentID     int       `json:"payment_id"`
	PaymentUID    string	`json:"payment_uid"`
	BookingID     int       `json:"booking_id"`
	UserID        int       `json:"user_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	PaymentDate   time.Time   `json:"payment_date"`
	Message       string    `json:"message"`
}


type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type UpdateBookingStatusRequest struct {
	BookingID int    `json:"booking_id"`
	Status    string `json:"status"`
}

type CallbackRequest struct {
	PaymentUID string `json:"payment_uid"`
}

type CreateRefundResponse struct {
	RefundID int    `json:"refund_id"`
	Message  string `json:"message"`
}

type CreateRefundRequest struct {
	UserID    *int `json:"user_id"`
	BookingID *int `json:"booking_id"`
}

type UpdateBookingStatusRefundRequest struct {
	UserID    *int    `json:"user_id"`
	BookingID *int    `json:"booking_id"`
	Status    string `json:"status"`
}