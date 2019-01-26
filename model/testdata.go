package model

import null "gopkg.in/guregu/null.v3"

var sampleBook = Book{
	Title:       "Book title example",
	Description: "Book description example",
	ISBN:        null.StringFrom("9781617293290"),
	Language:    "BR",
}
