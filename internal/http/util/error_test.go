package util

import (
	"encoding/json"
	"errors"
	"go-api-boilerplate/internal/constant"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		code           constant.ErrorCode
		err            error
		expectedStatus int
		expectedCode   constant.ErrorCode
		expectedMsg    string
	}{
		{
			name:           "validation error",
			status:         http.StatusBadRequest,
			code:           constant.ErrValidationCode,
			err:            errors.New("invalid input"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   constant.ErrValidationCode,
			expectedMsg:    "invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			NewError(c, tt.status, tt.code, tt.err)

			if w.Code != tt.expectedStatus {
				t.Errorf("NewError() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			var response HTTPError
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Code != tt.expectedCode {
				t.Errorf("NewError() code = %v, want %v", response.Code, tt.expectedCode)
			}

			if response.Message != tt.expectedMsg {
				t.Errorf("NewError() message = %v, want %v", response.Message, tt.expectedMsg)
			}
		})
	}
}
