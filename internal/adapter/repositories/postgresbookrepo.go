package repositories

import (
	"context"
	"go-api-demo/internal/application/port/out"
	"go-api-demo/internal/domain"

	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type PostgresBookRepo struct {
	db PgxIface
}

var _ out.BookRepository = &PostgresBookRepo{}

func NewPostgresBookRepo(db PgxIface) *PostgresBookRepo {
	return &PostgresBookRepo{db: db}
}

func (r *PostgresBookRepo) CreateBook(book domain.Book) error {
	_, err := r.db.Exec(
		context.Background(),
		"INSERT INTO books (title, author) VALUES ($1, $2)",
		book.Title,
		book.Author,
	)
	return err
}

func (r *PostgresBookRepo) GetBook(id int) (domain.Book, error) {
	var book domain.Book
	err := r.db.QueryRow(
		context.Background(),
		"SELECT id, title, author FROM books WHERE id = $1",
		id,
	).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Book{}, domain.ErrBookNotFound
		}
		return domain.Book{}, err
	}
	return book, nil
}

func (r *PostgresBookRepo) GetBooks(offset, limit int) ([]domain.Book, error) {
	rows, err := r.db.Query(
		context.Background(),
		"SELECT id, title, author FROM books ORDER BY id ASC LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return []domain.Book{}, err
	}
	defer rows.Close()

	books := []domain.Book{}

	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			return []domain.Book{}, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *PostgresBookRepo) UpdateBook(book domain.Book) error {
	cmdTag, err := r.db.Exec(
		context.Background(),
		"UPDATE books SET title = $1, author = $2 WHERE id = $3",
		book.Title,
		book.Author,
		book.ID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}

func (r *PostgresBookRepo) DeleteBook(id int) error {
	cmdTag, err := r.db.Exec(
		context.Background(),
		"DELETE FROM books WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}
