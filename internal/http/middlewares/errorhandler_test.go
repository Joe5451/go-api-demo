package middlewares

import (
	"encoding/json"
	"errors"
	"go-api-demo/internal/constant"
	"go-api-demo/internal/http/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		setupHandler   func(*gin.Context)
		expectedStatus int
		expectError    bool
		expectedCode   constant.ErrorCode
		expectedMsg    string
	}{
		{
			name: "no error - should not modify response",
			setupHandler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "has error - should return 500",
			setupHandler: func(c *gin.Context) {
				c.Error(errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
			expectedCode:   constant.ErrInternalServerError,
			expectedMsg:    "an unexpected error occurred",
		},
		{
			name: "multiple errors - should handle last error",
			setupHandler: func(c *gin.Context) {
				c.Error(errors.New("first error"))
				c.Error(errors.New("second error"))
				c.Error(errors.New("last error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
			expectedCode:   constant.ErrInternalServerError,
			expectedMsg:    "an unexpected error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Create a router to test the middleware
			r := gin.New()
			r.Use(ErrorHandler())
			r.GET("/test", func(c *gin.Context) {
				tt.setupHandler(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ErrorHandler() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectError {
				var response util.HTTPError
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}

				if response.Code != tt.expectedCode {
					t.Errorf("ErrorHandler() error code = %v, want %v", response.Code, tt.expectedCode)
				}

				if response.Message != tt.expectedMsg {
					t.Errorf("ErrorHandler() error message = %v, want %v", response.Message, tt.expectedMsg)
				}
			}
		})
	}
}
