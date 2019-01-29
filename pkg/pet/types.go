package pet

// Pet defines the data structure corresponding to a pet
type Pet struct {
	ID      uint32                 `json:"id"`
	Name    string                 `json:"name"`
	Species string                 `json:"species"`
	Owner   string                 `json:"owner"`
	Extra   map[string]interface{} `json:"extra"`
}

// Storer defines standard CRUD operations for Pets
type Storer interface {
	CreatePet(*Pet) error
	ReadPet(ID uint32) (*Pet, error)
	UpdatePet(ID uint32, pet *Pet) error
	DeletePet(ID uint32) (bool, error)
}
