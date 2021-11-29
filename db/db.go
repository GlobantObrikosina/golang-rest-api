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

func NewDatabase(username, password, database string) Database {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
		return db
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
		return db
	}
	log.Println("Database connection established")
	return db
}

func (db Database) GetAllBooks(booksFilter map[string][]string) (*models.BookList, error) {
	list := &models.BookList{}
	var query string
	var err error
	rows := &sql.Rows{}

	_, okGenre := booksFilter["genre"]
	_, okName := booksFilter["name"]
	if len(booksFilter) == 0 {
		query = "SELECT * FROM books WHERE amount > 0 ORDER BY ID DESC"
		rows, err = db.Conn.Query(query)
	} else if okGenre && okName {
		query = "SELECT * FROM books WHERE amount > 0 AND genre = $1 AND name = $2 ORDER BY ID DESC"
		rows, err = db.Conn.Query(query, booksFilter["genre"][0], booksFilter["name"][0])
	} else if okGenre {
		query = "SELECT * FROM books WHERE amount > 0 AND genre = $1 ORDER BY ID DESC"
		rows, err = db.Conn.Query(query, booksFilter["genre"][0])
	} else if okName {
		query = "SELECT * FROM books WHERE amount > 0 AND name = $1 ORDER BY ID DESC"
		rows, err = db.Conn.Query(query, booksFilter["name"][0])
	}

	if err != nil {
		return list, err
	}
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Name, &book.Genre, &book.Price, &book.Amount)
		if err != nil {
			return list, err
		}
		list.Books = append(list.Books, book)
	}
	return list, nil
}

func (db Database) CreateBook(book *models.Book) error {
	var id int
	query := `INSERT INTO books (name, genre, price, amount) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.Conn.QueryRow(query, book.Name, book.Genre, book.Price, book.Amount).Scan(&id)
	if err != nil {
		book.ID = 0
		return err
	}
	book.ID = id
	return nil
}

func (db Database) GetBookByID(bookId int) (models.Book, error) {
	book := models.Book{}
	query := `SELECT * FROM books WHERE id = $1;`
	row := db.Conn.QueryRow(query, bookId)
	switch err := row.Scan(&book.ID, &book.Name, &book.Genre, &book.Price, &book.Amount); err {
	case sql.ErrNoRows:
		return book, ErrNoMatch
	default:
		return book, err
	}
}

func (db Database) DeleteBookByID(bookId int) error {
	query := `DELETE FROM books WHERE id = $1;`
	_, err := db.Conn.Exec(query, bookId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateBookByID(bookId int, bookData models.Book) (int, error) {
	query := `UPDATE books SET name=$1, genre=$2, price=$3, amount=$4 WHERE id=$5 RETURNING id;`
	var newBookID int
	err := db.Conn.QueryRow(
		query, bookData.Name, bookData.Genre, bookData.Price, bookData.Amount, bookId).Scan(&newBookID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrNoMatch
		}
		return 0, err
	}
	return newBookID, nil
}

func (db Database) Close() error {
	return db.Conn.Close()
}
