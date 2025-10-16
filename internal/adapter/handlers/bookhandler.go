package handlers

import (
	"errors"
	"fmt"
	"go-api-demo/internal/application/port/in"
	"go-api-demo/internal/constant"
	"go-api-demo/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	CreateBookReq struct {
		Title  string `json:"title" binding:"required" example:"The Great Gatsby"`
		Author string `json:"author" binding:"required" example:"John Doe"`
	}
	GetBooksReq struct {
		Page    int `form:"page,default=1" binding:"min=1" example:"1"`
		PerPage int `form:"per_page,default=10" binding:"min=1,max=100" example:"10"`
	}
	UpdateBookReq struct {
		Title  string `json:"title" binding:"required" example:"The Great Gatsby"`
		Author string `json:"author" binding:"required" example:"John Doe"`
	}
)

type BookHandler struct {
	bookService in.BookUseCase
}

func NewBookHandler(bookService in.BookUseCase) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var json CreateBookReq
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(fmt.Errorf("%w: %v", constant.ErrValidation, err))
		return
	}

	book := domain.Book{
		Title:  json.Title,
		Author: json.Author,
	}

	err := h.bookService.CreateBook(book)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", constant.ErrValidation, "ID should be a number"))
		return
	}

	book, err := h.bookService.GetBook(ID)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    constant.ErrNotFoundCode,
				"message": "Book not found",
			})
			return
		}
		c.Error(err)
		return
	}

	res := struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Author string `json:"author"`
	}{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
	}
	c.JSON(http.StatusOK, res)
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	var query GetBooksReq
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(fmt.Errorf("%w: %v", constant.ErrValidation, err))
		return
	}

	books, err := h.bookService.GetBooks(query.Page, query.PerPage)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", constant.ErrValidation, "ID should be a number"))
		return
	}

	var json UpdateBookReq
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(fmt.Errorf("%w: %v", constant.ErrValidation, err))
		return
	}

	err = h.bookService.UpdateBook(domain.Book{
		ID:     ID,
		Title:  json.Title,
		Author: json.Author,
	})
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    constant.ErrNotFoundCode,
				"message": "Book not found",
			})
			return
		}
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", constant.ErrValidation, "ID should be a number"))
		return
	}

	err = h.bookService.DeleteBook(ID)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    constant.ErrNotFoundCode,
				"message": "Book not found",
			})
			return
		}
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
