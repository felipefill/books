package main

import (
	"errors"
	"strconv"
	"testing"

	"github.com/felipefill/books/model"
	"github.com/felipefill/books/utils"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSearchHandlerFindsBook(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	request.PathParameters = make(map[string]string)
	request.PathParameters["id"] = strconv.Itoa(int(sampleBook.ID))

	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(sampleBook.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
				AddRow(sampleBook.ID, sampleBook.Title, sampleBook.Description, sampleBook.ISBN.String, sampleBook.Language),
		)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       sampleBookAsJSONString,
		StatusCode: 200,
	}

	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, expectedError, actualError)
}

func TestSearchHandlerDoesNotFindBook(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	request.PathParameters = make(map[string]string)
	request.PathParameters["id"] = "20"

	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(20).
		WillReturnError(gorm.ErrRecordNotFound)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: 404,
	}

	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, expectedError, actualError)
}

func TestSearchHandlerFailsDueToDatabase(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	request.PathParameters = make(map[string]string)
	request.PathParameters["id"] = "20"

	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(20).
		WillReturnError(gorm.ErrCantStartTransaction)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       `{"error": "Failed to retrieve book with ID: 20"}`,
		StatusCode: 500,
	}

	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, expectedError, actualError)
}

func TestSearchHandlerFailsIDMissing(t *testing.T) {
	request := events.APIGatewayProxyRequest{}

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       `{"error": "Missing "id" parameter"}`,
		StatusCode: 400,
	}

	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, expectedError, actualError)
}

func TestSearchHandlerFailsIDNotInt(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	request.PathParameters = make(map[string]string)
	request.PathParameters["id"] = "not_int"

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       `{"error": ""id" parameter must be an integer"}`,
		StatusCode: 400,
	}

	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, expectedError, actualError)
}

func TestFindBookByIDFindsBook(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	var expectedBook = sampleBook
	var expectedError error

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(22).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "isbn", "language"}).
				AddRow(expectedBook.ID, expectedBook.Title, expectedBook.Description, expectedBook.ISBN.String, expectedBook.Language),
		)

	actualBook, actualError := findBookByID(22)

	assert.Equal(t, &expectedBook, actualBook)
	assert.Equal(t, expectedError, actualError)
}

func TestFindBookByIDDoesNotFindBook(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(22).
		WillReturnError(gorm.ErrRecordNotFound)

	var expectedBook *model.Book
	var expectedError error

	actualBook, actualError := findBookByID(22)

	assert.Equal(t, expectedBook, actualBook)
	assert.Equal(t, expectedError, actualError)
}

func TestFindBookByIDFailsDueToDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(22).
		WillReturnError(gorm.ErrCantStartTransaction)

	var expectedBook *model.Book
	expectedError := errors.New("Failed to retrieve book with ID: 22")

	actualBook, actualError := findBookByID(22)

	assert.Equal(t, expectedBook, actualBook)
	assert.Equal(t, expectedError, actualError)
}

func TestRetrieveIDFromRequest(t *testing.T) {
	request := events.APIGatewayProxyRequest{}

	expectedID := -1
	expectedError := errors.New("Missing \"id\" parameter")

	actualID, actualError := retrieveIDFromRequest(request)

	assert.Equal(t, expectedID, actualID)
	assert.Equal(t, expectedError, actualError)

	request.PathParameters = make(map[string]string)
	request.PathParameters["id"] = "not_int"

	expectedID = -1
	expectedError = errors.New("\"id\" parameter must be an integer")

	actualID, actualError = retrieveIDFromRequest(request)

	assert.Equal(t, expectedID, actualID)
	assert.Equal(t, expectedError, actualError)

	request.PathParameters["id"] = "18"

	expectedID = 18
	expectedError = nil

	actualID, actualError = retrieveIDFromRequest(request)

	assert.Equal(t, expectedID, actualID)
	assert.Equal(t, expectedError, actualError)
}
