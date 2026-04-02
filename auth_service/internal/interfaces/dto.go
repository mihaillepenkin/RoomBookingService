package interfaces

type RegistrDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
}

type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type DummyLoginDTO struct {
	Role string `json:"role"`
}

type OutputDTO struct {
	Status int `json:"status"`
	Data map[string]interface{} `json:"data"`
}