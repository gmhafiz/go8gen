package usecase

import (
	"context"

	"abc/internal/domain/book"
	"abc/internal/model"
	"abc/internal/resource"
)

type BookUseCase struct {
	bookRepo book.Repository
}

func NewBookUseCase(bookRepo book.Repository) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

func (u *BookUseCase) Create(ctx context.Context, Book *model.Book) (*model.Book, error) {
	panic("implement me")
}

func (u *BookUseCase) All(ctx context.Context) ([]resource.BookDB, error) {
	panic("implement me")
}

func (u *BookUseCase) Read(ctx context.Context) (int64, error) {
	panic("implement me")
}

func (u *BookUseCase) Update(ctx context.Context) (*model.Book, error) {
	panic("implement me")
}

func (u *BookUseCase) Delete(ctx context.Context) (*model.Book, error) {
	panic("implement me")
}