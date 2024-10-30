package dto


type ErrorResponse struct {
	Message string `json:"message"`
}


type CreateBookingRequest struct {
	UserID        int     `json:"user_id"`
	RoomID        int     `json:"room_id"`
	CheckinDate   string  `json:"checkin_date"`
	CheckoutDate  string  `json:"checkout_date"`
	TotalPrice    float64 `json:"total_price"`
}

type CreateRoomRequest struct {
	HotelID       int     `json:"hotel_id"`
	RoomNumber    string  `json:"room_number"`
	RoomType      string  `json:"room_type"`
	PricePerNight float64 `json:"price_per_night"`
	Description   string  `json:"description"`
	Status        string  `json:"status"`
}

type CreateHotelRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type CreateHotelResponse struct {
	HotelId int `json:"hotel_id"`
	Message string `json:"message"`
}

type CreateRoomResponse struct {
	RoomId int `json:"room_id"`
	Message string `json:"message"`
}


type CreateBookingResponse struct {
	BookingId int `json:"booking_id"`
	Message string `json:"message"`
}

type UpdateBookingStatusRequest struct {
	BookingID int    `json:"booking_id"`
	Status    string `json:"status"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}