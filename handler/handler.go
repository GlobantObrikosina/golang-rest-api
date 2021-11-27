package handler

import (
	"github.com/GlobantObrikosina/golang-rest-api/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

var dbInstance db.DatabaseBooksManager

func NewHandler(db db.DatabaseBooksManager) http.Handler {
	router := chi.NewRouter()
	dbInstance = db
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Route("/books", books)
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
