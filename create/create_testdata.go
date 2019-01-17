package main

import null "gopkg.in/guregu/null.v3"

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
