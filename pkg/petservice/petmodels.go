package petservice

// Pet defines the data structure corresponding to a pet
type Pet struct {
	ID      int32                  `json:"id"`
	Name    string                 `json:"name"`
	Species string                 `json:"owner"`
	Owner   string                 `json:"species"`
	Extra   map[string]interface{} `json:"extra"`
}
