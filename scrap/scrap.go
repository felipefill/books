package main

import (
	"regexp"
	"strings"

	null "gopkg.in/guregu/null.v3"

	"github.com/felipefill/books/model"

	"github.com/gocolly/colly"
)

func scrapBooksElements(booksIndex string) (booksElements [][]*colly.HTMLElement, scrapingError error) {
	booksElements = make([][]*colly.HTMLElement, 0)

	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		scrapingError = err
	})

	var currentBookElements []*colly.HTMLElement
	c.OnHTML("article", func(article *colly.HTMLElement) {
		// Article is the parent for all books
		// This will find all children elements of article
		article.ForEach("*", func(index int, element *colly.HTMLElement) {
			if element.Name == "h2" {
				if len(currentBookElements) > 0 {
					// Means we already have collected elements for the first book
					booksElements = append(booksElements, currentBookElements)
				}
				// Found new book
				currentBookElements = make([]*colly.HTMLElement, 0)
			}
			currentBookElements = append(currentBookElements, element)
		})

		if len(currentBookElements) > 0 {
			booksElements = append(booksElements, currentBookElements)
		}
	})

	c.Visit(booksIndex)
	c.Wait()

	if scrapingError != nil {
		return nil, scrapingError
	}

	return booksElements, nil
}

func scrapISBN(link string) (isbn string, scrapingError error) {
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		scrapingError = err
	})

	// How to stop OnHTML once I've found the ISBN? Is there a way to get ONLY the whole HTML?
	c.OnHTML("body", func(element *colly.HTMLElement) {
		indexOf := strings.Index(element.Text, "978")
		if indexOf != -1 {
			isbn = element.Text[indexOf : indexOf+26]
			isbn = strings.Replace(isbn, "-", "", -1)
			isbn = isbn[0:13]
		} else {
			indexOf = strings.Index(element.Text, "979")
			if indexOf != -1 {
				isbn = element.Text[indexOf : indexOf+26]
				isbn = strings.Replace(isbn, "-", "", -1)
				isbn = isbn[0:13]
			} else {
				isbn = "Unavailable"
			}
		}
	})

	c.Visit(link)
	c.Wait()

	return
}

func scrapBooksISBNs(booksElements [][]*colly.HTMLElement) (booksISBNs []string, scrapingError error) {
	// We already know the len and cap of this by len(booksElements)
	booksISBNs = make([]string, 0)

	for _, currentBookElements := range booksElements {
		var isbnLink string
		var isbn = "Unavailable"

		for _, currentElement := range currentBookElements {
			if currentElement.Name == "a" {
				isbnLink = currentElement.Attr("href")
				break
			}
		}

		if isbnLink != "" {
			isbn, scrapingError = scrapISBN(isbnLink)
			if scrapingError != nil {
				return nil, scrapingError
			}
		}

		booksISBNs = append(booksISBNs, isbn)
	}

	return
}

func combineBooksElementsAndISBNsIntoBooks(booksElements [][]*colly.HTMLElement, booksISBNs []string) (books []model.Book) {
	books = make([]model.Book, 0)

	for index, bookElements := range booksElements {
		currentBook := model.Book{}
		bookDescription := ""

		for _, element := range bookElements {
			if element.Name == "h2" {
				currentBook.Title = strings.TrimSpace(element.Text)
			}

			if element.Name == "p" {
				text := strings.Replace(element.Text, "\n", " ", -1)
				text = strings.Replace(text, "\t", " ", -1) + " "

				bookDescription += text
			}

			if element.Name == "div" {
				currentBook.Language = strings.ToUpper(element.Text)
			}
		}

		bookDescription = strings.TrimSpace(bookDescription)
		bookDescription = regexp.MustCompile(`[\s\p{Zs}]{2,}`).ReplaceAllString(bookDescription, " ")

		currentBook.Description = bookDescription
		currentBook.ISBN = null.StringFrom(booksISBNs[index])

		books = append(books, currentBook)
	}

	return
}

// FindKotlinBooks scraps and scraps Kotlin website's books section searching for new books for our library
func FindKotlinBooks(kotlinBooksURL string) ([]model.Book, error) {
	scrappedBooks, err := scrapBooksElements(kotlinBooksURL)
	if err != nil {
		return nil, err
	}

	scrapBooksISBN, err := scrapBooksISBNs(scrappedBooks)
	if err != nil {
		return nil, err
	}

	books := combineBooksElementsAndISBNsIntoBooks(scrappedBooks, scrapBooksISBN)

	return books, nil
}
