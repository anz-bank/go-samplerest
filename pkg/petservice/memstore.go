package petservice

import (
	"fmt"
)

// MemStore is an in-memory storage that implements PetStorer
// NOT persistent
type MemStore struct {
	store map[int32]Pet
}

// NewMemStore creates a new in-memory store with initialised map
func NewMemStore() *MemStore {
	return &MemStore{
		store: make(map[int32]Pet),
	}
}

// GetPet retrieves a pet from the store with the given ID
// Returns an error if none exist
func (m *MemStore) GetPet(petID int32) (*Pet, error) {
	pet, ok := m.store[petID]
	if !ok {
		return nil, fmt.Errorf("No pet exists with id %d", petID)
	}
	return &pet, nil
}

// PostPet creates a new pet with the given ID
// TODO: have the store create an ID
func (m *MemStore) PostPet(pet *Pet) error {
	_, ok := m.store[pet.ID]
	if ok {
		return fmt.Errorf("Pet with id %d already exists", pet.ID)
	}
	m.store[pet.ID] = *pet
	return nil
}

// PutPet replaces the pet with the given ID
// If the pet does not exist, will create it
// Currently no rules about replacing pet data
func (m *MemStore) PutPet(petID int32, pet *Pet) error {
	m.store[pet.ID] = *pet
	return nil
}

// DeletePet deletes the pet for a given id
// if the id does not exists, does nothing
// bool return value indicates whether a pet was deleted
func (m *MemStore) DeletePet(petID int32) (bool, error) {
	_, ok := m.store[petID]
	if !ok {
		return false, nil
	}
	delete(m.store, petID)
	return true, nil
}
