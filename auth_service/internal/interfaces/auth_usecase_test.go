package interfaces_test

import (
	"context"
	"errors"
	"testing"

	"auth_service/internal/domain"
	"auth_service/internal/interfaces"

	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepo struct {
	getUserFunc func(ctx context.Context, email string) (*domain.User, error)
	createFunc  func(ctx context.Context, email, password, role string) (*domain.User, error)
}

func (m *mockAuthRepo) GetUser(ctx context.Context, email string) (*domain.User, error) {
	return m.getUserFunc(ctx, email)
}
func (m *mockAuthRepo) CreateUser(ctx context.Context, email, password, role string) (*domain.User, error) {
	return m.createFunc(ctx, email, password, role)
}

type mockJWT struct {
	generateFunc func(user *domain.User) (string, error)
	validateFunc func(token string) (*domain.User, error)
}

func (m *mockJWT) GenerateAccessToken(user *domain.User) (string, error) {
	return m.generateFunc(user)
}
func (m *mockJWT) ValidateToken(token string) (*domain.User, error) {
	return m.validateFunc(token)
}

func TestAuthUsecase_Login_Success(t *testing.T){
	repo:=&mockAuthRepo{
		getUserFunc:func(ctx context.Context,email string)(*domain.User,error){
			return &domain.User{ID:"u1",Email:email,Password:hashForTest("pass123"),Role:"user"},nil
		},
	}
	jwt:=&mockJWT{
		generateFunc:func(user *domain.User)(string,error){return "token123",nil},
	}
	uc:=interfaces.NewAuthUsecase(repo,jwt)
	out:=uc.Login(context.Background(),interfaces.LoginDTO{Email:"test@example.com",Password:"pass123"})
	if out.Status!=200{t.Errorf("expected 200,got %d",out.Status)}
	if out.Data["token"]!="token123"{t.Errorf("expected token123,got %v",out.Data["token"])}
}

func hashForTest(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost) // MinCost для скорости
	return string(h)
}

func TestAuthUsecase_Login_WrongPassword(t *testing.T){
	repo:=&mockAuthRepo{
		getUserFunc:func(ctx context.Context,email string)(*domain.User,error){
			return &domain.User{ID:"u1",Email:email,Password:"$2a$10$WrongHashWrongHashWrongHashWrongHashWrongHashWrongHashWro",Role:"user"},nil
		},
	}
	jwt:=&mockJWT{}
	uc:=interfaces.NewAuthUsecase(repo,jwt)
	out:=uc.Login(context.Background(),interfaces.LoginDTO{Email:"test@example.com",Password:"wrongpass"})
	if out.Status!=401{t.Errorf("expected 401,got %d",out.Status)}
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T){
	repo:=&mockAuthRepo{
		getUserFunc:func(ctx context.Context,email string)(*domain.User,error){
			return nil,errors.New("user not found")
		},
	}
	jwt:=&mockJWT{}
	uc:=interfaces.NewAuthUsecase(repo,jwt)
	out:=uc.Login(context.Background(),interfaces.LoginDTO{Email:"notfound@example.com",Password:"pass"})
	if out.Status!=401{t.Errorf("expected 401,got %d",out.Status)}
}

func TestAuthUsecase_DummyLogin_Admin(t *testing.T){
	repo:=&mockAuthRepo{
		getUserFunc:func(ctx context.Context,email string)(*domain.User,error){
			return nil,errors.New("not found")
		},
		createFunc:func(ctx context.Context,email,password,role string)(*domain.User,error){
			return &domain.User{ID:"admin1",Email:email,Role:role},nil
		},
	}
	jwt:=&mockJWT{generateFunc:func(u *domain.User)(string,error){return "admintoken",nil}}
	uc:=interfaces.NewAuthUsecase(repo,jwt)
	out:=uc.DummyLogin(context.Background(),interfaces.DummyLoginDTO{Role:"admin"})
	if out.Status!=200{t.Errorf("expected 200,got %d",out.Status)}
	if out.Data["token"]!="admintoken"{t.Errorf("expected admintoken")}
}

func TestAuthUsecase_DummyLogin_InvalidRole(t *testing.T){
	repo:=&mockAuthRepo{}
	jwt:=&mockJWT{}
	uc:=interfaces.NewAuthUsecase(repo,jwt)
	out:=uc.DummyLogin(context.Background(),interfaces.DummyLoginDTO{Role:"superadmin"})
	if out.Status!=400{t.Errorf("expected 400,got %d",out.Status)}
}