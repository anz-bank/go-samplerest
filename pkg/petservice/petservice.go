package petservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// PetStorer defines standard CRUD operations for petservice
type PetStorer interface {
	GetPet(int32) (*Pet, error)
	PostPet(*Pet) error
	PutPet(int32, *Pet) error
	DeletePet(int32) (bool, error)
}

// PetService defines a rest api for interaction with a PetStorer
type PetService struct {
	store PetStorer
}

// NewPetService creates a new pet service with an in-memory store
func NewPetService() *PetService {
	return &PetService{
		store: NewMemStore(),
	}
}

// nilStatus is used when it is expected to be
// overriden in the response
// eg: function is returning an error
const nilStatus = http.StatusInternalServerError

// GetPet handles a GET request to retrieve a pet
func (ps *PetService) GetPet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return nilStatus, nil, err
	}
	pet, err := ps.store.GetPet(petID)
	if err != nil {
		return nilStatus, nil, petIDNotFoundError(petID, err)
	}
	return http.StatusOK, pet, nil
}

// PostPet handles a POST request to add a new pet
func (ps *PetService) PostPet(r *http.Request) (int, interface{}, error) {
	newPet, err := readPetBody(r)
	if err != nil {
		return nilStatus, nil, err
	}

	if err = ps.store.PostPet(newPet); err != nil {
		return nilStatus, nil, APIError{
			InternalMessage: fmt.Sprintf("Error posting new pet to store. %v", err),
			ResponseMessage: "Could not POST pet to store",
			StatusCode:      http.StatusInternalServerError,
		}
	}

	return http.StatusCreated, nil, nil
}

// PutPet handles a PUT request to modify a pet
func (ps *PetService) PutPet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return nilStatus, nil, err
	}
	pet, err := readPetBody(r)
	if err != nil {
		return nilStatus, nil, err
	}

	if err = ps.store.PutPet(petID, pet); err != nil {
		return nilStatus, nil, APIError{
			InternalMessage: fmt.Sprintf("Could not put pet data. ID: %d, error: %v", petID, err),
			ResponseMessage: "Could not PUT new pet data",
			StatusCode:      http.StatusInternalServerError,
		}
	}
	return http.StatusOK, nil, nil
}

// DeletePet handles a DELETE request to delete a pet
func (ps *PetService) DeletePet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return nilStatus, nil, err
	}

	petDeleted, err := ps.store.DeletePet(petID)
	if err != nil {
		return nilStatus, nil, APIError{
			InternalMessage: fmt.Sprintf("Could not delete pet. ID: %d, error: %v", petID, err),
			ResponseMessage: "Could not delete pet data",
			StatusCode:      http.StatusInternalServerError,
		}
	}
	if petDeleted {
		return http.StatusOK, nil, nil
	}
	return http.StatusNoContent, nil, nil
}

func readPetID(r *http.Request) (int32, error) {
	petID, ok := r.Context().Value(keyID).(string)
	if !ok {
		return 0, keyExtractionError
	}
	if petID == "" {
		return 0, keyMissingError
	}
	intID, err := strconv.ParseInt(petID, 10, 32)
	if err != nil {
		return 0, invalidPetIDError(petID)
	}
	return int32(intID), nil
}

func readPetBody(r *http.Request) (*Pet, error) {
	if r.Body == nil {
		return nil, badRequestBodyError(errors.New("Missing Body"))
	}
	petData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, badRequestBodyError(err)
	}

	var pet Pet
	if err = json.Unmarshal(petData, &pet); err != nil {
		return nil, badPetDataError(err)
	}

	// TODO: validation

	return &pet, nil
}
