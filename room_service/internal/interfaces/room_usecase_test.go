package interfaces_test

import (
	"context"
	"errors"
	"testing"

	"room_service/internal/domain"
	"room_service/internal/interfaces"
)

type mockJWT struct {
	generateFunc func(user *domain.User) (string, error)
	validateFunc func(token string) (*domain.User, error)
}

func (m *mockJWT) GenerateAccessToken(user *domain.User) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(user)
	}
	return "mock-token", nil
}

func (m *mockJWT) ValidateToken(token string) (*domain.User, error) {
	if m.validateFunc != nil {
		return m.validateFunc(token)
	}
	return &domain.User{ID: "mock-user", Role: "user"}, nil
}

type mockRoomRepo struct {
	getAllFunc func(ctx context.Context)(*[]domain.Room,error)
	createFunc func(ctx context.Context,name,desc string,cap int)(*domain.Room,error)
	getFunc    func(ctx context.Context,id string)(*domain.Room,error)
}
func (m*mockRoomRepo)GetAllRooms(ctx context.Context)(*[]domain.Room,error){return m.getAllFunc(ctx)}
func (m*mockRoomRepo)CreateRoom(ctx context.Context,name,desc string,cap int)(*domain.Room,error){return m.createFunc(ctx,name,desc,cap)}
func (m*mockRoomRepo)GetRoom(ctx context.Context,id string)(*domain.Room,error){return m.getFunc(ctx,id)}

func TestRoomUsecase_ListRooms_Success(t *testing.T){
	repo:=&mockRoomRepo{getAllFunc:func(ctx context.Context)(*[]domain.Room,error){
		r:=[]domain.Room{{ID:"r1",Name:"Room1"}}
		return &r,nil
	}}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{Role:"admin"},nil}}
	uc:=interfaces.NewRoomUsecase(jwt,repo)
	out:=uc.ListRooms(context.Background(),"token")
	if out.Status!=200{t.Errorf("expected 200,got %d",out.Status)}
}

func TestRoomUsecase_ListRooms_InvalidToken(t *testing.T){
	repo:=&mockRoomRepo{}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return nil,errors.New("invalid")}}
	uc:=interfaces.NewRoomUsecase(jwt,repo)
	out:=uc.ListRooms(context.Background(),"badtoken")
	if out.Status!=401{t.Errorf("expected 401,got %d",out.Status)}
}

func TestRoomUsecase_CreateRoom_Admin_Success(t *testing.T){
	repo:=&mockRoomRepo{createFunc:func(ctx context.Context,name,desc string,cap int)(*domain.Room,error){
		return &domain.Room{ID:"new1",Name:name},nil
	}}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{Role:"admin"},nil}}
	uc:=interfaces.NewRoomUsecase(jwt,repo)
	out:=uc.CreateRoom(context.Background(),interfaces.CreateDTO{Name:"NewRoom"}, "token")
	if out.Status!=201{t.Errorf("expected 201,got %d",out.Status)}
}

func TestRoomUsecase_CreateRoom_NonAdmin(t *testing.T){
	repo:=&mockRoomRepo{}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{Role:"user"},nil}}
	uc:=interfaces.NewRoomUsecase(jwt,repo)
	out:=uc.CreateRoom(context.Background(),interfaces.CreateDTO{Name:"New"},"token")
	if out.Status!=403{t.Errorf("expected 403,got %d",out.Status)}
}