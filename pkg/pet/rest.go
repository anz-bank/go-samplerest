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
	r.Route("/api", func(r chi.Router) {
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
			// TODO: Logging
			renderHTTPErrorResponse(w, r, err)
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

// nilStatus is used when it is expected to be
// overriden in the response
// eg: function is returning an error
const nilStatus = http.StatusInternalServerError

// GetPet handles a GET request to retrieve a pet
func (ps *Service) GetPet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return nilStatus, nil, err
	}
	pet, err := ps.store.GetPet(petID)
	if err != nil {
		return nilStatus, nil, err
	}
	return http.StatusOK, pet, nil
}

// PostPet handles a POST request to add a new pet
func (ps *Service) PostPet(r *http.Request) (int, interface{}, error) {
	newPet, err := readPetBody(r)
	if err != nil {
		return nilStatus, nil, err
	}

	if err = ps.store.PostPet(newPet); err != nil {
		return nilStatus, nil, err
	}

	return http.StatusCreated, nil, nil
}

// PutPet handles a PUT request to modify a pet
func (ps *Service) PutPet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return nilStatus, nil, err
	}
	pet, err := readPetBody(r)
	if err != nil {
		return nilStatus, nil, err
	}

	if err = ps.store.PutPet(petID, pet); err != nil {
		return nilStatus, nil, err
	}
	return http.StatusOK, nil, nil
}

// DeletePet handles a DELETE request to delete a pet
func (ps *Service) DeletePet(r *http.Request) (int, interface{}, error) {
	petID, err := readPetID(r)
	if err != nil {
		return nilStatus, nil, err
	}

	petDeleted, err := ps.store.DeletePet(petID)
	if err != nil {
		return nilStatus, nil, err
	}
	if petDeleted {
		return http.StatusOK, nil, nil
	}
	return http.StatusNoContent, nil, nil
}

func readPetID(r *http.Request) (int32, error) {
	petID, ok := r.Context().Value(idKey).(string)
	if !ok {
		return 0, APIErrorf(ErrUnknown, "pet ID was lost somewhere")
	}
	if petID == "" {
		return 0, APIErrorf(ErrUnknown, "pet ID was empty")
	}
	intID, err := strconv.ParseInt(petID, 10, 32)
	if err != nil {
		return 0, APIErrorEf(0, err, "Inalid pet ID %d. ID should be a number", intID)
	}
	return int32(intID), nil
}

func readPetBody(r *http.Request) (*Pet, error) {
	if r.Body == nil {
		return nil, APIErrorf(ErrBadRequest, "No request body")
	}
	petData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, APIErrorEf(ErrBadRequest, err, "Bad request body")
	}

	var pet Pet
	if err = json.Unmarshal(petData, &pet); err != nil {
		return nil, APIErrorEf(ErrBadRequest, err, "Invalid pet data")
	}

	// TODO: validation

	return &pet, nil
}
