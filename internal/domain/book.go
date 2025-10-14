package domain

import (
	"errors"
)

var (
	ErrTitleRequired  = errors.New("title is required")
	ErrAuthorRequired = errors.New("author is required")
	ErrBookNotFound   = errors.New("book not found")
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return ErrTitleRequired
	}
	if b.Author == "" {
		return ErrAuthorRequired
	}
	return nil
}
