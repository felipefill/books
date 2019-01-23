package main

import (
	"encoding/json"
	"errors"

	"github.com/felipefill/books/model"
	"github.com/felipefill/books/utils"

	null "gopkg.in/guregu/null.v3"
)

// CreateBookRequest represent a request to create a new book
type CreateBookRequest struct {
	Title       null.String `json:"title"`
	Description null.String `json:"description"`
	ISBN        null.String `json:"isbn"`
	Language    null.String `json:"language"`
}

// NewCreateBookRequestFromJSONString tries to create a new CreateBookRequest from a given JSON string
func NewCreateBookRequestFromJSONString(jsonString string) (*CreateBookRequest, error) {
	request := new(CreateBookRequest)

	if err := json.Unmarshal([]byte(jsonString), request); err != nil {
		return nil, errors.New("Failed to parse JSON string into CreateBookRequest")
	}

	return request, nil
}

// StoreInDatabase stores request content in database as a new book
func (request *CreateBookRequest) StoreInDatabase() (*model.Book, error) {
	book, err := request.ToBook()
	if err != nil {
		return nil, err
	}

	//TODO: should I handle duplicates? how?
	db := utils.GetDB().Create(book)
	if db.Error != nil {
		return nil, db.Error
	}

	return book, nil
}

// ToBook converts CreateBookRequest into a Book, runs validation before doing so
func (request *CreateBookRequest) ToBook() (*model.Book, error) {
	if err := request.validate(); err != nil {
		return nil, err
	}

	book := model.Book{
		Title:       request.Title.String,
		Description: request.Description.String,
		ISBN:        request.ISBN,
		Language:    request.Language.String,
	}

	return &book, nil
}

func (request *CreateBookRequest) validate() error {
	errorString := ""

	if !request.Title.Valid || request.Title.String == "" {
		errorString = "Title cannot be null nor empty"
	}

	if !request.Description.Valid || request.Description.String == "" {
		if errorString != "" {
			errorString += "; "
		}

		errorString += "Description cannot be null nor empty"
	}

	if !request.ISBN.Valid || request.ISBN.String == "" {
		if errorString != "" {
			errorString += "; "
		}

		errorString += "ISBN cannot be null nor empty"
	}

	if !request.Language.Valid || request.Language.String == "" {
		if errorString != "" {
			errorString += "; "
		}

		errorString += "Language cannot be null nor empty"
	}

	if errorString != "" {
		return errors.New(errorString)
	}

	return nil
}
