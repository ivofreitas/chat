package user

import (
	"github.com/ivofreitas/chat/pkg/config"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repository Repository
}

func NewHandler(repository Repository) *Handler {
	return &Handler{repository: repository}
}

// RegisterUser handles user registration.
// @Summary Register a new user
// @Description Creates a new user account with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body object true "User registration request"
// @Success 201 {object} map[string]string "message: user created successfully"
// @Failure 400 {object} map[string]string "error: invalid request"
// @Failure 500 {object} map[string]string "error: failed to hash password | error: failed to create user"
// @Router /users/register [post]
func (h *Handler) RegisterUser(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	err = h.repository.CreateUser(req.Email, string(hashedPassword))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "user created successfully"})
}

// LoginUser handles user authentication.
// @Summary Login a user
// @Description Authenticates a user and returns a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body object true "User login request"
// @Success 200 {object} map[string]string "token: JWT token"
// @Failure 400 {object} map[string]string "error: invalid request"
// @Failure 401 {object} map[string]string "error: invalid credentials"
// @Failure 500 {object} map[string]string "error: failed to generate token"
// @Router /users/login [post]
func (h *Handler) LoginUser(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	user, err := h.repository.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// generateJWT generates a JWT token for authentication.
func generateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetEnv().Security.JWTSecretKey))
}
