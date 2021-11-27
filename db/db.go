package db

import (
	"database/sql"
	"fmt"
	"github.com/GlobantObrikosina/golang-rest-api/models"
	_ "github.com/lib/pq"
	"log"
)

const (
	HOST = "database"
	PORT = 5432
)

// ErrNoMatch is returned when we request a row that doesn't exist
var ErrNoMatch = fmt.Errorf("no matching record")

type DatabaseBooksManager interface {
	GetAllBooks(filterCondition map[string][]string) (*models.BookList, error)
	CreateBook(book *models.Book) error
	GetBookByID(bookId int) (models.Book, error)
	DeleteBookByID(bookId int) error
	UpdateBookByID(bookId int, bookData models.Book) (int, error)
	Close() error
}

type Database struct {
	Conn *sql.DB
}

func Initialize(username, password, database string) (DatabaseBooksManager, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	return db, nil
}
