package models

import (
	"fmt"
	"net/http"
)

type Book struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Genre  int     `json:"genre"`
	Price  float64 `json:"price"`
	Amount int     `json:"amount"`
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
