package main

import (
	"encoding/json"
	"errors"

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

// Validate checks struct for errors
func (request *CreateBookRequest) Validate() error {
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
