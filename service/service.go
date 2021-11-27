package service

import (
	"github.com/GlobantObrikosina/golang-rest-api/db"
	"github.com/GlobantObrikosina/golang-rest-api/models"
)

type DatabaseBooksManager interface {
	GetAllBooks(filterCondition map[string][]string) (*models.BookList, error)
	CreateBook(book *models.Book) error
	GetBookByID(bookId int) (models.Book, error)
	DeleteBookByID(bookId int) error
	UpdateBookByID(bookId int, bookData models.Book) (int, error)
	Close() error
}

type BooksManagerService struct {
	repo db.DatabaseBooksManager
}

func NewService(repo db.DatabaseBooksManager) *BooksManagerService {
	return &BooksManagerService{repo: repo}
}

func (s *BooksManagerService) CreateBook(book *models.Book) error {
	return s.repo.CreateBook(book)
}

func (s *BooksManagerService) GetBookByID(id int) (models.Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *BooksManagerService) GetAllBooks(filterCondition map[string][]string) (*models.BookList, error) {
	return s.repo.GetAllBooks(filterCondition)
}

func (s *BooksManagerService) DeleteBookByID(id int) error {
	return s.repo.DeleteBookByID(id)
}

func (s *BooksManagerService) UpdateBookByID(id int, book models.Book) (int, error) {
	return s.repo.UpdateBookByID(id, book)
}
