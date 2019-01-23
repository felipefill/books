package main

import (
	"encoding/json"

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
	books, err := FindKotlinBooks()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Something went wrong while searching for books", StatusCode: 500}, nil
	}

	json, _ := json.Marshal(books)
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
