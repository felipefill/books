package main

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/felipefill/books/utils"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookHandlerHappyPath(t *testing.T) {
	request := events.APIGatewayProxyRequest{Body: validCreateBookRequestAsJSONString}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(sampleBook.Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(sampleBook.ISBN.String, sampleBook.Title, sampleBook.Description, sampleBook.Language).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{Body: `{"book_id": 1}`, StatusCode: 201}
	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateBookHandlerFailsBodyIsEmpty(t *testing.T) {
	request := events.APIGatewayProxyRequest{Body: ""}

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{Body: `{"error": "Body cannot be empty"}`, StatusCode: 400}
	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateBookHandlerFailsBodyIsInvalid(t *testing.T) {
	request := events.APIGatewayProxyRequest{Body: "not even a valid request body"}

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{Body: `{"error": "Failed to parse JSON string into CreateBookRequest"}`, StatusCode: 400}
	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateBookHandlerFailsDatabaseError(t *testing.T) {
	request := events.APIGatewayProxyRequest{Body: validCreateBookRequestAsJSONString}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDB(db)

	mock.
		ExpectQuery("SELECT (.+) FROM \"books\" (.+)").
		WithArgs(sampleBook.Title).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.
		ExpectQuery("INSERT INTO \"books\" \\(\"isbn\",\"title\",\"description\",\"language\"\\)").
		WithArgs(sampleBook.ISBN.String, sampleBook.Title, sampleBook.Description, sampleBook.Language).
		WillReturnError(errors.New("some error"))

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{Body: `{"error": "Failed to store book"}`, StatusCode: 500}
	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
}
