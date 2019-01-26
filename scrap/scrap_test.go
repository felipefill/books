package main

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/felipefill/books/model"

	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
)

func createTestServer() *httptest.Server {
	listener, _ := net.Listen("tcp", "127.0.0.1:8080")

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := ioutil.ReadFile("html/" + r.URL.Path)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		w.Write(file)
	}))

	ts.Listener.Close()
	ts.Listener = listener

	ts.Start()

	return ts
}

func TestScrapBooksElementsSucceeds(t *testing.T) {
	ts := createTestServer()
	defer ts.Close()

	var expectedErr error
	expectedFoundBooksNumber := 3
	expectedElementsInBook1 := 8
	expectedElementsInBook2 := 7
	expectedElementsInBook3 := 3

	elements, actualErr := scrapBooksElements(ts.URL + "/index.html")

	assert.Equal(t, expectedErr, actualErr)
	assert.Equal(t, expectedFoundBooksNumber, len(elements))
	assert.Equal(t, expectedElementsInBook1, len(elements[0]))
	assert.Equal(t, expectedElementsInBook2, len(elements[1]))
	assert.Equal(t, expectedElementsInBook3, len(elements[2]))
}

func TestScrapBooksElementsFails(t *testing.T) {
	var expectedElements [][]*colly.HTMLElement
	expectedError := &url.Error{
		Op:  "Get",
		URL: "http://not_a_url",
		Err: errors.New("http: no Host in request URL"),
	}

	actualElements, actualError := scrapBooksElements("not_a_url")

	assert.Equal(t, expectedElements, actualElements)
	assert.Equal(t, expectedError, actualError)
}

func TestScrapISBNSucceedsAndFindsISBN(t *testing.T) {
	ts := createTestServer()
	defer ts.Close()

	var expectedError error
	expectedISBN := "9783161484100"

	actualISBN, actualError := scrapISBN(ts.URL + "/book1.html")

	assert.Equal(t, expectedISBN, actualISBN)
	assert.Equal(t, expectedError, actualError)
}

func TestScrapISBNSucceedsAndDoesNotFindISBN(t *testing.T) {
	ts := createTestServer()
	defer ts.Close()

	var expectedError error
	expectedISBN := "Unavailable"

	actualISBN, actualError := scrapISBN(ts.URL + "/book2.html")

	assert.Equal(t, expectedISBN, actualISBN)
	assert.Equal(t, expectedError, actualError)
}

func TestScrapISBNFails(t *testing.T) {
	expectedISBN := ""
	expectedError := &url.Error{
		Op:  "Get",
		URL: "http://not_a_url",
		Err: errors.New("http: no Host in request URL"),
	}

	actualISBN, actualError := scrapISBN("not_a_url")

	assert.Equal(t, expectedISBN, actualISBN)
	assert.Equal(t, expectedError, actualError)
}

func TestScrapBooksISBNs(t *testing.T) {
	ts := createTestServer()
	defer ts.Close()

	scrappedBooksElements, scrappedBooksElementsError := scrapBooksElements(ts.URL + "/index.html")
	assert.Equal(t, nil, scrappedBooksElementsError)

	var expectedError error
	expectedISBNs := sampleBooksISBNs

	actualISBNs, actualError := scrapBooksISBNs(scrappedBooksElements)

	assert.Equal(t, expectedISBNs, actualISBNs)
	assert.Equal(t, expectedError, actualError)
}

func TestCombineBooksElementsAndISBNsIntoBooks(t *testing.T) {
	ts := createTestServer()
	defer ts.Close()

	scrappedBooksElements, scrappedBooksElementsError := scrapBooksElements(ts.URL + "/index.html")
	assert.Equal(t, nil, scrappedBooksElementsError)

	expectedBooks := sampleBooksUsedInLocalWebsite
	actualBooks := combineBooksElementsAndISBNsIntoBooks(scrappedBooksElements, sampleBooksISBNs)

	assert.Equal(t, expectedBooks, actualBooks)
}

func TestFindKotlinBooksSucceeds(t *testing.T) {
	ts := createTestServer()
	defer ts.Close()

	var expectedError error
	expectedBooks := sampleBooksUsedInLocalWebsite

	actualBooks, actualError := FindKotlinBooks(ts.URL + "/index.html")

	assert.Equal(t, expectedBooks, actualBooks)
	assert.Equal(t, expectedError, actualError)
}

func TestFindKotlinBooksFailsToScrap(t *testing.T) {
	var expectedBooks []model.Book
	expectedError := &url.Error{
		Op:  "Get",
		URL: "http://not_a_url",
		Err: errors.New("http: no Host in request URL"),
	}

	actualBooks, actualError := FindKotlinBooks("not_a_url")

	assert.Equal(t, expectedBooks, actualBooks)
	assert.Equal(t, expectedError, actualError)
}
