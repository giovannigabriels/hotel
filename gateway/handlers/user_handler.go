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

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var userServiceURL string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	userServiceURL = os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		log.Fatal("USER_SERVICE_URL not set in environment")
	}
}

func Register(c echo.Context) error {
	var req dto.RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process request"})
	}

	resp, err := http.Post(fmt.Sprintf("%s/register", userServiceURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to connect to user-service"})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to read response from user-service"})
	}

	return c.JSONBlob(resp.StatusCode, respBody)
}


func Login(c echo.Context) error {
	var req dto.LoginUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process request"})
	}

	resp, err := http.Post(fmt.Sprintf("%s/login", userServiceURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to connect to user-service"})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to read response from user-service"})
	}

	return c.JSONBlob(resp.StatusCode, respBody)
}

func GetUserByIDHandler(c echo.Context) error {
	idFloat, ok := c.Get("id").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Unauthorized"})
	}

	userID := fmt.Sprintf("%.0f", idFloat)

	resp, err := http.Get(fmt.Sprintf("%s/user/%s", userServiceURL, userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to connect to user-service"})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to read response from user-service"})
	}

	return c.JSONBlob(resp.StatusCode, respBody)
}


func GetListBooking(c echo.Context) error {
	userID, ok := c.Get("id").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Unauthorized"})
	}

	url := fmt.Sprintf("%s/booking/%d", BookingServiceURL, int(userID))
	resp, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to connect to booking service"})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to read response from booking service"})
	}

	return c.JSONBlob(resp.StatusCode, respBody)
}

func GetDetailBooking(c echo.Context) error {
	bookingID := c.Param("booking_id")
	if bookingID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "booking_id is required"})
	}

	url := fmt.Sprintf("%s/booking/detail/%s", BookingServiceURL, bookingID)
	resp, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to connect to booking service"})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to read response from booking service"})
	}

	return c.JSONBlob(resp.StatusCode, respBody)
}
