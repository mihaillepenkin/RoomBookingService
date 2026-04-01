package interfaces

import (
	"room_service/internal/domain"
)

type RoomService interface {
	ListRooms() OutputDTO
	CreateRoom(input CreateDTO, token string) OutputDTO
	GetRoom(input GetDTO) OutputDTO
}

type RoomRepository interface {
	GetAllRooms() (*[]domain.Room, error)
	CreateRoom(name, description string, capacity int) (*domain.Room, error)
	GetRoom(name string) (*domain.Room, error)
}

type JWTService interface {
	ValidateToken(token string)  (*domain.User, error)
}