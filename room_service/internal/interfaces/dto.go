package interfaces


type CreateDTO struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Capacity int `json:"capacity"`
}

type OutputDTO struct {
	Status int `json:"status"`
	Data map[string]interface{} `json:"data"`
}

type GetDTO struct {
	Name string `json:"name"`
}