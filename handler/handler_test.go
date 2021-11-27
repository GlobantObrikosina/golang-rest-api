package handler

//import (
//	"bytes"
//	"errors"
//	"github.com/GlobantObrikosina/golang-rest-api/models"
//	"github.com/GlobantObrikosina/golang-rest-api/service"
//	mock "github.com/GlobantObrikosina/golang-rest-api/service/mocks"
//	"github.com/go-chi/chi"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestHandler_CreateBook(t *testing.T) {
//	type mockBehavior func(s *mock.MockDatabaseBooksManager, book models.Book)
//	tests := []struct {
//		name                 string
//		inputBody            string
//		inputBook            models.Book
//		mockBehavior         mockBehavior
//		expectedStatusCode   int
//		expectedResponseBody string
//	}{
//		{
//			name:      "Ok",
//			inputBody: `{"name": "Book1", "price": 0, "genre": 1, "amount": 0}`,
//			inputBook: models.Book{
//				ID: 1,
//				Name:   "Book1",
//				Price:  0,
//				Genre:  1,
//				Amount: 0,
//			},
//			mockBehavior: func(r *mock.MockDatabaseBooksManager, book models.Book) {
//				r.EXPECT().CreateBook(book).Return(1, nil)
//			},
//			expectedStatusCode:   http.StatusOK,
//			expectedResponseBody: `{"id":1}`,
//		},
//		{
//			name:                 "Some fields missing",
//			inputBody:            `{"price": 67.88, "genre": 1, "amount": 5}`,
//			mockBehavior:         func(r *mock.MockDatabaseBooksManager, book models.Book) {},
//			expectedStatusCode:   http.StatusBadRequest,
//			expectedResponseBody: `{"error":"invalid input"}`,
//		},
//		{
//			name:                 "Invalid genre",
//			inputBody:            `{"name": "hello", "price": 67.88, "genre": 6, "amount": 7}`,
//			mockBehavior:         func(r *mock.MockDatabaseBooksManager, book models.Book) {},
//			expectedStatusCode:   http.StatusBadRequest,
//			expectedResponseBody: `{"error":"invalid input"}`,
//		},
//		{
//			name:      "Unique name",
//			inputBody: `{"name": "hello", "price": 67.88, "genre": 1, "amount": 7}`,
//			inputBook: models.Book{
//				ID: 1,
//				Name:   "hello",
//				Price:  67.88,
//				Genre:  1,
//				Amount: 7,
//			},
//			mockBehavior: func(r *mock.MockDatabaseBooksManager, book models.Book) {
//				r.EXPECT().CreateBook(book).Return(0, errors.New("input book name isn't unique"))
//			},
//			expectedStatusCode:   http.StatusInternalServerError,
//			expectedResponseBody: `{"error":"input book name isn't unique"}`,
//		},
//		{
//			name:      "Unique name",
//			inputBody: `{"name": "hello", "price": 67.88, "genre": 1, "amount": 7}`,
//			inputBook: models.Book{
//				ID: 1,
//				Name:   "hello",
//				Price:  67.88,
//				Genre:  1,
//				Amount: 7,
//			},
//			mockBehavior: func(r *mock.MockDatabaseBooksManager, book models.Book) {
//				r.EXPECT().CreateBook(book).Return(0, errors.New("input book name isn't unique"))
//			},
//			expectedStatusCode:   http.StatusInternalServerError,
//			expectedResponseBody: `{"error":"input book name isn't unique"}`,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			mockManager := mock.MockDatabaseBooksManager(c)
//			test.mockBehavior(&mockManager, test.inputBook)
//
//			r := chi.NewRouter()
//			r.Post("/books", CreateBook)
//
//			w := httptest.NewRecorder()
//			req := httptest.NewRequest("POST", "/books",
//				bytes.NewBufferString(test.inputBody))
//
//			r.ServeHTTP(w, req)
//
//			assert.Equal(t, test.expectedStatusCode, w.Code)
//			assert.Equal(t, test.expectedResponseBody, w.Body.String())
//		})
//	}
//
//}
//
