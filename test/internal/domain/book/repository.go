package book

import (
	"context"

	"abc/internal/model"
	"abc/internal/resource"
)

type Repository interface {
	Create(ctx context.Context, Book *model.Book) (*model.Book, error)
	All(ctx context.Context) ([]resource.BookDB, error)
	Read(ctx context.Context) (int64, error)
	Update(ctx context.Context) (*model.Book, error)
	Delete(ctx context.Context) (*model.Book, error)
	Close()
	Drop() error
	Up() error
}