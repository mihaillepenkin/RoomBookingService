package interfaces

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	jwtRepo  JWTService
	authRepo AuthRepository
}

func NewAuthUsecase(repo AuthRepository, jwt JWTService) *AuthUsecase {
	return &AuthUsecase{jwtRepo: jwt, authRepo: repo}
}

func (a *AuthUsecase) Login(ctx context.Context, input LoginDTO) OutputDTO {
	user, err := a.authRepo.GetUser(ctx, input.Email)
	if err != nil {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": err.Error()}}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "invalid email or password"}}
	}
	token, err := a.jwtRepo.GenerateAccessToken(user)
	if err != nil {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"token": token}}
}

func (a *AuthUsecase) Registr(ctx context.Context, input RegistrDTO) OutputDTO {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	user, err := a.authRepo.CreateUser(ctx, input.Email, string(hashed), input.Role)
	if err != nil {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": err.Error()}}
	}
	//чтобы не выводить хешированный
	user.Password = input.Password
	return OutputDTO{Status: 200, Data: map[string]interface{}{"user": *user}}
}

func (a *AuthUsecase) DummyLogin(ctx context.Context, input DummyLoginDTO) OutputDTO {
	testEmail := "misha1@example.com"
	password := "qwerty123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	switch input.Role {
	case "admin":
		testEmail = "misha1@example.com"
	case "user":
		testEmail = "misha2@example.com"
	default:
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": "invalid role for DummyLogin"}}
	}
	user, err := a.authRepo.GetUser(ctx, testEmail)
	if err != nil {
		user, err = a.authRepo.CreateUser(ctx, testEmail, string(hashed), input.Role)
		if (err != nil) {
			return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
		}
	}
	token, err := a.jwtRepo.GenerateAccessToken(user)
	if err != nil {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"token": token}}
}