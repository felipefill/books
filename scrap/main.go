package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/felipefill/books/model"
	"github.com/felipefill/books/utils"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	workingMode := retrieveWorkingMode(request)

	switch workingMode {
	case ScrapOnly:
		return scrapBooksAndReturn()
	case ScrapAndStore:
		return scrapAndStoreBooksThenReturn()
	default:
		return retrieveAllStoredBooks()
	}
}

func scrapBooksAndReturn() (events.APIGatewayProxyResponse, error) {
	scrappedBooks, err := FindKotlinBooks("https://kotlinlang.org/docs/books.html")
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Something went wrong while searching for books", StatusCode: 500}, nil
	}

	books := model.Books{
		NumberBooks: uint(len(scrappedBooks)),
		Books:       scrappedBooks,
	}

	json, _ := json.Marshal(books)
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200}, nil
}

func scrapAndStoreBooksThenReturn() (events.APIGatewayProxyResponse, error) {
	scrappedBooks, err := FindKotlinBooks("https://kotlinlang.org/docs/books.html")
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Something went wrong while searching for books", StatusCode: 500}, nil
	}

	for _, book := range scrappedBooks {
		if err = book.StoreOrRetrieveByTitle(utils.GetDB()); err != nil {
			return events.APIGatewayProxyResponse{Body: "Something went wrong while storing scrapped books", StatusCode: 500}, nil
		}
	}

	return retrieveAllStoredBooks()
}

func retrieveAllStoredBooks() (events.APIGatewayProxyResponse, error) {
	storedBooks := model.Books{}
	if err := storedBooks.GetAll(utils.GetDB()); err != nil {
		return events.APIGatewayProxyResponse{Body: "Something went wrong while retrieving books from database", StatusCode: 500}, nil
	}

	json, _ := json.Marshal(storedBooks)
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200}, nil
}

func retrieveWorkingMode(request events.APIGatewayProxyRequest) WorkingMode {
	value, _ := request.QueryStringParameters["mode"]
	return WorkingModeFromString(value)
}

func main() {
	lambda.Start(Handler)
}
