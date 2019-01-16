package pet

import (
	"sync"
)

// MemStore is an in-memory implementation of PetStorer
type MemStore struct {
	mutex *sync.Mutex
	store sync.Map
}

// NewMemStore creates a new in-memory store with map intialised
func NewMemStore() *MemStore {
	return &MemStore{
		mutex: &sync.Mutex{},
		store: sync.Map{},
	}
}

// GetPet gets a pet from the store given an ID
func (m *MemStore) GetPet(petID int32) (*Pet, error) {
	petData, ok := m.store.Load(petID)
	if !ok {
		return nil, Errorf(ErrIDNotFound, "No pet exists with id %d", petID)
	}
	pet, ok := petData.(Pet)
	if !ok {
		return nil, Errorf(ErrUnknown, "Could not read pet data from store")
	}
	return &pet, nil
}

// PostPet adds a new pet to the store
func (m *MemStore) PostPet(pet *Pet) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.store.Load(pet.ID); ok {
		return Errorf(ErrIDAlreadyExists, "Pet with id %d already exists", pet.ID)
	}
	m.store.Store(pet.ID, *pet)
	return nil
}

// PutPet puts new pet data to the store, either creating a new one or overriding an old
func (m *MemStore) PutPet(petID int32, pet *Pet) error {
	m.store.Store(petID, *pet)
	return nil
}

// DeletePet deletes a pet from the store
func (m *MemStore) DeletePet(petID int32) (bool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.store.Load(petID); !ok {
		return false, nil
	}
	m.store.Delete(petID)
	return true, nil
}
