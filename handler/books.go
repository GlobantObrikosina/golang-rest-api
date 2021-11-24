package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/idoko/bucketeer/db"
	"gitlab.com/idoko/bucketeer/models"
	"log"
	"net/http"
	"strconv"
)

var bookIDKey = "bookID"

func books(router chi.Router) {
	log.Printf("books called successfully")
	router.Get("/", getAllBooks)
	router.Post("/", createBook)
	router.Route("/{bookID}", func(router chi.Router) {
		router.Use(BookContext)
		router.Get("/", getBook)
		router.Put("/", updateBook)
		router.Delete("/", deleteBook)
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

func createBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	log.Printf("createBook called %s", r)
	if err := render.Bind(r, book); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	log.Printf("createBook called 1")
	if err := dbInstance.AddBook(book); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	log.Printf("createBook called 2")
	if err := render.Render(w, r, book); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := dbInstance.GetAllBooks()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, books); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	book, err := dbInstance.GetBookById(bookID)
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

func deleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	err := dbInstance.DeleteBook(bookID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.Context().Value(bookIDKey).(int)
	bookData := models.Book{}
	if err := render.Bind(r, &bookData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	book, err := dbInstance.UpdateBook(bookID, bookData)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &book); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
