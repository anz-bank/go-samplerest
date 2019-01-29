package pet

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type storerImpl int

const (
	mem storerImpl = iota
)

var (
	pet10 = func() Pet {
		return Pet{
			ID:      10,
			Name:    "Slinky",
			Species: "Toy dog",
			Owner:   "Andy",
			Extra:   nil,
		}
	}
	modifiedPet10 = func() Pet {
		return Pet{
			ID:      10,
			Name:    "Mr Potato Head",
			Species: "Potato Head",
			Owner:   "Andy",
			Extra:   map[string]interface{}{"temperament": "aggressive"},
		}
	}
	pet11 = func() Pet {
		return Pet{
			ID:      11,
			Name:    "Bo Peep",
			Species: "Sheep herder",
			Owner:   "Molly",
			Extra:   map[string]interface{}{"likes": "Woody"},
		}
	}
)

type storerSuite struct {
	suite.Suite
	store Storer
	impl  storerImpl
}

func (s *storerSuite) SetupTest() {
	switch s.impl {
	case mem:
		// create a fresh memstore
		s.store = NewMemStore()
		pet := pet10()
		s.store.CreatePet(&pet)
	default:
		panic("Unrecognised storer implementation")
	}
}

func (s *storerSuite) TestReadPetSuccessful() {
	// given
	assert := tassert.New(s.T())

	// when
	pet, err := s.store.ReadPet(10)

	// then
	expectedPet := pet10()
	if assert.NoError(err, "Should be able to read pet10 from store") {
		assert.Equal(&expectedPet, pet, "store should return a pet identical to pet10")
	}
}

func (s *storerSuite) TestReadPet_IDDoesNotExist() {
	// given
	assert := tassert.New(s.T())

	// when
	_, err := s.store.ReadPet(11)

	// then
	assert.Error(err, "Should get an error when attempting to read an non-existing pet")
}

func (s *storerSuite) TestCreatePetSuccessful() {
	// given
	assert := tassert.New(s.T())
	newPet := pet11()

	// when
	err := s.store.CreatePet(&newPet)

	// then
	assert.NoError(err, "Should not get an error creating a pet to a free ID")
	createdPet, err := s.store.ReadPet(11)
	if assert.NoError(err, "Should be able to read an newly created pet") {
		assert.Equal(&newPet, createdPet, "Created pet should be identical to the one passed to Create")
	}
}

func (s *storerSuite) TestCreatePet_IDAlreadyTaken() {
	// given
	assert := tassert.New(s.T())
	oldPet := pet10()
	newPet := modifiedPet10()

	// when
	err := s.store.CreatePet(&newPet)

	// then
	assert.Error(err, "Create should return an error if attempting to create to an already existing ID")
	currentPet, err := s.store.ReadPet(10)
	if assert.NoError(err, "Should be able to read the old pet after a failed overwrite attempt") {
		assert.Equal(&oldPet, currentPet, "Stored pet should be identical to the old pet after a failed overwrite attempt")
	}
}

func (s *storerSuite) TestUpdatePetSuccessful() {
	// given
	assert := tassert.New(s.T())
	testModifiedPet := modifiedPet10()

	// when
	err := s.store.UpdatePet(10, &testModifiedPet)

	// then
	assert.NoError(err, "UpdatePet should successfully update a pet")
	storedPet, err := s.store.ReadPet(10)
	if assert.NoError(err, "Should be able to read a modified pet") {
		assert.Equal(&testModifiedPet, storedPet, "Stored pet should be equal to the modified pet")
	}
}

func (s *storerSuite) TestUpdatePet_IDDoesNotExist() {
	// given
	assert := tassert.New(s.T())
	testPet := pet11()

	// when
	err := s.store.UpdatePet(11, &testPet)

	// then
	assert.NoError(err, "Updating to non-existing pet ID is not an error")
	newPet, err := s.store.ReadPet(11)
	if assert.NoError(err, "Should be able to read a newly added pet via update") {
		assert.Equal(&testPet, newPet, "Newly added pet should be equal to test pet")
	}
}

func (s *storerSuite) TestDeletePetSuccessful() {
	// when
	assert := tassert.New(s.T())
	deleted, err := s.store.DeletePet(10)

	// then
	assert.NoError(err, "Delete should successfully delete a pet")
	assert.True(deleted, "Delete should return true indicating a pet was deleted")
	_, err = s.store.ReadPet(10)
	assert.Error(err, "Should not be able to read a deleted ID")
}

func (s *storerSuite) TestDeletePet_IDDoesNotExist() {
	// when
	assert := tassert.New(s.T())
	deleted, err := s.store.DeletePet(11)

	// then
	assert.NoError(err, "Deleting a non-existing ID is not an error")
	assert.False(deleted, "Delete should return false indicating no pet was deleted")
}

func TestStorer(t *testing.T) {
	memSuite := storerSuite{
		store: NewMemStore(),
		impl:  mem,
	}
	suite.Run(t, &memSuite)
}
