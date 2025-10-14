package in

import "go-api-demo/internal/domain"

type BookUseCase interface {
	CreateBook(book domain.Book) error
	GetBook(id int) (domain.Book, error)
	GetBooks(page, perPage int) ([]domain.Book, error)
	UpdateBook(book domain.Book) error
	DeleteBook(id int) error
}
