package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"abc/internal/domain/book"
	"abc/internal/model"
	"abc/internal/resource"
)

type repository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) book.Repository {
	return &repository{db: db}
}

func (r *repository) All(ctx context.Context) ([]resource.BookDB, error) {
	panic("implement me")
}

func (r *repository) Read(ctx context.Context) (int64, error) {
	panic("implement me")
}

func (r *repository) Update(ctx context.Context) (*model.Book, error) {
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context) (*model.Book, error) {
	panic("implement me")
}

func (r *repository) Create(ctx context.Context, book *model.Book) (*model.Book, error) {
	panic("implement me")
}

// Close attaches the provider and close the connection
func (r *repository) Close() {
	r.db.Close()
}

// Up attaches the provider and create the table
func (r *repository) Up() error {
	ctx := context.Background()

	query := "CREATE table books(book_id bigserial, title varchar(255) not null, published_date timestamp with time zone not null, image_url varchar(255), description text not null, created_at timestamp with time zone default current_timestamp, updated_at timestamp with time zone default current_timestamp, deleted_at timestamp with time zone, primary key (book_id))"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	return err
}

// Drop attaches the provider and drop the table
func (r *repository) Drop() error {
	ctx := context.Background()

	query := "DROP TABLE IF EXISTS books cascade"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	return err
}
