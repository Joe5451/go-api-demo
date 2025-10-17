package handlers

import (
	"errors"
	"go-api-demo/internal/application/port/in"
	"go-api-demo/internal/constant"
	"go-api-demo/internal/domain"
	"go-api-demo/internal/http/util"
	"net/http"

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

// CreateBook godoc
// @Summary      Create a book
// @Description  Create a book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        request  body		CreateBookReq	true "Create book"
// @Success      204  {object}	nil
// @Failure      400  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Router       /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var json CreateBookReq
	if err := c.ShouldBindJSON(&json); err != nil {
		util.NewError(c, http.StatusBadRequest, constant.ErrValidationCode, err)
		return
	}

	book := domain.Book{
		Title:  json.Title,
		Author: json.Author,
	}

	err := h.bookService.CreateBook(c.Request.Context(), book)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetBook godoc
// @Summary      Get a book
// @Description  Get a book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Book ID"
// @Success      200  {object}  domain.Book
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Router       /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	type params struct {
		ID int `uri:"id" binding:"required"`
	}
	var p params
	if err := c.ShouldBindUri(&p); err != nil {
		util.NewError(c, http.StatusBadRequest, constant.ErrValidationCode, err)
		return
	}

	book, err := h.bookService.GetBook(c.Request.Context(), p.ID)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			util.NewError(c, http.StatusNotFound, constant.ErrNotFoundCode, domain.ErrBookNotFound)
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

// GetBooks godoc
// @Summary      Get books
// @Description  Get books
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        page  query  int  false  "Page"
// @Param        per_page  query  int  false  "Per Page"
// @Success      200  {object}  []domain.Book
// @Failure      400  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Router       /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	var query GetBooksReq
	if err := c.ShouldBindQuery(&query); err != nil {
		util.NewError(c, http.StatusBadRequest, constant.ErrValidationCode, err)
		return
	}

	books, err := h.bookService.GetBooks(c.Request.Context(), query.Page, query.PerPage)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, books)
}

// UpdateBook godoc
// @Summary      Update a book
// @Description  Update a book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Book ID"
// @Param        request  body  UpdateBookReq  true  "Update book"
// @Success      204  {object}  nil
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Router       /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	type params struct {
		ID int `uri:"id" binding:"required"`
	}
	var p params
	if err := c.ShouldBindUri(&p); err != nil {
		util.NewError(c, http.StatusBadRequest, constant.ErrValidationCode, err)
		return
	}

	var json UpdateBookReq
	if err := c.ShouldBindJSON(&json); err != nil {
		util.NewError(c, http.StatusBadRequest, constant.ErrValidationCode, err)
		return
	}

	err := h.bookService.UpdateBook(c.Request.Context(), domain.Book{
		ID:     p.ID,
		Title:  json.Title,
		Author: json.Author,
	})
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			util.NewError(c, http.StatusNotFound, constant.ErrNotFoundCode, domain.ErrBookNotFound)
			return
		}
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteBook godoc
// @Summary      Delete a book
// @Description  Delete a book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Book ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Router       /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	type params struct {
		ID int `uri:"id" binding:"required"`
	}
	var p params
	if err := c.ShouldBindUri(&p); err != nil {
		util.NewError(c, http.StatusBadRequest, constant.ErrValidationCode, err)
		return
	}

	err := h.bookService.DeleteBook(c.Request.Context(), p.ID)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			util.NewError(c, http.StatusNotFound, constant.ErrNotFoundCode, domain.ErrBookNotFound)
			return
		}
		util.NewError(c, http.StatusInternalServerError, constant.ErrInternalServerError, err)
		return
	}
	c.Status(http.StatusNoContent)
}
