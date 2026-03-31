package interfaces

import (
	"auth_service/internal/domain"
)

type AuthService interface {
	Login(input LoginDTO) OutputDTO
	Registr(input RegistrDTO) OutputDTO
	DummyLogin(input DummyLoginDTO) OutputDTO
}

type AuthRepository interface {
	CreateUser(email, password, role string) (*domain.User, error)
	GetUser(email string) (*domain.User, error)
}

type JWTService interface {
	GenerateAccessToken(user *domain.User) (string, error)
}