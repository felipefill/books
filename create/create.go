package main

import (
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

// NewCreateBookRequestFromJSONString
func NewCreateBookRequestFromJSONString() (CreateBookRequest, error) {

}

// Validate checks struct for errors
func (request *CreateBookRequest) Validate() error {
	errorString := ""

	if !request.Title.Valid || request.Title.Value == "" {
		errorString = "Title cannot be null nor empty"
	}

	if !request.Description.Valid || request.Description.Value == "" {
		if errorString != "" {
			errorString += "; "
		}

		errorString += "Description cannot be null nor empty"
	}

	if !request.ISBN.Valid || request.ISBN.Value == "" {
		if errorString != "" {
			errorString += "; "
		}

		errorString += "ISBN cannot be null nor empty"
	}

	if !request.Language.Valid || request.Language.Value == "" {
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
