package interfaces

type RegistrDTO struct {
	email string `json:"email"`
	password string `json:"password"`
	role string `json:"role"`
}

type LoginDTO struct {
	email string `json:"email"`
	password string `json:"password"`
}

type DummyLoginDTO struct {
	role string `json:"role"`
}

type OutputDTO struct {
	Status int `json:"status"`
	Data map[string]interface{} `json:"data"`
}