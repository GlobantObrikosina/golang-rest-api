package models

import (
	"fmt"
	"net/http"
)

type Book struct {
	ID     int     `json:"id"`
	Name   string  `json:"name" binding:"min=1,max=100"`
	Genre  int     `json:"genre" binding:"min=1,max=3"`
	Price  float64 `json:"price" binding:"min=0"`
	Amount int     `json:"amount" binding:"min=0"`
}
type BookList struct {
	Books []Book `json:"books"`
}

func (i *Book) Bind(r *http.Request) error {
	if i.Name == "" {
		return fmt.Errorf("name is a required field")
	} else if len(i.Name) > 100 {
		return fmt.Errorf("name has to be less than 100 characters")
	} else if i.Price < 0 {
		return fmt.Errorf("price can't be less than 0")
	} else if i.Genre < 1 || i.Genre > 3 {
		return fmt.Errorf("genre has to be between 1 and 3")
	} else if i.Amount < 0 {
		return fmt.Errorf("amount can't be less than 0")
	}
	return nil
}

func (i *BookList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Book) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type CreateBookResponse struct {
	BookID int `json:"id"`
}

func (CreateBookResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetBooks struct {
	Name  string `json:"name"`
	Genre int    `json:"genre"`
}

func (i *GetBooks) Bind(r *http.Request) error {
	if len(i.Name) > 100 {
		return fmt.Errorf("name has to be less than 100 characters")
	}
	return nil
}
