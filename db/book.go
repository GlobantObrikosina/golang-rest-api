package db

import (
	"database/sql"
	"github.com/GlobantObrikosina/golang-rest-api/models"
)

func (db Database) GetAllBooks() (*models.BookList, error) {
	list := &models.BookList{}
	rows, err := db.Conn.Query("SELECT * FROM books ORDER BY ID DESC")
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

func (db Database) AddBook(book *models.Book) error {
	var id int
	query := `INSERT INTO books (name, genre, price, amount) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.Conn.QueryRow(query, book.Name, book.Genre, book.Price, book.Amount).Scan(&id)
	if err != nil {
		return err
	}
	book.ID = id
	return nil
}

func (db Database) GetBookById(bookId int) (models.Book, error) {
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

func (db Database) DeleteBook(bookId int) error {
	query := `DELETE FROM books WHERE id = $1;`
	_, err := db.Conn.Exec(query, bookId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateBook(bookId int, bookData models.Book) (models.Book, error) {
	book := models.Book{}
	query := `UPDATE books SET name=$1, genre=$2, price=$3, amount=$4 WHERE id=$3 RETURNING id, name, genre, price, amount;`
	err := db.Conn.QueryRow(
		query, bookData.Name, bookData.Genre, bookData.Price, bookData.Amount, bookId).Scan(
		&book.ID, &book.Name, &book.Genre, &book.Price, &book.Amount)

	if err != nil {
		if err == sql.ErrNoRows {
			return book, ErrNoMatch
		}
		return book, err
	}
	return book, nil
}
