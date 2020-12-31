package http

import (
	"github.com/go-chi/chi"

	"abc/internal/domain/book"
	"abc/internal/middleware"
)

func RegisterHTTPEndPoints(router *chi.Mux, uc book.UseCase) {
	h := NewHandler(uc)

	router.Route("/api/v1/books", func(router chi.Router) {
		router.With(middleware.Paginate).Get("/", h.All)
		router.Post("/", h.Create)
	})
}
