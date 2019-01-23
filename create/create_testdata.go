package main

import (
	"github.com/felipefill/books/model"
	null "gopkg.in/guregu/null.v3"
)

var validCreateBookRequestAsJSONString = `{
"title": "Book title example",
"description": "Book description example",
"isbn": "9781617293290",
"language": "BR"
}`

var validCreateBookRequest = CreateBookRequest{
	Title:       null.StringFrom("Book title example"),
	Description: null.StringFrom("Book description example"),
	ISBN:        null.StringFrom("9781617293290"),
	Language:    null.StringFrom("BR"),
}

var invalidCreateBookRequest = CreateBookRequest{
	Title:       null.StringFrom(""),
	Description: null.String{},
	ISBN:        null.StringFrom("9781617293290"),
	Language:    null.StringFrom("BR"),
}

var sampleBook = model.Book{
	Title:       "Book title example",
	Description: "Book description example",
	ISBN:        null.StringFrom("9781617293290"),
	Language:    "BR",
}
