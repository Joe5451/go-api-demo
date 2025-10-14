package application

import (
	"go-api-demo/internal/application/port/in"
	"go-api-demo/internal/application/port/out"
	"go-api-demo/internal/domain"
)

type BookService struct {
	bookRepo out.BookRepository
}

var _ in.BookUseCase = &BookService{}

func NewBookService(bookRepo out.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) CreateBook(book domain.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	return s.bookRepo.CreateBook(book)
}

func (s *BookService) GetBook(id int) (domain.Book, error) {
	return s.bookRepo.GetBook(id)
}

func (s *BookService) GetBooks(page, perPage int) ([]domain.Book, error) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	limit := perPage

	return s.bookRepo.GetBooks(offset, limit)
}

func (s *BookService) UpdateBook(book domain.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	return s.bookRepo.UpdateBook(book)
}

func (s *BookService) DeleteBook(id int) error {
	return s.bookRepo.DeleteBook(id)
}
