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

var BookingServiceURL string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	BookingServiceURL = os.Getenv("Booking_Service_URL")
	if BookingServiceURL == "" {
		log.Fatal("Booking_Service_URL not set in environment")
	}
}

func GetListHotelsHandler(c echo.Context) error {

	resp, err := http.Get(BookingServiceURL + "/hotel")
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

func GetHotelsHandler(c echo.Context) error {
	id := c.Param("id")

	url := fmt.Sprintf("%s/hotel/%s", BookingServiceURL, id)
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

func ListRoomsByHotelIdHandler(c echo.Context) error {
	hotelID := c.QueryParam("hotel_id")
	url := fmt.Sprintf("%s/hotel/room?hotel_id=%s", BookingServiceURL, hotelID)
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

func GetRoomByIDHandler(c echo.Context) error {
	id := c.Param("id")
	url := fmt.Sprintf("%s/room/%s", BookingServiceURL, id)
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

func CreateHotelHandler(c echo.Context) error {
	var req dto.CreateHotelRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process request data"})
	}

	resp, err := http.Post(BookingServiceURL+"/hotel", "application/json", bytes.NewBuffer(jsonData))
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

func CreateRoomHandler(c echo.Context) error {
	var req dto.CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process request data"})
	}

	url := BookingServiceURL + "/room"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
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

func CreateBookingHandler(c echo.Context) error {
	userID, ok := c.Get("id").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Unauthorized"})
	}

	var req dto.CreateBookingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	req.UserID = int(userID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process request data"})
	}

	url := BookingServiceURL + "/booking"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
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

func UpdateCheckinStatusHandler(c echo.Context) error {
	var req dto.UpdateCheckinStatusRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to process request data"})
	}

	url := fmt.Sprintf("%s/booking/checkin-status", BookingServiceURL)
	client := &http.Client{}
	reqToBookingService, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to create request"})
	}
	reqToBookingService.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(reqToBookingService)
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


