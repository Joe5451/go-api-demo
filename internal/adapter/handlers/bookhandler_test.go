package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-api-demo/internal/domain"
	"go-api-demo/internal/http/middlewares"
	"go-api-demo/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter() *gin.Engine {
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	return r
}

func TestBookHandler_CreateBook(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setup      func(*mocks.MockBookUseCase)
		wantStatus int
	}{
		{
			name: "success",
			body: CreateBookReq{Title: "Test Book", Author: "Test Author"},
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().CreateBook(gomock.Any(), domain.Book{Title: "Test Book", Author: "Test Author"}).Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid json",
			body:       `invalid`,
			setup:      func(m *mocks.MockBookUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation error - missing title",
			body:       CreateBookReq{Author: "Test Author"},
			setup:      func(m *mocks.MockBookUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: CreateBookReq{Title: "Test Book", Author: "Test Author"},
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().CreateBook(gomock.Any(), gomock.Any()).Return(errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockBookUseCase(ctrl)
			tt.setup(mockService)

			h := NewBookHandler(mockService)

			r := setupTestRouter()
			r.POST("/books", h.CreateBook)

			w := httptest.NewRecorder()
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateBook() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestBookHandler_GetBook(t *testing.T) {
	tests := []struct {
		name       string
		bookID     string
		setup      func(*mocks.MockBookUseCase)
		wantStatus int
	}{
		{
			name:   "success",
			bookID: "1",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().GetBook(gomock.Any(), 1).Return(domain.Book{ID: 1, Title: "Test Book", Author: "Test Author"}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			bookID:     "invalid",
			setup:      func(m *mocks.MockBookUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "book not found",
			bookID: "999",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().GetBook(gomock.Any(), 999).Return(domain.Book{}, domain.ErrBookNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "service error",
			bookID: "1",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().GetBook(gomock.Any(), 1).Return(domain.Book{}, errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockBookUseCase(ctrl)
			tt.setup(mockService)

			h := NewBookHandler(mockService)

			r := setupTestRouter()
			r.GET("/books/:id", h.GetBook)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/books/"+tt.bookID, nil)

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetBook() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestBookHandler_GetBooks(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		setup      func(*mocks.MockBookUseCase)
		wantStatus int
	}{
		{
			name:  "success with defaults",
			query: "",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().GetBooks(gomock.Any(), 1, 10).Return([]domain.Book{
					{ID: 1, Title: "Book 1", Author: "Author 1"},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "success with pagination",
			query: "?page=2&per_page=20",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().GetBooks(gomock.Any(), 2, 20).Return([]domain.Book{}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid page",
			query:      "?page=0",
			setup:      func(m *mocks.MockBookUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "service error",
			query: "",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().GetBooks(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockBookUseCase(ctrl)
			tt.setup(mockService)

			h := NewBookHandler(mockService)

			r := setupTestRouter()
			r.GET("/books", h.GetBooks)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/books"+tt.query, nil)

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetBooks() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestBookHandler_UpdateBook(t *testing.T) {
	tests := []struct {
		name       string
		bookID     string
		body       interface{}
		setup      func(*mocks.MockBookUseCase)
		wantStatus int
	}{
		{
			name:   "success",
			bookID: "1",
			body:   UpdateBookReq{Title: "Updated Book", Author: "Updated Author"},
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().UpdateBook(gomock.Any(), domain.Book{ID: 1, Title: "Updated Book", Author: "Updated Author"}).Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id",
			bookID:     "invalid",
			body:       UpdateBookReq{Title: "Updated Book", Author: "Updated Author"},
			setup:      func(m *mocks.MockBookUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation error - missing title",
			bookID:     "1",
			body:       UpdateBookReq{Author: "Updated Author"},
			setup:      func(m *mocks.MockBookUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "book not found",
			bookID: "999",
			body:   UpdateBookReq{Title: "Updated Book", Author: "Updated Author"},
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().UpdateBook(gomock.Any(), gomock.Any()).Return(domain.ErrBookNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "service error",
			bookID: "1",
			body:   UpdateBookReq{Title: "Updated Book", Author: "Updated Author"},
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().UpdateBook(gomock.Any(), gomock.Any()).Return(errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockBookUseCase(ctrl)
			tt.setup(mockService)

			h := NewBookHandler(mockService)

			r := setupTestRouter()
			r.PUT("/books/:id", h.UpdateBook)

			w := httptest.NewRecorder()
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/books/"+tt.bookID, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("UpdateBook() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestBookHandler_DeleteBook(t *testing.T) {
	tests := []struct {
		name       string
		bookID     string
		setup      func(*mocks.MockBookUseCase)
		wantStatus int
	}{
		{
			name:   "success",
			bookID: "1",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().DeleteBook(gomock.Any(), 1).Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id",
			bookID:     "invalid",
			setup:      func(m *mocks.MockBookUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "book not found",
			bookID: "999",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().DeleteBook(gomock.Any(), 999).Return(domain.ErrBookNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "service error",
			bookID: "1",
			setup: func(m *mocks.MockBookUseCase) {
				m.EXPECT().DeleteBook(gomock.Any(), 1).Return(errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockBookUseCase(ctrl)
			tt.setup(mockService)

			h := NewBookHandler(mockService)

			r := setupTestRouter()
			r.DELETE("/books/:id", h.DeleteBook)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/books/"+tt.bookID, nil)

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("DeleteBook() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
