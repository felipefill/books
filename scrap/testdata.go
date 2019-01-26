package main

import (
	"github.com/felipefill/books/model"
	null "gopkg.in/guregu/null.v3"
)

var sampleBooksISBNs = []string{
	"9783161484100", // First book's page has ISBN
	"Unavailable",   // Second book's page does not have ISBN
	"Unavailable",   // Third book has no page thus no ISBN
}

var sampleBooksUsedInLocalWebsite = []model.Book{
	model.Book{
		ID:          0,
		Title:       "Awesome book number one",
		Description: "This book was created by me and it's really great, please read it. Oh, this was also my first pharagraph. So, this is my second paragraph and I think I'll write another one after this. Yep, last one I swear. Oh, by the way, here's another link to my book1.",
		ISBN:        null.StringFrom("9783161484100"),
		Language:    "EN",
	},

	model.Book{
		ID:          0,
		Title:       "Awesome book number two",
		Description: "This book was created by me and it's really great, not as great as the first one. Sequels, right? Yep, last paragraph I swear. Oh, by the way, here's another link to my book2. I fooled you! Here's another paragraph.",
		ISBN:        null.StringFrom("Unavailable"),
		Language:    "EN",
	},

	model.Book{
		ID:          0,
		Title:       "My not so awesome book",
		Description: "I won't link this to its own page because it doesn't even have one",
		ISBN:        null.StringFrom("Unavailable"),
		Language:    "EN",
	},
}
