package main

import (
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
