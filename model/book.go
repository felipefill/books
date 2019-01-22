package model

import (
	"database/sql"
)

// Book represents a book record in database
type Book struct {
	ID          uint           `gorm:"primary_key"`
	ISBN        sql.NullString `gorm:"size:13"`
	Title       string         `gorm:"size:255"`
	Description string
	Language    string `gorm:"size:2"`
}
