package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookHandler(t *testing.T) {
	request := events.APIGatewayProxyRequest{}

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{Body: "", StatusCode: 200}
	actualResponse, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponse)
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
	var expectedCreateBookRequest *CreateBookRequest = nil

	actualCreateBookRequest, actualError := NewCreateBookRequestFromJSONString(jsonString)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedCreateBookRequest, actualCreateBookRequest)
}

func TestCreateBookRequestValidate(t *testing.T) {
	request := CreateBookRequest{}
	expectedError := errors.New("Title cannot be null nor empty; Description cannot be null nor empty; ISBN cannot be null nor empty; Language cannot be null nor empty")
	actualError := request.Validate()
	assert.Equal(t, expectedError, actualError)

	_ = request.Description.Scan("This is a description")
	expectedError = errors.New("Title cannot be null nor empty; ISBN cannot be null nor empty; Language cannot be null nor empty")
	actualError = request.Validate()
	assert.Equal(t, expectedError, actualError)

	_ = request.Title.Scan("This is a title")
	expectedError = errors.New("ISBN cannot be null nor empty; Language cannot be null nor empty")
	actualError = request.Validate()
	assert.Equal(t, expectedError, actualError)

	_ = request.Language.Scan("EN")
	expectedError = errors.New("ISBN cannot be null nor empty")
	actualError = request.Validate()
	assert.Equal(t, expectedError, actualError)

	_ = request.ISBN.Scan("9781234567890")
	expectedError = nil
	actualError = request.Validate()
	assert.Equal(t, expectedError, actualError)
}
