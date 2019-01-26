package model

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestStoreOrRetrieveByTitleRetrieve(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open("postgres", db)

	var expectedError error
	expectedBook := sampleBook

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(expectedBook.Title).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
			AddRow(expectedBook.ID, expectedBook.Title, expectedBook.Description, expectedBook.ISBN.String, expectedBook.Language),
		)

	actualBook := Book{
		Title: sampleBook.Title,
	}

	actualError := actualBook.StoreOrRetrieveByTitle(gormDB)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBook, actualBook)
}

func TestStoreOrRetrieveByTitleStore(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open("postgres", db)

	var expectedError error
	expectedBook := sampleBook
	expectedBook.ID = 1

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(expectedBook.Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(expectedBook.ISBN.String, expectedBook.Title, expectedBook.Description, expectedBook.Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	actualBook := Book{
		Title:       sampleBook.Title,
		Description: sampleBook.Description,
		Language:    sampleBook.Language,
		ISBN:        sampleBook.ISBN,
	}

	actualError := actualBook.StoreOrRetrieveByTitle(gormDB)

	sampleBook.ID = 0

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBook, actualBook)
}

func TestGetAllNoRecords(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open("postgres", db)

	var expectedError error
	expectedBooks := Books{
		NumberBooks: 0,
		Books:       []Book{},
	}

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\"").
		WillReturnError(gorm.ErrRecordNotFound)

	actualBooks := Books{}
	actualError := actualBooks.GetAll(gormDB)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBooks, actualBooks)
}

func TestGetAllFails(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open("postgres", db)

	expectedError := errors.New("Failed to retrieve books from database")
	expectedBooks := Books{}

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\"").
		WillReturnError(errors.New("database error"))

	actualBooks := Books{}
	actualError := actualBooks.GetAll(gormDB)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBooks, actualBooks)
}

func TestGetAllRetrievesBook(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open("postgres", db)

	var expectedError error
	expectedBooks := Books{
		NumberBooks: 1,
		Books:       []Book{sampleBook},
	}

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\"").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
			AddRow(sampleBook.ID, sampleBook.Title, sampleBook.Description, sampleBook.ISBN.String, sampleBook.Language),
		)

	actualBooks := Books{}
	actualError := actualBooks.GetAll(gormDB)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedBooks, actualBooks)
}
