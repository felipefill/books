package model

import (
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
