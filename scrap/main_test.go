package main

import (
	"encoding/json"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/felipefill/books/model"
	"github.com/felipefill/books/utils"
	"github.com/jinzhu/gorm"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveWorkingMode(t *testing.T) {
	expected := RetrieveAll
	actual := retrieveWorkingMode(events.APIGatewayProxyRequest{})

	assert.Equal(t, expected, actual)
}

func TestRetrieveAllStoredBooksSucceeds(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	books := sampleBooksUsedInLocalWebsite
	books[0].ID = 1
	books[1].ID = 2
	books[2].ID = 3

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\"").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
				AddRow(books[0].ID, books[0].Title, books[0].Description, books[0].ISBN.String, books[0].Language).
				AddRow(books[1].ID, books[1].Title, books[1].Description, books[1].ISBN.String, books[1].Language).
				AddRow(books[2].ID, books[2].Title, books[2].Description, books[2].ISBN.String, books[2].Language),
		)

	booksResponse := model.Books{
		NumberBooks: uint(len(books)),
		Books:       books,
	}

	booksResponseJSON, _ := json.Marshal(&booksResponse)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       string(booksResponseJSON),
		StatusCode: 200,
	}

	actualResponse, actualError := retrieveAllStoredBooks()

	books[0].ID = 0
	books[1].ID = 0
	books[2].ID = 0

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestRetrieveAllStoredBooksFails(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\"").
		WillReturnError(errors.New("database error"))

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       `{"error": "Something went wrong while retrieving books from database"}`,
		StatusCode: 500,
	}

	actualResponse, actualError := retrieveAllStoredBooks()

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestScrapAndStoreBooksThenReturnSucceeds(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	ts := createTestServer()
	defer ts.Close()

	utils.InjectDB(db)

	books := sampleBooksUsedInLocalWebsite
	books[0].ID = 1
	books[1].ID = 2
	books[2].ID = 3

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(books[0].Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(books[0].ISBN.String, books[0].Title, books[0].Description, books[0].Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(books[0].ID))

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(books[1].Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(books[1].ISBN.String, books[1].Title, books[1].Description, books[1].Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(books[1].ID))

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(books[2].Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(books[2].ISBN.String, books[2].Title, books[2].Description, books[2].Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(books[2].ID))

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\"").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
				AddRow(books[0].ID, books[0].Title, books[0].Description, books[0].ISBN.String, books[0].Language).
				AddRow(books[1].ID, books[1].Title, books[1].Description, books[1].ISBN.String, books[1].Language).
				AddRow(books[2].ID, books[2].Title, books[2].Description, books[2].ISBN.String, books[2].Language),
		)

	booksResponse := model.Books{
		NumberBooks: uint(len(books)),
		Books:       books,
	}

	booksResponseJSON, _ := json.Marshal(&booksResponse)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       string(booksResponseJSON),
		StatusCode: 200,
	}

	actualResponse, actualError := scrapAndStoreBooksThenReturn(ts.URL + "/index.html")

	books[0].ID = 0
	books[1].ID = 0
	books[2].ID = 0

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestScrapAndStoreBooksThenReturnFailsToScrapBooks(t *testing.T) {
	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       `{"error": "Something went wrong while searching for books"}`,
		StatusCode: 500,
	}

	actualResponse, actualError := scrapAndStoreBooksThenReturn("not_a_url")

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestScrapAndStoreBooksThenReturnFailsDueToDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	ts := createTestServer()
	defer ts.Close()

	utils.InjectDB(db)

	books := sampleBooksUsedInLocalWebsite

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(books[0].Title).
		WillReturnError(errors.New("database error"))

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       `{"error": Something went wrong while storing scrapped books"}`,
		StatusCode: 500,
	}

	actualResponse, actualError := scrapAndStoreBooksThenReturn(ts.URL + "/index.html")

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestScrapBooksAndReturnSucceeds(t *testing.T) {
	ts := createTestServer()
	defer ts.Close()

	books := sampleBooksUsedInLocalWebsite

	booksResponse := model.Books{
		NumberBooks: uint(len(books)),
		Books:       books,
	}

	booksResponseJSON, _ := json.Marshal(&booksResponse)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       string(booksResponseJSON),
		StatusCode: 200,
	}

	actualResponse, actualError := scrapBooksAndReturn(ts.URL + "/index.html")

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestScrapBooksAndReturnFailsToScrapBooks(t *testing.T) {
	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       `{"error": "Something went wrong while searching for books"}`,
		StatusCode: 500,
	}

	actualResponse, actualError := scrapBooksAndReturn("not_a_url")

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestHandlerScrapOnly(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	request.QueryStringParameters = make(map[string]string)
	request.QueryStringParameters["mode"] = "scrap_only"

	ts := createTestServer()
	defer ts.Close()

	kotlinBooksURL = ts.URL + "/index.html"

	books := sampleBooksUsedInLocalWebsite

	booksResponse := model.Books{
		NumberBooks: uint(len(books)),
		Books:       books,
	}

	booksResponseJSON, _ := json.Marshal(&booksResponse)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       string(booksResponseJSON),
		StatusCode: 200,
	}

	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestHandlerScrapAndStore(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	request.QueryStringParameters = make(map[string]string)
	request.QueryStringParameters["mode"] = "scrap_and_store"

	ts := createTestServer()
	defer ts.Close()

	db, mock, _ := sqlmock.New()
	defer db.Close()

	kotlinBooksURL = ts.URL + "/index.html"
	utils.InjectDB(db)

	books := sampleBooksUsedInLocalWebsite
	books[0].ID = 1
	books[1].ID = 2
	books[2].ID = 3

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(books[0].Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(books[0].ISBN.String, books[0].Title, books[0].Description, books[0].Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(books[0].ID))

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(books[1].Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(books[1].ISBN.String, books[1].Title, books[1].Description, books[1].Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(books[1].ID))

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(books[2].Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(books[2].ISBN.String, books[2].Title, books[2].Description, books[2].Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(books[2].ID))

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\"").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
				AddRow(books[0].ID, books[0].Title, books[0].Description, books[0].ISBN.String, books[0].Language).
				AddRow(books[1].ID, books[1].Title, books[1].Description, books[1].ISBN.String, books[1].Language).
				AddRow(books[2].ID, books[2].Title, books[2].Description, books[2].ISBN.String, books[2].Language),
		)

	booksResponse := model.Books{
		NumberBooks: uint(len(books)),
		Books:       books,
	}

	booksResponseJSON, _ := json.Marshal(&booksResponse)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       string(booksResponseJSON),
		StatusCode: 200,
	}

	actualResponse, actualError := Handler(request)

	books[0].ID = 0
	books[1].ID = 0
	books[2].ID = 0

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestHandlerRetrieveAll(t *testing.T) {
	ts := createTestServer()
	defer ts.Close()

	db, mock, _ := sqlmock.New()
	defer db.Close()

	kotlinBooksURL = ts.URL + "/index.html"
	utils.InjectDB(db)

	books := sampleBooksUsedInLocalWebsite
	books[0].ID = 1
	books[1].ID = 2
	books[2].ID = 3

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\"").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
				AddRow(books[0].ID, books[0].Title, books[0].Description, books[0].ISBN.String, books[0].Language).
				AddRow(books[1].ID, books[1].Title, books[1].Description, books[1].ISBN.String, books[1].Language).
				AddRow(books[2].ID, books[2].Title, books[2].Description, books[2].ISBN.String, books[2].Language),
		)

	booksResponse := model.Books{
		NumberBooks: uint(len(books)),
		Books:       books,
	}

	booksResponseJSON, _ := json.Marshal(&booksResponse)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       string(booksResponseJSON),
		StatusCode: 200,
	}

	actualResponse, actualError := Handler(events.APIGatewayProxyRequest{})

	books[0].ID = 0
	books[1].ID = 0
	books[2].ID = 0

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}
