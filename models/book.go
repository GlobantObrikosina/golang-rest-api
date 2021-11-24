package models

import (
	"fmt"
	"net/http"
)

type Book struct {
	ID     int     `json:"id"`
	Name   string  `json:"name" binding:"min=1,max=100"`
	Genre  int     `json:"genre" binding:"min=0"`
	Price  float64 `json:"price" binding:"min=1,max=3"`
	Amount int     `json:"amount" binding:"min=0"`
}
type BookList struct {
	Books []Book `json:"books"`
}

func (i *Book) Bind(r *http.Request) error {
	if i.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}

func (*BookList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Book) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
