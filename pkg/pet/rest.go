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
// renderHTTPErrorResponse defaults to internal server error
// if a specific error code is not defined.
var errStatusMap = map[int]int{
	ErrBadRequest: http.StatusBadRequest,
	ErrIDNotFound: http.StatusNotFound,
}

// renderHTTPErrorResponse handles http responses in the case of an error
func renderHTTPErrorResponse(w http.ResponseWriter, err error) {
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
	return
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
		r.Post("/", makeHandler(s.PostPet))
		r.Route("/{id}", func(r chi.Router) {
			r.Use(urlParamContextSaverMiddleware("id", idKey))
			r.Get("/", makeHandler(s.GetPet))
			r.Put("/", makeHandler(s.PutPet))
			r.Delete("/", makeHandler(s.DeletePet))
		})
	})
}

func makeHandler(serveRequest func(*http.Request) (int, interface{}, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode, response, err := serveRequest(r)
		if err != nil {
			renderHTTPErrorResponse(w, err)
			return
		}
		render.Status(r, statusCode)
		render.JSON(w, r, response)
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

// errorResponseData returns the response data in the case an error occurs
// The handler functions are expected to return either
// 1. a response code and body
// 2. an error
// In the case an error occurs, the response status and body will be overriden
// by the error handler
func errorResponseData(err error) (int, interface{}, error) {
	return 0, nil, err
}

// These functions take a request and return the appropriate response and status code
// In the case these return an error, the error will be passed

// GetPet handles a GET request to retrieve a pet
func (ps *Service) GetPet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return errorResponseData(err)
	}
	pet, err := ps.store.ReadPet(petID)
	if err != nil {
		return errorResponseData(err)
	}
	return http.StatusOK, pet, nil
}

// PostPet handles a POST request to add a new pet
func (ps *Service) PostPet(r *http.Request) (int, interface{}, error) {
	newPet, err := readPetBody(r)
	if err != nil {
		return errorResponseData(err)
	}

	if err = ps.store.CreatePet(newPet); err != nil {
		return errorResponseData(err)
	}

	return http.StatusCreated, nil, nil
}

// PutPet handles a PUT request to create or modify a pet
func (ps *Service) PutPet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return errorResponseData(err)
	}
	pet, err := readPetBody(r)
	if err != nil {
		return errorResponseData(err)
	}
	if err = ps.store.UpdatePet(petID, pet); err != nil {
		return errorResponseData(err)
	}
	return http.StatusCreated, nil, nil
}

// DeletePet handles a DELETE request to delete a pet
func (ps *Service) DeletePet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return errorResponseData(err)
	}
	petDeleted, err := ps.store.DeletePet(petID)
	if err != nil {
		return errorResponseData(err)
	}
	if petDeleted {
		return http.StatusOK, nil, nil
	}
	return http.StatusNoContent, nil, nil
}

func readPetID(r *http.Request) (uint32, error) {
	petID := r.Context().Value(idKey)
	if petID == nil {
		// Reaching this indicates a bug. At this point, request context should contain an id
		return uint32(ErrUnknown), Errorf(ErrUnknown, "pet ID was lost somewhere")
	}
	intID, err := strconv.ParseInt(petID.(string), 10, 32)
	if err != nil {
		return uint32(ErrBadRequest), Errorf(ErrBadRequest, "Invalid pet ID %v. ID should be a number", petID)
	}
	return uint32(intID), nil
}

func readPetBody(r *http.Request) (*Pet, error) {
	if r.Body == nil {
		return nil, Errorf(ErrBadRequest, "No request body")
	}
	petData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrorEf(ErrBadRequest, err, "Bad request body")
	}
	var pet Pet
	if err = json.Unmarshal(petData, &pet); err != nil {
		return nil, ErrorEf(ErrBadRequest, err, "Invalid pet data")
	}
	return &pet, nil
}
