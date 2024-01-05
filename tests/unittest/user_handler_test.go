package unittest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"notes/internal/model"
	"notes/internal/server"
)

// helper function to create a new server and echo instance
func setupServer(config server.Config) (*server.Server, *echo.Echo) {
	mockRepo := NewMockRepository()
	s := &server.Server{
		Repository: mockRepo,
		Config:     config,
	}
	e := s.RegisterRoutes()
	e.Logger.SetOutput(io.Discard)
	return s, e
}

// TestCreateUser tests the CreateUser handler
func TestCreateUser(t *testing.T) {
	_, e := setupServer(server.NewConfig("", 8080, 10, "secret"))

	// Test cases
	tests := []struct {
		name           string
		payload        model.UserCreateDTO
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Successful User Creation",
			payload: model.UserCreateDTO{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "Secure@Passwprd123",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "Invalid Email",
			payload: model.UserCreateDTO{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "Secure@Passwprd123",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Invalid Password",
			payload: model.UserCreateDTO{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "", // invalid password
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Binding Failure",
			payload:        model.UserCreateDTO{}, // empty payload
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectError {
				assert.NotEqual(t, "", rec.Body.String())
			}
		})
	}
}

func TestRateLimit(t *testing.T) {
	// set rate limit to 1
	_, e := setupServer(server.NewConfig("", 8080, 1, "secret"))

	model := model.UserCreateDTO{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Secure@Passwprd123",
	}
	jsonBody, _ := json.Marshal(model)

	// First request should be OK
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	// Second request should be rate limited
	req = httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(jsonBody))
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}
