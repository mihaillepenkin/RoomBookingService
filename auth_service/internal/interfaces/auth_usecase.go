package interfaces

import "golang.org/x/crypto/bcrypt"

type AuthUsecase struct {
	jwtRepo  JWTService
	authRepo AuthRepository
}

func NewAuthUsecase(repo AuthRepository, jwt JWTService) *AuthUsecase {
	return &AuthUsecase{jwtRepo: jwt, authRepo: repo}
}

func (a *AuthUsecase) Login(input LoginDTO) OutputDTO {
	user, err := a.authRepo.GetUser(input.email)
	if err != nil {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": err.Error()}}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.password)); err != nil {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": err.Error()}}
	}
	token, err := a.jwtRepo.GenerateAccessToken(user)
	if err != nil {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"token": token}}
}

func (a *AuthUsecase) Registr(input RegistrDTO) OutputDTO {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.password), bcrypt.DefaultCost)
	if err != nil {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	user, err := a.authRepo.CreateUser(input.email, string(hashed), input.role)
	if err != nil {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"user": *user}}
}

func (a *AuthUsecase) DummyLogin(input DummyLoginDTO) OutputDTO {
	testEmail := "misha1@example.com"
	switch input.role {
	case "admin":
		testEmail = "misha1@example.com"
	case "user":
		testEmail = "misha2@example.com"
	default:
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": "invalid role for DummyLogin"}}
	}
	user, err := a.authRepo.GetUser(testEmail)
	if err != nil {
		//здесь 500 код, так как эти пользователи должны в миграции добавиться сразу
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	token, err := a.jwtRepo.GenerateAccessToken(user)
	if err != nil {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"token": token}}
}