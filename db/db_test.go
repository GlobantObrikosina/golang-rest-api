package db

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GlobantObrikosina/golang-rest-api/models"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func MockDB() (DatabaseBooksManager, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return Database{db}, mock
}

func TestGetAllBooks(t *testing.T) {
	repo, mock := MockDB()

	type mockBehavior func(filterCondition map[string][]string)
	tests := []struct {
		name            string
		mockBehavior    mockBehavior
		filterCondition map[string][]string
		expectedBooks   []models.Book
		expectError     bool
	}{
		{
			name: "Ok",
			mockBehavior: func(filterCondition map[string][]string) {
				mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "genre", "price", "amount"}).
						AddRow(1, "book1", 1, 3.7, 1).
						AddRow(2, "book2", 2, 4.7, 2).
						AddRow(3, "book3", 3, 5.7, 3))
			},
			expectedBooks: []models.Book{
				{ID: 1, Name: "book1", Genre: 1, Price: 3.7, Amount: 1},
				{ID: 2, Name: "book2", Genre: 2, Price: 4.7, Amount: 2},
				{ID: 3, Name: "book3", Genre: 3, Price: 5.7, Amount: 3},
			},
		},
		{
			name:            "Filter",
			filterCondition: map[string][]string{"genre": {"1"}},
			mockBehavior: func(filterCondition map[string][]string) {
				genreId := filterCondition["genre"][0]
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(genreId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "genre", "amount"}).
						AddRow(1, "book1", 1, 3.7, 1))
			},
			expectedBooks: []models.Book{
				{ID: 1, Name: "book1", Genre: 1, Price: 3.7, Amount: 1},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.filterCondition)
			books, err := repo.GetAllBooks(test.filterCondition)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBooks, books.Books)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddBook(t *testing.T) {
	repo, mock := MockDB()
	type mockBehavior func(mock sqlmock.Sqlmock, returnedId int, book models.Book)
	tests := []struct {
		name         string
		inputBook    models.Book
		returnedId   int
		mockBehavior mockBehavior
		expectError  bool
	}{
		{
			name: "OK",
			inputBook: models.Book{
				ID:     1,
				Name:   "hello",
				Price:  45.99,
				Genre:  1,
				Amount: 8,
			},
			returnedId: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, returnedId int, book models.Book) {
				mock.ExpectQuery(`INSERT INTO books`).
					WithArgs(book.Name, book.Genre, book.Price, book.Amount).
					WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(returnedId))
			},
			expectError: false,
		},
		{
			name: "Empty field",
			inputBook: models.Book{
				Name:   "",
				Price:  45.99,
				Genre:  1,
				Amount: 8,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, returnedId int, book models.Book) {
				mock.ExpectQuery(`INSERT INTO books`).
					WithArgs(book.Name, book.Genre, book.Price, book.Amount).
					WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(returnedId).
						RowError(0, errors.New("insert error")))
			},
			expectError: true,
		},
		{
			name: "OK 2, check ID",
			inputBook: models.Book{
				Name:   "hello world",
				Price:  45.99,
				Genre:  1,
				Amount: 8,
			},
			returnedId: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, returnedId int, book models.Book) {
				mock.ExpectQuery(`INSERT INTO books`).
					WithArgs(book.Name, book.Genre, book.Price, book.Amount).
					WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(returnedId))
			},
			expectError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(mock, test.returnedId, test.inputBook)
			err := repo.CreateBook(&test.inputBook)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.returnedId, test.inputBook.ID)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetBookById(t *testing.T) {
	repo, mock := MockDB()
	type mockBehavior func(inputId int)
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		inputId      int
		expectedBook models.Book
		expectError  bool
	}{
		{
			name: "Ok",
			mockBehavior: func(inputId int) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(inputId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "genre", "amount"}).
						AddRow(1, "book1", 2, 1.11, 9))
			},
			inputId: 2,
			expectedBook: models.Book{
				ID:     1,
				Name:   "book1",
				Genre:  2,
				Price:  1.11,
				Amount: 9,
			},
		},
		{
			name: "Id not found",
			mockBehavior: func(inputId int) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(inputId).WillReturnError(errors.New("id not found"))
			},
			inputId:     2,
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.inputId)
			book, err := repo.GetBookByID(test.inputId)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBook, book)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteBook(t *testing.T) {
	repo, mock := MockDB()
	type mockBehavior func(inputId int)
	tests := []struct {
		name         string
		inputId      int
		mockBehavior mockBehavior
		expectError  bool
	}{
		{
			name:    "Ok",
			inputId: 3,
			mockBehavior: func(inputId int) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs(inputId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Id not found",
			mockBehavior: func(inputId int) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE")).
					WithArgs(inputId).WillReturnError(errors.New("id not found"))
			},
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.inputId)
			err := repo.DeleteBookByID(test.inputId)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateBook(t *testing.T) {
	repo, mock := MockDB()
	type mockBehavior func(inputId int, inputBook models.Book)
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		inputId      int
		inputBook    models.Book
		expectError  bool
	}{
		{
			name: "Ok",
			mockBehavior: func(inputId int, inputBook models.Book) {
				mock.ExpectQuery(`UPDATE`).
					WithArgs(inputBook.Name, inputBook.Genre, inputBook.Price, inputBook.Amount, inputId).
					WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(inputId))
			},
			inputId: 1,
			inputBook: models.Book{
				ID:     1,
				Name:   "book1",
				Genre:  2,
				Price:  1.11,
				Amount: 9,
			},
		},
		{
			name: "id not found",
			mockBehavior: func(inputId int, inputBook models.Book) {
				mock.ExpectQuery(`UPDATE`).
					WithArgs(inputBook.Name, inputBook.Genre, inputBook.Price, inputBook.Amount, inputId).
					WillReturnError(errors.New("id not found"))
			},
			inputId: 1,
			inputBook: models.Book{
				ID:     1,
				Name:   "book1",
				Genre:  2,
				Price:  1.11,
				Amount: 9,
			},
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.inputId, test.inputBook)
			updatedBook, err := repo.UpdateBookByID(test.inputId, test.inputBook)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, updatedBook, test.inputBook.ID)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
