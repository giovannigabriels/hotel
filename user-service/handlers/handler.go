package handler

import (
	"database/sql"
	"net/http"
	"regexp"
	"time"
	"user-service/config"
	"user-service/dto"
	"user-service/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isUserExists(db *sql.DB, username, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 OR email = $2)`
	err := db.QueryRow(query, username, email).Scan(&exists)
	return exists, err
}

func RegisterUser(c echo.Context) error {
	var req dto.RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	if req.Username == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Username is required"})
	}

	if req.Email == "" || !isValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid email format"})
	}

	if req.Password == "" || len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Password must be at least 6 characters"})
	}

	if req.Role != "user" && req.Role != "admin" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Role must be either 'user' or 'admin"})
	}

	exists, err := isUserExists(config.DB, req.Username, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to check user existence"})
	}
	if exists {
		return c.JSON(http.StatusConflict, dto.ErrorResponse{Message: "Username or email already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to hash password"})
	}

	newUser := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO users (username, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err = config.DB.QueryRow(query, newUser.Username, newUser.Email, newUser.Password, newUser.Role, newUser.CreatedAt, newUser.UpdatedAt).Scan(&newUser.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to register user"})
	}

	return c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		UserID:   newUser.UserID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
		Message:  "User registered successfully",
	})
}

var jwtSecret = []byte("rahasia")

func createJWT(userID int, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func LoginUser(c echo.Context) error {
	var req dto.LoginUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}

	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Email is required"})
	}
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Password is required"})
	}

	var user models.User
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	err := config.DB.QueryRow(query, req.Email).Scan(&user.UserID, &user.Email, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Invalid email or password"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to query user"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Invalid email or password"})
	}

	token, err := createJWT(user.UserID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to create token"})
	}

	return c.JSON(http.StatusOK, dto.LoginUserResponse{
		Token: token,
	})
}

func GetUserByID(c echo.Context) error {
	userID := c.Param("id")

	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "User ID is required"})
	}

	query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`

	var user models.User

	err := config.DB.QueryRow(query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "User not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to retrieve user"})
	}

	return c.JSON(http.StatusOK, user)
}

