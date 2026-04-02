package interfaces

import (
	"context"
	"room_service/internal/domain"
)

type RoomService interface {
	ListRooms(ctx context.Context, token string) OutputDTO
	CreateRoom(ctx context.Context, input CreateDTO, token string) OutputDTO
	GetRoom(ctx context.Context, input GetDTO) OutputDTO
}

type RoomRepository interface {
	GetAllRooms(ctx context.Context) (*[]domain.Room, error)
	CreateRoom(ctx context.Context, name, description string, capacity int) (*domain.Room, error)
	GetRoom(ctx context.Context, name string) (*domain.Room, error)
}

type JWTService interface {
	ValidateToken(token string)  (*domain.User, error)
}