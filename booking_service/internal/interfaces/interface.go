package interfaces

import (
	"booking_service/internal/domain"
	"context"
	"time"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, input CreateScheduleDTO, token string) OutputDTO
}

type ScheduleRepository interface {
	CreateSchedule(ctx context.Context, RoomID string, DayOfWeek []int64, startTime *time.Time, endTime *time.Time) (*domain.Schedule, error)
	GetScheduleByRoomID(ctx context.Context, RoomID string) (*domain.Schedule, error)
}

type JWTService interface {
	ValidateToken(token string) (*domain.User, error)
}