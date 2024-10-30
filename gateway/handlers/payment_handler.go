package handler

import (
	"api-gateway/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)


var PaymentServiceURL string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	PaymentServiceURL = os.Getenv("PAYMENT_SERVICE_URL")
	if PaymentServiceURL == "" {
		log.Fatal("PAYMENT_SERVICE_URL not set in environment")
	}
}

func CreatePaymentHandler(c echo.Context) error {
	userID, ok := c.Get("id").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Unauthorized"})
	}

	var req dto.CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	req.UserID = int(userID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process request data"})
	}

	url := fmt.Sprintf("%s/payment", PaymentServiceURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to connect to payment service"})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to read response from payment service"})
	}

	return c.JSONBlob(resp.StatusCode, respBody)
}

func CreateRefundHandler(c echo.Context) error {
	userID, ok := c.Get("id").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Unauthorized"})
	}

	bookingIDParam := c.Param("booking_id")
	bookingID, err := strconv.Atoi(bookingIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid booking_id"})
	}

	refundRequest := dto.CreateRefundRequest{
		UserID:    int(userID),
		BookingID: bookingID,
	}

	jsonData, err := json.Marshal(refundRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process refund request"})
	}

	url := fmt.Sprintf("%s/refund", PaymentServiceURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to connect to payment service"})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to read response from payment service"})
	}

	return c.JSONBlob(resp.StatusCode, respBody)
}
