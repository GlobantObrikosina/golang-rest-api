package handler

import (
	"context"
	"fmt"
	"github.com/GlobantObrikosina/golang-rest-api/models"
	"github.com/GlobantObrikosina/golang-rest-api/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.BooksManagerService
}

func NewHandler(service *service.BooksManagerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() http.Handler {
	router := chi.NewRouter()
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Route("/books", h.books)
	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, ErrMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, r, ErrNotFound)
}

var bookIDKey = "bookID"

func (h *Handler) books(router chi.Router) {
	router.Get("/", h.GetAllBooks)
	router.Post("/", h.CreateBook)
	router.Route("/{bookID}", func(router chi.Router) {
		router.Use(h.BookContext)
		router.Get("/", h.GetBook)
		router.Put("/", h.UpdateBook)
		router.Delete("/", h.DeleteBookByID)
	})
}

func (h *Handler) BookContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")

		id, err := strconv.Atoi(bookID)
		if err != nil {
			_ = render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid book ID")))
			return
		}
		ctx := context.WithValue(r.Context(), bookIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	filterCondition := r.URL.Query()
	if len(filterCondition) != 0 {
		if _, ok := filterCondition["genre"]; !ok {
			_ = render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid filter condition")))
			return
		}
		genreID, err := strconv.Atoi(filterCondition.Get("genre"))
		if err != nil || genreID < 1 || genreID > 3 {
			_ = render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid filter condition")))
			return
		}
	}
	books, err := h.service.GetAllBooks(filterCondition)
	if err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, books); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
	}
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	if err := render.Bind(r, book); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := h.service.CreateBook(book); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, models.CreateBookResponse{BookID: book.ID}); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	book, err := h.service.GetBookByID(bookID)
	if err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, &book); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func (h *Handler) DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	err := h.service.DeleteBookByID(bookID)
	if err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
	}
	render.NoContent(w, r)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	bookData := models.Book{}
	if err := render.Bind(r, &bookData); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}
	newBookID, err := h.service.UpdateBookByID(bookID, bookData)
	if err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, models.CreateBookResponse{BookID: newBookID}); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
