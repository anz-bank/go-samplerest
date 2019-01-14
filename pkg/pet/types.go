package pet

// Pet defines the data structure corresponding to a pet
type Pet struct {
	ID      int32                  `json:"id"`
	Name    string                 `json:"name"`
	Species string                 `json:"owner"`
	Owner   string                 `json:"species"`
	Extra   map[string]interface{} `json:"extra"`
}

// Storer defines standard CRUD operations for Pets
type Storer interface {
	GetPet(ID int32) (*Pet, error)
	PostPet(*Pet) error
	PutPet(ID int32, pet *Pet) error
	DeletePet(ID int32) (bool, error)
}
