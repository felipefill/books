package main

import (
	"errors"
	"testing"

	"github.com/felipefill/books/model"
	"github.com/jinzhu/gorm"

	"github.com/felipefill/books/utils"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookRequestStoreInDatabaseFailsInvalidRequest(t *testing.T) {
	request := invalidCreateBookRequest

	var expectedBook *model.Book
	expectedError := errors.New("Title cannot be null nor empty; Description cannot be null nor empty")

	actualBook, actualError := request.StoreInDatabase()

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBook, actualBook)
}

func TestCreateBookRequestStoreInDatabaseFailsDatabaseError(t *testing.T) {
	request := validCreateBookRequest

	var expectedBook *model.Book
	expectedError := errors.New("some database error")

	db, mock, _ := sqlmock.New()
	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(sampleBook.Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(sampleBook.ISBN.String, sampleBook.Title, sampleBook.Description, sampleBook.Language).
		WillReturnError(errors.New("some database error"))

	actualBook, actualError := request.StoreInDatabase()

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBook, actualBook)
}

func TestCreateBookRequestStoreInDatabaseSucceeds(t *testing.T) {
	request := validCreateBookRequest

	expectedBook := sampleBook
	expectedBook.ID = 1

	var expectedError error

	db, mock, _ := sqlmock.New()
	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(sampleBook.Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(sampleBook.ISBN.String, sampleBook.Title, sampleBook.Description, sampleBook.Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	actualBook, actualError := request.StoreInDatabase()

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, &expectedBook, actualBook)
}

func TestCreateBookRequestStoreInDatabaseSucceedsAlreadyInDatabase(t *testing.T) {
	request := validCreateBookRequest

	expectedBook := sampleBook
	expectedBook.ID = 1

	var expectedError error

	db, mock, _ := sqlmock.New()
	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(sampleBook.Title).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
				AddRow(expectedBook.ID, expectedBook.Title, expectedBook.Description, expectedBook.ISBN.String, expectedBook.Language),
		)

	actualBook, actualError := request.StoreInDatabase()

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, &expectedBook, actualBook)
}

func TestCreateBookRequestToBook(t *testing.T) {
	request := invalidCreateBookRequest

	var expectedBook *model.Book
	expectedError := errors.New("Title cannot be null nor empty; Description cannot be null nor empty")

	actualBook, actualError := request.ToBook()

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBook, actualBook)

	request = validCreateBookRequest
	expectedBook = &sampleBook
	expectedError = nil

	actualBook, actualError = request.ToBook()

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBook, actualBook)
}

func TestNewCreateBookRequestFromJSONString(t *testing.T) {
	jsonString := validCreateBookRequestAsJSONString

	var expectedError error
	expectedCreateBookRequest := &validCreateBookRequest

	actualCreateBookRequest, actualError := NewCreateBookRequestFromJSONString(jsonString)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedCreateBookRequest, actualCreateBookRequest)
}

func TestNewCreateBookRequestFromJSONStringFailsWithInvalidString(t *testing.T) {
	jsonString := "This is not a JSON string"

	expectedError := errors.New("Failed to parse JSON string into CreateBookRequest")
	var expectedCreateBookRequest *CreateBookRequest

	actualCreateBookRequest, actualError := NewCreateBookRequestFromJSONString(jsonString)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedCreateBookRequest, actualCreateBookRequest)
}

func TestCreateBookRequestValidate(t *testing.T) {
	request := CreateBookRequest{}
	expectedError := errors.New("Title cannot be null nor empty; Description cannot be null nor empty; ISBN cannot be null nor empty; Language cannot be null nor empty")
	actualError := request.validate()
	assert.Equal(t, expectedError, actualError)

	_ = request.Description.Scan("This is a description")
	expectedError = errors.New("Title cannot be null nor empty; ISBN cannot be null nor empty; Language cannot be null nor empty")
	actualError = request.validate()
	assert.Equal(t, expectedError, actualError)

	_ = request.Title.Scan("This is a title")
	expectedError = errors.New("ISBN cannot be null nor empty; Language cannot be null nor empty")
	actualError = request.validate()
	assert.Equal(t, expectedError, actualError)

	_ = request.Language.Scan("EN")
	expectedError = errors.New("ISBN cannot be null nor empty")
	actualError = request.validate()
	assert.Equal(t, expectedError, actualError)

	_ = request.ISBN.Scan("9781234567890")
	expectedError = nil
	actualError = request.validate()
	assert.Equal(t, expectedError, actualError)
}
