package pet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"

	tassert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	pet1000 = func() Pet {
		return Pet{
			ID:      1000,
			Name:    "Nemo",
			Species: "Goldfish",
			Owner:   "Marlin",
			Extra:   nil,
		}
	}
	pet1001 = func() Pet {
		return Pet{
			ID:      1001,
			Name:    "Scruff",
			Species: "Golden Retriever",
			Owner:   "Bob Bobson",
			Extra:   map[string]interface{}{"food": "meat"},
		}
	}
	modifiedPet1000 = func() Pet {
		return Pet{
			ID:      1000,
			Name:    "Dory",
			Species: "Bluefish",
			Owner:   "Marlin",
			Extra:   nil,
		}
	}
)

type petServiceConfig struct {
	suite.Suite
	service *Service
	router  chi.Router
}

// resets the store for every test
func (p *petServiceConfig) SetupTest() {
	p.service.store = NewMemStore() // wipes any pet currently stored
}

func (p *petServiceConfig) TestGetPetSuccessful() {
	// given
	assert := tassert.New(p.T())
	testPet := pet1000()
	err := p.service.store.CreatePet(&testPet)
	assert.NoError(err, "Error initializing store. Check memstore errors")

	req, _ := http.NewRequest("GET", "/api/pet/1000", nil)
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusOK, resp.Code, "Response status should be 200 OK")
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Body should be readable")
	assert.NotEmpty(body, "Body should contain data")
	var retrievedPet Pet
	err = json.Unmarshal(body, &retrievedPet)
	assert.NoError(err, "Body should be able to unmarshal to a pet struct")
	assert.Equal(testPet, testPet, "Retrieved pet should be equal to test pet")
}

func (p *petServiceConfig) TestGetPet_NoPetExistsWithGivenID() {
	// given
	assert := tassert.New(p.T())
	req, _ := http.NewRequest("GET", "/api/pet/1000", nil)
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusNotFound, resp.Code, "Response status should be 404 Not Found")
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Body should be readable")
	assert.NotEmpty(body, "Body should be empty")
}

func (p *petServiceConfig) TestGetPet_InvalidID() {
	// given
	assert := tassert.New(p.T())
	req, _ := http.NewRequest("GET", "/api/pet/111x", nil)
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusBadRequest, resp.Code, "Response status should be 400 Bad Request")
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Body should be readable")
	assert.NotEmpty(body, "Body should be empty")
}

func (p *petServiceConfig) TestPostPetSuccessful() {
	// given
	assert := tassert.New(p.T())
	testPet := pet1000()

	requestBody, err := json.Marshal(testPet)
	if err != nil {
		panic(fmt.Errorf("Error in test code, could not marshal testPet to json. %v", err))
	}

	req, _ := http.NewRequest("POST", "/api/pet", bytes.NewBuffer(requestBody))
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusCreated, resp.Code, "Response status should be 201 Created")
	responseBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Body should be readable")
	assert.NotEmpty(responseBody, "Body should be empty")
}

func (p *petServiceConfig) TestPostPet_PetWithIDAlreadyExists() {
	// given
	assert := tassert.New(p.T())
	testPet := pet1000()
	if err := p.service.store.CreatePet(&testPet); err != nil {
		panic("Erro in test code, could not add initial data to test.")
	}

	requestBody, err := json.Marshal(testPet)
	if err != nil {
		panic(fmt.Errorf("Error in test code, could not marshal test pet to json. %v", err))
	}

	req, _ := http.NewRequest("POST", "/api/pet", bytes.NewBuffer(requestBody))
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusInternalServerError, resp.Code, "Response status should be 500 Internal Server Error")
	responseBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Body should be readable")
	assert.NotEmpty(responseBody, "Body should be empty")
}

func (p *petServiceConfig) TestPutPetSuccessful() {
	// given
	assert := tassert.New(p.T())
	testPet := pet1000()
	testModifiedPet := modifiedPet1000()
	if err := p.service.store.CreatePet(&testPet); err != nil {
		panic("Error in test code, could not add initial data to test.")
	}

	requestBody, err := json.Marshal(&testModifiedPet)
	if err != nil {
		panic(fmt.Errorf("Error in test code, could not marshal testModifiedPet to json. %v", err))
	}

	req, _ := http.NewRequest("PUT", "/api/pet/1000", bytes.NewBuffer(requestBody))
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusCreated, resp.Code, "Response status should be 201 Created")
	responseBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Response body should be readable")
	assert.NotEmpty(responseBody, "response body should not be empty")

	newPet, err := p.service.store.ReadPet(1000)
	assert.NoError(err, "Pet with id 1000 should be retrievable")
	assert.Equal(&testModifiedPet, newPet, "Pet with ID 1000 should be modified to be identical to testModifiedPet")
}

func (p *petServiceConfig) TestPutPet_CreatesNewPetIfNoneExist() {
	// given
	assert := tassert.New(p.T())
	testPet := pet1000()
	testNewPet := pet1001()
	if err := p.service.store.CreatePet(&testPet); err != nil {
		panic("Error in test code, could not add initial data to test")
	}

	requestBody, err := json.Marshal(&testNewPet)
	if err != nil {
		panic(fmt.Sprintf("Error in test code, could not marshal pet1001 to json. %v", err))
	}
	// Check no pet currently exists at id 1001
	if pet, _ := p.service.store.ReadPet(1001); pet != nil {
		panic("Error in test code, No pet with ID 1001 should exist")
	}

	req, _ := http.NewRequest("PUT", "/api/pet/1001", bytes.NewBuffer(requestBody))
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusCreated, resp.Code, "Response status should be 201 Created")
	responseBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Response body should be readable")
	assert.NotEmpty(responseBody, "response body should not be empty")

	newPet, err := p.service.store.ReadPet(1001)
	assert.NoError(err, "Pet with ID 1001 should be retrievable")
	assert.Equal(&testNewPet, newPet, "Pet with ID 1001 should be identical to pet in request payload")
}

func (p *petServiceConfig) TestDeletePetSuccessful() {
	// given
	assert := tassert.New(p.T())
	testPet := pet1000()
	if err := p.service.store.CreatePet(&testPet); err != nil {
		panic("Error in test code, Could not add initial data to test")
	}
	req, _ := http.NewRequest("DELETE", "/api/pet/1000", nil)
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusOK, resp.Code, "Response code should be 200 OK")
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Body should be readable")
	assert.NotEmpty(body, "Body should be empty")
}

func (p *petServiceConfig) TestDeletePet_NoPetExistsWithGivenID() {
	// given
	assert := tassert.New(p.T())
	req, _ := http.NewRequest("DELETE", "/api/pet/0", nil)
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusNoContent, resp.Code, "Response code should be 204 No Content")
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Body should be readable")
	assert.NotEmpty(body, "body should be empty")
}

func (p *petServiceConfig) TestDeletePet_InvalidID() {
	// given
	assert := tassert.New(p.T())
	req, _ := http.NewRequest("DELETE", "/api/pet/111x", nil)
	resp := httptest.NewRecorder()

	// when
	p.router.ServeHTTP(resp, req)

	// then
	assert.Equal(http.StatusBadRequest, resp.Code, "Response status should be 400 Bad Request")
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Body should be readable")
	assert.NotEmpty(body, "Body should be empty")
}

// Initialises config, run the suite
func TestPetService(t *testing.T) {
	store := NewMemStore()
	service := NewPetService(store)
	router := chi.NewRouter()
	SetupRoutes(router, service)
	config := &petServiceConfig{
		service: service,
		router:  router,
	}
	suite.Run(t, config)
}
