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

type SlotService interface {
	GetAllSlots(ctx context.Context, input GetAllSlotsDTO, token string) OutputDTO
}

type SlotRepository interface {
	GetAllSlotsByRoomIDAndDate(ctx context.Context, RoomID string, Date string) (*[]domain.Slot, error)
}

type BookingRepository interface {
	CreateBook(ctx context.Context, slot domain.Slot, userId string, slotId string) (*domain.Booking, error)
	GetAllBooksByUserID(ctx context.Context, userId string) (*[]domain.Booking, *[]time.Time, error)
	UpdateBookStatusByID(ctx context.Context, id string, userId string) (*domain.Booking, error)
	GetAllBookings(ctx context.Context, page, pageSize int) (*[]domain.Booking, int, error)
}

type BookingService interface {
	CancelBook(ctx context.Context, input CancelBookingDTO, token string) OutputDTO
	CreateBook(ctx context.Context, input CreateBookingDTO, token string) OutputDTO
	GetMyBooks(ctx context.Context, token string) OutputDTO
	GetAllBookings(ctx context.Context, input GetAllBooksDTO, token string) OutputDTO
}

type JWTService interface {
	ValidateToken(token string) (*domain.User, error)
}