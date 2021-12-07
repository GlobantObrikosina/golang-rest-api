package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/GlobantObrikosina/golang-rest-api/models"
	"github.com/GlobantObrikosina/golang-rest-api/service"
	mock "github.com/GlobantObrikosina/golang-rest-api/service/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetAllBooks(t *testing.T) {
	type mockBehavior func(s *mock.MockDatabaseBooksManager, filterCondition map[string][]string)
	tests := []struct {
		name                 string
		filterCondition      map[string][]string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Invalid filter condition",
			filterCondition:      map[string][]string{"avadakedavra": {"7"}},
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, filterCondition map[string][]string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"invalid filter condition\"}\n",
		},
		{
			name:                 "Invalid genre id in filter condition",
			filterCondition:      map[string][]string{"genre": {"0"}},
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, filterCondition map[string][]string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"invalid filter condition\"}\n",
		},
		{
			name:            "Get All Ok",
			filterCondition: map[string][]string{},
			mockBehavior: func(r *mock.MockDatabaseBooksManager, filterCondition map[string][]string) {
				r.EXPECT().GetAllBooks(filterCondition).Return(&models.BookList{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"books\":null}\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock.NewMockDatabaseBooksManager(c)
			test.mockBehavior(mockManager, test.filterCondition)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := chi.NewRouter()
			r.Get("/books", handler.GetAllBooks)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/books", nil)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			values := url.Values(test.filterCondition)
			req.URL.RawQuery = values.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestCreateBook(t *testing.T) {
	type mockBehavior func(s *mock.MockDatabaseBooksManager, book *models.Book)
	tests := []struct {
		name                 string
		inputBody            string
		inputBook            *models.Book
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Everything's Ok",
			inputBody: `{"name": "Book1", "genre": 1, "price": 1.1, "amount": 1}`,
			inputBook: &models.Book{
				ID:     0,
				Name:   "Book1",
				Genre:  1,
				Price:  1.1,
				Amount: 1,
			},
			mockBehavior: func(r *mock.MockDatabaseBooksManager, book *models.Book) {
				r.EXPECT().CreateBook(book).Return(0, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":0}\n",
		},
		{
			name:                 "Fields missing",
			inputBody:            `{"price": 67.88, "genre": 1, "amount": 5}`,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, book *models.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"name is a required field\"}\n",
		},
		{
			name:                 "Invalid genre",
			inputBody:            `{"name": "hello", "price": 67.88, "genre": 112, "amount": 7}`,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, book *models.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"genre has to be between 1 and 3\"}\n",
		},
		{
			name:                 "Invalid genre",
			inputBody:            `{"name": "hello", "price": -67.88, "genre": 1, "amount": 7}`,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, book *models.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"price can't be less than 0\"}\n",
		},
		{
			name:                 "Invalid genre",
			inputBody:            `{"name": "hello", "price": 67.88, "genre": 1, "amount": -7}`,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, book *models.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"amount can't be less than 0\"}\n",
		},
		{
			name:      "Not unique name",
			inputBody: `{"name": "hello", "price": 67.88, "genre": 1, "amount": 7}`,
			inputBook: &models.Book{
				ID:     0,
				Name:   "hello",
				Price:  67.88,
				Genre:  1,
				Amount: 7,
			},
			mockBehavior: func(r *mock.MockDatabaseBooksManager, book *models.Book) {
				r.EXPECT().CreateBook(book).Return(0, errors.New("input book name isn't unique"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "{\"status_text\":\"Internal server error\",\"message\":\"input book name isn't unique\"}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock.NewMockDatabaseBooksManager(c)
			test.mockBehavior(mockManager, test.inputBook)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := chi.NewRouter()
			r.Post("/books", handler.CreateBook)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/books",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}

}

func TestDeleteBook(t *testing.T) {
	type mockBehavior func(s *mock.MockDatabaseBooksManager, id int)
	tests := []struct {
		name                 string
		inputId              int
		falseId              string
		useFalseID           bool
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Id invalid",
			falseId:              "ababababa",
			useFalseID:           true,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, id int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"invalid book ID\"}\n",
		},
		{
			name:    "Id invalid",
			inputId: 10,
			mockBehavior: func(r *mock.MockDatabaseBooksManager, id int) {
				r.EXPECT().DeleteBookByID(id).Return(errors.New("no matching record"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"no matching record\"}\n",
		},
		{
			name:    "Id OK",
			inputId: 1,
			mockBehavior: func(r *mock.MockDatabaseBooksManager, id int) {
				r.EXPECT().DeleteBookByID(id).Return(nil)
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock.NewMockDatabaseBooksManager(c)
			test.mockBehavior(mockManager, test.inputId)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := chi.NewRouter()
			r.Route("/{bookID}", func(router chi.Router) {
				router.Use(handler.BookContext)
				router.Delete("/", handler.DeleteBookByID)
			})

			var target string
			if test.useFalseID {
				target = fmt.Sprintf("/%v", test.falseId)
			} else {
				target = fmt.Sprintf("/%v", test.inputId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", target, nil)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestGetBook(t *testing.T) {
	type mockBehavior func(s *mock.MockDatabaseBooksManager, id int)
	tests := []struct {
		name                 string
		inputId              int
		falseId              string
		useFalseID           bool
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Id invalid",
			falseId:              "knekndijf",
			useFalseID:           true,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, id int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"invalid book ID\"}\n",
		},
		{
			name:    "Id OK",
			inputId: 1,
			mockBehavior: func(r *mock.MockDatabaseBooksManager, id int) {
				r.EXPECT().GetBookByID(id).Return(models.Book{
					ID:     1,
					Name:   "hello",
					Price:  4.32,
					Genre:  2,
					Amount: 9,
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1,\"name\":\"hello\",\"genre\":2,\"price\":4.32,\"amount\":9}\n",
		},
		{
			name:    "Id not found",
			inputId: 1,
			mockBehavior: func(r *mock.MockDatabaseBooksManager, id int) {
				r.EXPECT().GetBookByID(id).Return(models.Book{}, fmt.Errorf("no matching record"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"no matching record\"}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock.NewMockDatabaseBooksManager(c)
			test.mockBehavior(mockManager, test.inputId)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := chi.NewRouter()
			r.Route("/{bookID}", func(router chi.Router) {
				router.Use(handler.BookContext)
				router.Get("/", handler.GetBook)
			})
			var target string
			if test.useFalseID {
				target = fmt.Sprintf("/%v", test.falseId)
			} else {
				target = fmt.Sprintf("/%v", test.inputId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", target, nil)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestUpdateBook(t *testing.T) {
	type mockBehavior func(s *mock.MockDatabaseBooksManager, id int, book models.Book)
	tests := []struct {
		name                 string
		inputBody            string
		inputId              int
		falseId              string
		useFalseID           bool
		inputBook            models.Book
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Update invalid id",
			falseId:              "78.99",
			useFalseID:           true,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, id int, book models.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"invalid book ID\"}\n",
		},
		{
			name:      "Update invalid price input",
			inputId:   1,
			inputBody: `{"name": "Book1", "genre": 1, "price": -12.12, "amount": 0}`,
			mockBehavior: func(r *mock.MockDatabaseBooksManager, id int, book models.Book) {
				//r.EXPECT().UpdateBookByID(id, book).Return(0, errors.New("price can't be less than 0"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"price can't be less than 0\"}\n",
		},
		{
			name:                 "Update invalid genre input",
			inputId:              1,
			inputBody:            `{"name": "Book1", "genre": 111, "price": 12.12, "amount": 0}`,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, id int, book models.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"genre has to be between 1 and 3\"}\n",
		},
		{
			name:                 "Update invalid amount input",
			inputId:              1,
			inputBody:            `{"name": "Book1", "genre": 1, "price": 12.12, "amount": -10}`,
			mockBehavior:         func(r *mock.MockDatabaseBooksManager, id int, book models.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"amount can't be less than 0\"}\n",
		},
		{
			name:      "Update id not found",
			inputId:   1,
			inputBody: `{"name": "Book1", "genre": 1, "price": 0, "amount": 0}`,
			inputBook: models.Book{
				Name:   "Book1",
				Genre:  1,
				Price:  0,
				Amount: 0,
			},
			mockBehavior: func(r *mock.MockDatabaseBooksManager, id int, book models.Book) {
				r.EXPECT().UpdateBookByID(id, book).Return(0, errors.New("id not found"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"status_text\":\"Bad request\",\"message\":\"id not found\"}\n",
		},
		{
			name:      "Update ok",
			inputId:   1,
			inputBody: `{"name": "Book1", "genre": 1, "price": 0, "amount": 0}`,
			inputBook: models.Book{
				Name:   "Book1",
				Price:  0,
				Genre:  1,
				Amount: 0,
			},
			mockBehavior: func(r *mock.MockDatabaseBooksManager, id int, book models.Book) {
				r.EXPECT().UpdateBookByID(id, book).Return(id, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1}\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock.NewMockDatabaseBooksManager(c)
			test.mockBehavior(mockManager, test.inputId, test.inputBook)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := chi.NewRouter()
			r.Route("/{bookID}", func(router chi.Router) {
				router.Use(handler.BookContext)
				router.Put("/", handler.UpdateBook)
			})
			var target string
			if test.useFalseID {
				target = fmt.Sprintf("/%v", test.falseId)
			} else {
				target = fmt.Sprintf("/%v", test.inputId)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", target, bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
