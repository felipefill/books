package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/felipefill/books/model"
	"github.com/felipefill/books/utils"

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
	params := request.PathParameters
	idAsString, ok := params["id"]
	if !ok {
		return events.APIGatewayProxyResponse{Body: `{"error": "Missing \"id\" parameter"}`, StatusCode: 400}, nil
	}

	id, err := strconv.Atoi(idAsString)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: `{"error": "\"id\" param must be an integer"}`, StatusCode: 400}, nil
	}

	book, err := FindBookByID(id)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf(`{"error": %s}`, err.Error()), StatusCode: 400}, nil
	}

	if book == nil {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 404}, nil
	}

	json, err := json.Marshal(book)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: `{"error": "Sorry, something went wrong on our side"}`, StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200}, nil
}

// FindBookByID tries to find a book by it's ID, returns nil if not found
func FindBookByID(id int) (*model.Book, error) {
	book := model.Book{}
	db := utils.GetDB()

	if db.Where("id = ?", id).First(&book).RecordNotFound() {
		return nil, nil
	}

	if db.Error != nil {
		return nil, fmt.Errorf("Failed to retrieve book with ID: %d", id)
	}

	return &book, nil
}

func main() {
	lambda.Start(Handler)
}
