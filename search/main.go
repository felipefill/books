package main

import (
	"encoding/json"
	"errors"
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
	id, err := retrieveIDFromRequest(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf(`{"error": "%s"}`, err.Error()), StatusCode: 400}, nil
	}

	book, err := findBookByID(id)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf(`{"error": "%s"}`, err.Error()), StatusCode: 500}, nil
	}

	if book == nil {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 404}, nil
	}

	json, _ := json.Marshal(book)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: `{"error": "Sorry, something went wrong on our side"}`, StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200}, nil
}

func findBookByID(id int) (*model.Book, error) {
	book := model.Book{}
	db := utils.GetDB()
	dbc := db.Where("id = ?", id).Find(&book)

	if dbc.RecordNotFound() {
		return nil, nil
	} else if len(dbc.GetErrors()) > 0 {
		return nil, fmt.Errorf("Failed to retrieve book with ID: %d", id)
	}

	return &book, nil
}

func retrieveIDFromRequest(request events.APIGatewayProxyRequest) (int, error) {
	params := request.PathParameters
	idAsString, ok := params["id"]
	if !ok {
		return -1, errors.New("Missing \"id\" parameter")
	}

	id, err := strconv.Atoi(idAsString)
	if err != nil {
		return -1, errors.New("\"id\" parameter must be an integer")
	}

	return id, nil
}

func main() {
	lambda.Start(Handler)
}
