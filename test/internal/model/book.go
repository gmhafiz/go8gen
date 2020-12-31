package model

import (
	"time"

	"github.com/volatiletech/null/v8"
)

// Book is an object representing the database table.
type Book struct {
	BookID        int64       `boil:"book_id" json:"book_id" toml:"book_id" yaml:"book_id"`
	Title         string      `boil:"title" json:"title" toml:"title" yaml:"title"`
	PublishedDate time.Time   `boil:"published_date" json:"published_date" toml:"published_date" yaml:"published_date"`
	ImageURL      null.String `boil:"image_url" json:"image_url,omitempty" toml:"image_url" yaml:"image_url,omitempty"`
	Description   null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	CreatedAt     null.Time   `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt     null.Time   `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	DeletedAt     null.Time   `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
}

var BookColumns = struct {
	BookID        string
	Title         string
	PublishedDate string
	ImageURL      string
	Description   string
	CreatedAt     string
	UpdatedAt     string
	DeletedAt     string
}{
	BookID:        "book_id",
	Title:         "title",
	PublishedDate: "published_date",
	ImageURL:      "image_url",
	Description:   "description",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}