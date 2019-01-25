package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	null "gopkg.in/guregu/null.v3"
)

// Book represents a book record in database
type Book struct {
	ID          uint        `gorm:"primary_key" json:"id"`
	ISBN        null.String `gorm:"size:13" json:"isbn"`
	Title       string      `gorm:"size:255" json:"title"`
	Description string      `json:"description"`
	Language    string      `gorm:"size:2" json:"language"`
}

// Books represents a collection of books and their count
type Books struct {
	NumberBooks uint   `json:"numberBooks"`
	Books       []Book `json:"books"`
}

// GetAll retrieve all books from database
func (b *Books) GetAll(db *gorm.DB) error {
	var books []Book

	if dbc := db.Find(&books); dbc.Error != nil {
		if dbc.RecordNotFound() {
			b.NumberBooks = 0
			b.Books = make([]Book, 0)
			return nil
		}

		return errors.New("Failed to retrieve books from database")
	}

	b.NumberBooks = uint(len(books))
	b.Books = books

	return nil
}
