package interfaces

import (
	"auth_service/internal/domain"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, input LoginDTO) OutputDTO
	Registr(ctx context.Context, input RegistrDTO) OutputDTO
	DummyLogin(ctx context.Context, input DummyLoginDTO) OutputDTO
}

type AuthRepository interface {
	CreateUser(ctx context.Context, email, password, role string) (*domain.User, error)
	GetUser(ctx context.Context, email string) (*domain.User, error)
}

type JWTService interface {
	GenerateAccessToken(user *domain.User) (string, error)
}