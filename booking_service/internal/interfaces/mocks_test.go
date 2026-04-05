package interfaces_test

import (
	"booking_service/internal/domain"
	"context"
	"database/sql"
	"time"
)


type mockBookingRepo struct {
	createFunc    func(ctx context.Context, slot domain.Slot, userID, slotID string) (*domain.Booking, error)
	getByUserFunc func(ctx context.Context, userID string) (*[]domain.Booking, *[]time.Time, error)
	updateFunc    func(ctx context.Context, id, userID string) (*domain.Booking, error)
	getAllFunc    func(ctx context.Context, page, pageSize int) (*[]domain.Booking, int, error)
	getBookedFunc func(ctx context.Context, roomID, date string) (*[]domain.Slot, error)
}

type mockJWT struct {
	generateFunc func(*domain.User) (string, error)
	validateFunc func(string) (*domain.User, error)
}
func (m *mockJWT) GenerateAccessToken(u *domain.User) (string, error) {
	if m.generateFunc != nil { return m.generateFunc(u) }
	return "mock-token", nil
}
func (m *mockJWT) ValidateToken(t string) (*domain.User, error) {
	if m.validateFunc != nil { return m.validateFunc(t) }
	return &domain.User{ID:"mock",Role:"user"}, nil
}
func (m *mockBookingRepo) GetAllSlotsByRoomIDAndDate(ctx context.Context, roomID, date string) (*[]domain.Slot, error) {
	if m.getBookedFunc != nil { return m.getBookedFunc(ctx, roomID, date) }
	s := []domain.Slot{}
	return &s, nil
}
func (m *mockBookingRepo) CreateBook(ctx context.Context, slot domain.Slot, userID, slotID string) (*domain.Booking, error) { return nil, nil }
func (m *mockBookingRepo) GetAllBooksByUserID(ctx context.Context, userID string) (*[]domain.Booking, *[]time.Time, error) { return nil, nil, nil }
func (m *mockBookingRepo) UpdateBookStatusByID(ctx context.Context, id, userID string) (*domain.Booking, error) { return nil, sql.ErrNoRows }
func (m *mockBookingRepo) GetAllBookings(ctx context.Context, page, pageSize int) (*[]domain.Booking, int, error) { return nil, 100, nil }
type mockScheduleRepo struct {
	getFunc    func(ctx context.Context, roomID string) (*domain.Schedule, error)
	createFunc func(ctx context.Context, roomID string, days []int64, start, end *time.Time) (*domain.Schedule, error)
}
func (m *mockScheduleRepo) GetScheduleByRoomID(ctx context.Context, roomID string) (*domain.Schedule, error) {
	if m.getFunc != nil { return m.getFunc(ctx, roomID) }
	return nil, nil
}
type mockRoomClient struct {
	checkFunc func(context.Context, string) (bool, error)
}
func (m *mockRoomClient) CheckRoomExists(ctx context.Context, id string) (bool, error) {
	if m.checkFunc != nil { return m.checkFunc(ctx, id) }
	return true, nil
}
func (m *mockRoomClient) Close() error { return nil }

func (m *mockScheduleRepo) CreateSchedule(ctx context.Context, roomID string, days []int64, start, end *time.Time) (*domain.Schedule, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, roomID, days, start, end)
	}
	return &domain.Schedule{RoomID: roomID}, nil
}
