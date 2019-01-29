package pet

import (
	"sync"
)

// MemStore is an in-memory implementation of PetStorer
type MemStore struct {
	sync.Mutex
	sync.Map
}

// NewMemStore creates a new in-memory store with map intialised
func NewMemStore() *MemStore {
	return &MemStore{}
}

// CreatePet adds a new pet to the store
func (m *MemStore) CreatePet(pet *Pet) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Load(pet.ID); ok {
		return Errorf(ErrDuplicate, "Pet with id %d already exists", pet.ID)
	}
	m.Store(pet.ID, *pet)
	return nil
}

// ReadPet gets a pet from the store given an ID
func (m *MemStore) ReadPet(petID uint32) (*Pet, error) {
	petData, ok := m.Load(uint32(petID))
	if !ok {
		return nil, Errorf(ErrNotFound, "No pet exists with id %d", petID)
	}
	pet, ok := petData.(Pet)
	if !ok {
		return nil, Errorf(ErrUnknown, "Could not read pet data from store")
	}
	return &pet, nil
}

// UpdatePet puts new pet data to the store, either creating a new one or overriding an old
func (m *MemStore) UpdatePet(petID uint32, pet *Pet) error {
	m.Store(petID, *pet)
	return nil
}

// DeletePet deletes a pet from the store
func (m *MemStore) DeletePet(petID uint32) (bool, error) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Load(petID); !ok {
		return false, nil
	}
	m.Delete(petID)
	return true, nil
}
