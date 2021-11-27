package handler

import (
	"context"
	"fmt"
	"github.com/GlobantObrikosina/golang-rest-api/db"
	"github.com/GlobantObrikosina/golang-rest-api/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
)

var bookIDKey = "bookID"

func books(router chi.Router) {
	log.Printf("books called successfully")
	router.Get("/", GetAllBooks)
	router.Post("/", CreateBook)
	router.Route("/{bookID}", func(router chi.Router) {
		router.Use(BookContext)
		router.Get("/", GetBook)
		router.Put("/", UpdateBook)
		router.Delete("/", DeleteBook)
	})
}

func BookContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")

		if bookID == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("book ID is required")))
			return
		}
		id, err := strconv.Atoi(bookID)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid book ID")))
		}
		ctx := context.WithValue(r.Context(), bookIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	if err := render.Bind(r, book); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := dbInstance.CreateBook(book); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, models.CreateBookResponse{BookID: book.ID}); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
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
	books, err := dbInstance.GetAllBooks(filterCondition)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, books); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	book, err := dbInstance.GetBookByID(bookID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &book); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	err := dbInstance.DeleteBookByID(bookID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
	}
	render.NoContent(w, r)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	bookData := models.Book{}
	if err := render.Bind(r, &bookData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	newBookID, err := dbInstance.UpdateBookByID(bookID, bookData)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, models.CreateBookResponse{BookID: newBookID}); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
