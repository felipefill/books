package main

import (
	"github.com/felipefill/books/model"
	null "gopkg.in/guregu/null.v3"
)

var sampleBook = model.Book{
	ID:          99,
	Title:       "Sample book",
	Description: "This is a great book, 10/10.",
	ISBN:        null.StringFrom("0123456789012"),
	Language:    "EN",
}

var sampleBookAsJSONString = `{"id":99,"isbn":"0123456789012","title":"Sample book","description":"This is a great book, 10/10.","language":"EN"}`
