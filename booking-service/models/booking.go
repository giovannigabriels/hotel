package model

type RoomStatus string

const (
    Available    RoomStatus = "available"
    Booked       RoomStatus = "booked"
    Occupied     RoomStatus = "occupied"
    Maintenance  RoomStatus = "maintenance"
)

type BookingStatus string

const (
    Pending   BookingStatus = "pending"
    Confirmed BookingStatus = "confirmed"
    Canceled  BookingStatus = "canceled"
    Refund    BookingStatus = "request_refund"
)

type CheckinStatus string

const (
    NotCheckedIn CheckinStatus = "not_checked_in"
    CheckedIn    CheckinStatus = "checked_in"
    CheckedOut   CheckinStatus = "checked_out"
)

type RoomType string

const (
    Single RoomType = "single"
    Double RoomType = "double"
    Suite  RoomType = "suite"
    Deluxe RoomType = "deluxe"
)


type Hotel struct {
    HotelID     int    `json:"id"`
    Name        string `json:"name"`
    Address     string `json:"address"`
    City        string `json:"city"`
    Country     string `json:"country"`
    PhoneNumber string `json:"phone_number"`
    Email       string `json:"email"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type Room struct {
    RoomID        int        `json:"id"`
    HotelID       int        `json:"hotel_id"`
    RoomNumber    string     `json:"room_number"`
    RoomType      RoomType   `json:"room_type"`
    PricePerNight float64    `json:"price_per_night"`
    Description   string     `json:"description"`
    Status        RoomStatus `json:"status"`
    CreatedAt     string     `json:"created_at"`
    UpdatedAt     string     `json:"updated_at"`
}

type Booking struct {
    BookingID    int           `json:"id"`
    UserID       int           `json:"user_id"`
    RoomID       int           `json:"room_id"`
    CheckinDate  string        `json:"checkin_date"`
    CheckoutDate string        `json:"checkout_date"`
    TotalPrice   float64       `json:"total_price"`
    Status       BookingStatus `json:"status"`
    CheckinStatus CheckinStatus `json:"checkin_status"`
    CreatedAt    string        `json:"created_at"`
    UpdatedAt    string        `json:"updated_at"`
}
