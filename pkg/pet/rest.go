package pet

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type contextKey int

const (
	idKey contextKey = iota
)

// maps from internal errors to response status codes
// renderErrorResponse defaults to internal server error
// if a specific error code is not defined.
var errStatusMap = map[int]int{
	ErrInvalidInput: http.StatusBadRequest,
	ErrNotFound:     http.StatusNotFound,
}

// renderErrorResponse handles http responses in the case of an error
func renderErrorResponse(w http.ResponseWriter, err error) {
	message := err.Error()
	responseStatus := http.StatusInternalServerError
	// pet service Errors store more specific response information
	if specificError, ok := err.(*Error); ok {
		message = specificError.Message
		// Attempt to get a more specific status code
		if status, ok := errStatusMap[specificError.Code]; ok {
			responseStatus = status
		}
	}
	http.Error(w, message, responseStatus)
}

// urlParamContextSaverMiddleware is a middleware that extracts a url parameter on an path
// and saves its value to the request context for downstream paths and endpoints
func urlParamContextSaverMiddleware(urlParam string, id contextKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			up := chi.URLParam(r, urlParam)
			ctx := context.WithValue(r.Context(), id, up)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SetupRoutes sets up pet service routes for the given router
func SetupRoutes(r chi.Router, s *Service) {
	r.Route("/api/pet", func(r chi.Router) {
		r.Post("/", s.PostPet)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(urlParamContextSaverMiddleware("id", idKey))
			r.Get("/", s.GetPet)
			r.Put("/", s.PutPet)
			r.Delete("/", s.DeletePet)
		})
	})
}

// Service defines a rest api for interaction with a PetStorer
type Service struct {
	store Storer
}

// NewPetService creates a new pet service with an in-memory store
func NewPetService(storer Storer) *Service {
	return &Service{
		store: storer,
	}
}

// These functions take a request and return the appropriate response and status code
// In the case these return an error, the error will be passed

// GetPet handles a GET request to retrieve a pet
func (ps *Service) GetPet(w http.ResponseWriter, r *http.Request) {
	petID, err := readPetID(r)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}
	pet, err := ps.store.ReadPet(petID)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, pet)
}

// PostPet handles a POST request to add a new pet
func (ps *Service) PostPet(w http.ResponseWriter, r *http.Request) {
	newPet, err := readPetBody(r)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}

	if err = ps.store.CreatePet(newPet); err != nil {
		renderErrorResponse(w, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, nil)
}

// PutPet handles a PUT request to create or modify a pet
func (ps *Service) PutPet(w http.ResponseWriter, r *http.Request) {
	petID, err := readPetID(r)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}
	pet, err := readPetBody(r)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}
	if err = ps.store.UpdatePet(petID, pet); err != nil {
		renderErrorResponse(w, err)
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, nil)
}

// DeletePet handles a DELETE request to delete a pet
func (ps *Service) DeletePet(w http.ResponseWriter, r *http.Request) {
	petID, err := readPetID(r)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}
	petDeleted, err := ps.store.DeletePet(petID)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}
	if petDeleted {
		render.Status(r, http.StatusOK)
	} else {
		render.Status(r, http.StatusNoContent)
	}
	render.JSON(w, r, nil)
}

func readPetID(r *http.Request) (uint32, error) {
	petID := r.Context().Value(idKey)
	if petID == nil {
		// Reaching this indicates a bug. At this point, request context should contain an id
		return uint32(ErrUnknown), Errorf(ErrUnknown, "pet ID was lost somewhere")
	}
	intID, err := strconv.ParseInt(petID.(string), 10, 32)
	if err != nil {
		return uint32(ErrInvalidInput), Errorf(ErrInvalidInput, "Invalid pet ID %v. ID should be a number", petID)
	}
	return uint32(intID), nil
}

func readPetBody(r *http.Request) (*Pet, error) {
	if r.Body == nil {
		return nil, Errorf(ErrInvalidInput, "No request body")
	}
	petData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrorEf(ErrInvalidInput, err, "Bad request body")
	}
	var pet Pet
	if err = json.Unmarshal(petData, &pet); err != nil {
		return nil, ErrorEf(ErrInvalidInput, err, "Invalid pet data")
	}
	return &pet, nil
}
