package interfaces


type CreateDTO struct {
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Capacity int `json:"capacity,omitempty"`
}

type OutputDTO struct {
	Status int `json:"status"`
	Data map[string]interface{} `json:"data"`
}

type GetDTO struct {
	ID string `json:"room_id"`
}