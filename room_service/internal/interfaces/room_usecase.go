package interfaces

import (
	"context"
	"fmt"
)

type RoomUsecase struct {
	jwt  JWTService
	repo RoomRepository
}

func NewRoomUsecase(jwt JWTService, repo RoomRepository) *RoomUsecase {
	return &RoomUsecase{jwt: jwt, repo: repo}
}

func (r *RoomUsecase) ListRooms(ctx context.Context, token string) OutputDTO {
	_, err := r.jwt.ValidateToken(token)
	if (err != nil) {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "invalid token"}}
	}
	list, err := r.repo.GetAllRooms(ctx)
	if err != nil {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"list": *list}}
}

func (r *RoomUsecase) CreateRoom(ctx context.Context, input CreateDTO, token string) OutputDTO {
	user, err := r.jwt.ValidateToken(token)
	if err != nil {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "invalid token"}}
	}
	if user.Role != "admin" {
		return OutputDTO{Status: 403, Data: map[string]interface{}{"error": fmt.Sprintf("for this action need admin role, your role : %s", user.Role)}}
	}
	room, err := r.repo.CreateRoom(ctx, input.Name, input.Description, input.Capacity)
	if err != nil {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 201, Data: map[string]interface{}{"room": *room}}
}

func (r *RoomUsecase) GetRoom(ctx context.Context, input GetDTO) OutputDTO {
	room, err := r.repo.GetRoom(ctx, input.ID)
	if err != nil {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"room": *room}}
}