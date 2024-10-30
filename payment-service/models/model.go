package models

import "time"

type Payment struct {
	ID            int       `json:"id"`
	BookingID     int       `json:"booking_id"`
	UserID        int       `json:"user_id"`
	PaymentUID    string    `json:"payment_uid"` 
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"` // "pending", "completed", "refunded"
	PaymentDate   time.Time `json:"payment_date"`
	RefundedAt    *time.Time `json:"refunded_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Refund struct {
	ID           int       `json:"id"`
	PaymentID    int       `json:"payment_id"` 
	RefundAmount float64   `json:"refund_amount"`
	RefundStatus string    `json:"refund_status"` // "requested", "completed", "denied"
	RefundDate   time.Time `json:"refund_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
