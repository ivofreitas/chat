package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ivofreitas/chat/internal/auth/application/user/mocks"
	"github.com/ivofreitas/chat/internal/auth/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test cases for user registration
func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		name         string
		email        string
		password     string
		expectedCode int
		mockSetup    func(*mocks.Repository)
	}{
		{
			name:         "Success",
			email:        "test@example.com",
			password:     "securepassword",
			expectedCode: http.StatusCreated,
			mockSetup: func(m *mocks.Repository) {
				m.On("CreateUser", "test@example.com", mock.AnythingOfType("string")).Return(nil)
			},
		},
		{
			name:         "Repository Error",
			email:        "test@example.com",
			password:     "securepassword",
			expectedCode: http.StatusInternalServerError,
			mockSetup: func(m *mocks.Repository) {
				m.On("CreateUser", "test@example.com", mock.AnythingOfType("string")).Return(errors.New("database error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			if tc.mockSetup != nil {
				tc.mockSetup(mockRepo)
			}
			h := NewHandler(mockRepo)
			e := echo.New()
			requestBody, _ := json.Marshal(map[string]string{"email": tc.email, "password": tc.password})
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h.RegisterUser(c)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

// Test cases for user authentication
func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name         string
		email        string
		password     string
		expectedCode int
		mockSetup    func(*mocks.Repository)
	}{
		{
			name:         "Success",
			email:        "test@example.com",
			password:     "securepassword",
			expectedCode: http.StatusOK,
			mockSetup: func(m *mocks.Repository) {
				m.On("GetUserByEmail", "test@example.com").Return(&domain.User{Email: "test@example.com", HashedPassword: "$2b$12$nNun97VbMYsPoGriI0FEpelZu88zKTPeOdX5aHjwEtXrggTQ.MHc2"}, nil)
			},
		},
		{
			name:         "User Not Found",
			email:        "unknown@example.com",
			password:     "securepassword",
			expectedCode: http.StatusUnauthorized,
			mockSetup: func(m *mocks.Repository) {
				m.On("GetUserByEmail", "unknown@example.com").Return(nil, errors.New("user not found"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			if tc.mockSetup != nil {
				tc.mockSetup(mockRepo)
			}
			h := NewHandler(mockRepo)
			e := echo.New()
			requestBody, _ := json.Marshal(map[string]string{"email": tc.email, "password": tc.password})
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h.LoginUser(c)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
