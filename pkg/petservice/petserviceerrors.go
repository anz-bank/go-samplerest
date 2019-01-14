package petservice

import (
	"fmt"
	"net/http"
)

// APIError defines an error that separates internal and external error messages
type APIError struct {
	InternalMessage string
	ResponseMessage string
	StatusCode      int
}

func (e APIError) Error() string {
	return e.InternalMessage
}

var (
	keyExtractionError = APIError{
		InternalMessage: "Unknown error extracting petID from request context",
		ResponseMessage: "An Unkown internal error occured",
		StatusCode:      http.StatusInternalServerError,
	}
	keyMissingError = APIError{
		InternalMessage: "Key missing",
		ResponseMessage: "An unkown internal error occured",
		StatusCode:      http.StatusBadRequest,
	}
)

func invalidPetIDError(petID string) APIError {
	return APIError{
		InternalMessage: fmt.Sprintf("Invalid pet ID %s", petID),
		ResponseMessage: fmt.Sprintf("Invalid pet ID %s", petID),
		StatusCode:      http.StatusBadRequest,
	}
}

func petIDNotFoundError(petID int32, err error) APIError {
	return APIError{
		InternalMessage: fmt.Sprintf("Pet ID not found. %v", err),
		ResponseMessage: fmt.Sprintf("No pet found with ID %d", petID),
		StatusCode:      http.StatusNotFound,
	}
}

func badRequestBodyError(err error) APIError {
	return APIError{
		InternalMessage: fmt.Sprintf("Error reading request body. %v", err),
		ResponseMessage: fmt.Sprintf("Could not read request body"),
		StatusCode:      http.StatusBadRequest,
	}
}

func badPetDataError(err error) APIError {
	return APIError{
		InternalMessage: fmt.Sprintf("Request pet data invalid. %v", err),
		ResponseMessage: fmt.Sprintf("Request pet data invalid"),
		StatusCode:      http.StatusBadRequest,
	}
}
