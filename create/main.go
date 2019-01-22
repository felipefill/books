package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" {
		return events.APIGatewayProxyResponse{Body: `{"error": "Body cannot be empty"}`, StatusCode: 400}, nil
	}

	createBookRequest, err := NewCreateBookRequestFromJSONString(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf(`{"error": "%s"}`, err.Error()), StatusCode: 400}, nil
	}

	book, err := createBookRequest.StoreInDatabase()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: `{"error": "Failed to store book"}`, StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: fmt.Sprintf(`{"book_id": %d}`, book.ID), StatusCode: 201}, nil
}

func main() {
	lambda.Start(Handler)
}
