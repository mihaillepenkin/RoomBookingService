package interfaces

import (
	"booking_service/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type BookingUsecase struct {
	repo       BookingRepository
	jwt        JWTService
	sRepo      ScheduleRepository
}

func NewBookingService(repo BookingRepository, jwt JWTService, sRepo ScheduleRepository) *BookingUsecase {
    return &BookingUsecase{
        repo:       repo,
		jwt:        jwt,
		sRepo:      sRepo,
    }
}

func parsedSlotId(id string) (string, string, string, error) {
	parsed := strings.Split(id, "|")
	if (len(parsed) != 3) {
		return "", "", "", fmt.Errorf("invalid slotId")
	}
	return parsed[0], parsed[1], parsed[2], nil
}

func (b *BookingUsecase) CreateBook(ctx context.Context, input CreateBookingDTO, token string) OutputDTO {
	user, err := b.jwt.ValidateToken(token)
	if (err != nil) {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "you are not authtorized"}}
	}
	if (user.Role != "user") {
		return OutputDTO{Status: 403, Data: map[string]interface{}{"error": "for this action need user role"}}
	}
	roomId, date, startTimeString, err := parsedSlotId(input.SlotID)
	if (err != nil) {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": err.Error()}}
	}
	startTime, err := time.Parse("15:04", startTimeString)
	if (err != nil) {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": "invalid slotId"}}
	}
	schedule, _ := b.sRepo.GetScheduleByRoomID(ctx, roomId)
	if (schedule == nil) {
		return OutputDTO{Status: 404, Data: map[string]interface{}{"error": "slot is not founded"}}
	}
	today, err := time.Parse("2006-01-02", date)
	now := time.Now().UTC()
	slotDateTime := time.Date(today.Year(), today.Month(), today.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	if slotDateTime.Before(now) {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": "cannot book past slots"}}
	}
	flag := true
	for _, el := range schedule.DaysOfWeek {
		if ((el % 7) == int64(int(today.Weekday()))) {
			flag = false
			break
		}
	}
	if (flag) {
		return OutputDTO{Status: 404, Data: map[string]interface{}{"error": "slot is not founded"}}
	}
	flag = true
	tt, _ := time.Parse("15:04", schedule.StartTime.Format("15:04"))
	for i := tt; i.Add(30*time.Minute).Compare(*schedule.EndTime) <= 0; i = i.Add(30*time.Minute) {
		if (startTime.Compare(i) == 0) {
			flag = false
			break
		}
		if (startTime.Compare(i) == -1) {
			break
		}
	}
	if (flag) {
		return OutputDTO{Status: 404, Data: map[string]interface{}{"error": "slot is not founded"}}
	}
	endTime := startTime.Add(30 * time.Minute)
	slot := domain.Slot{RoomID: roomId, StartTime: &startTime, EndTime: &endTime, Date: date}
	slot.GenerateID()
	book, err := b.repo.CreateBook(ctx, slot, user.ID, input.SlotID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return OutputDTO{Status: 409, Data: map[string]interface{}{"error": "slot already booked"}}
		}
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 201, Data: map[string]interface{}{"booking": *book}}
}

func (b *BookingUsecase) GetMyBooks(ctx context.Context, token string) OutputDTO {
	user, err := b.jwt.ValidateToken(token)
	if (err != nil) {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "you are not authtorized"}}
	}
	if (user.Role != "user") {
		return OutputDTO{Status: 403, Data: map[string]interface{}{"error": "for this action need user role"}}
	}
	listOfAllBooks, listOfStarts, err := b.repo.GetAllBooksByUserID(ctx, user.ID)
	if (err != nil) {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	now := time.Now()
	listOfActualBooks := make([]domain.Booking, 0)
	for ind, el := range *listOfStarts {
		if (el.Compare(now) >= -1) {
			listOfActualBooks = append(listOfActualBooks, (*listOfAllBooks)[ind])
		}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"list": listOfActualBooks}}
}

func (b *BookingUsecase) CancelBook(ctx context.Context, input CancelBookingDTO, token string) OutputDTO {
	user, err := b.jwt.ValidateToken(token)
	if (err != nil) {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "you are not authtorized"}}
	}
	if (user.Role != "user") {
		return OutputDTO{Status: 403, Data: map[string]interface{}{"error": "for this action need user role"}}
	}
	book, err := b.repo.UpdateBookStatusByID(ctx, input.BookingID, user.ID)
	if book == nil {
        book = &domain.Booking{}
    }
	if err != nil {
		if err == sql.ErrNoRows {
			return OutputDTO{Status: 404, Data: map[string]interface{}{"error": err.Error()}}
		}
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"booking": *book}} 
}

func (b *BookingUsecase) GetAllBookings(ctx context.Context, input GetAllBooksDTO, token string) OutputDTO {
    user, err := b.jwt.ValidateToken(token)
    if err != nil {
        return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "you are not authorized"}}
    }
    if user.Role != "admin" {
        return OutputDTO{Status: 403, Data: map[string]interface{}{"error": "for this action need admin role"}}
    }
    if input.Page < 1 {
        return OutputDTO{Status: 400, Data: map[string]interface{}{"error": "invalid page value"}}
    }
    if input.PageSize < 1 {
        return OutputDTO{Status: 400, Data: map[string]interface{}{"error": "invalid page size value"}}
    }
    if input.PageSize > 100 {
        return OutputDTO{Status: 400, Data: map[string]interface{}{"error": "invalid page size value"}}
    }
    bookings, totalCount, err := b.repo.GetAllBookings(ctx, input.Page, input.PageSize)
    if err != nil {
        return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
    }
    totalPages := (totalCount + input.PageSize - 1) / input.PageSize
	if bookings == nil {
        bookings = &[]domain.Booking{}
    }
    return OutputDTO{
        Status: 200,
        Data: map[string]interface{}{
            "bookings": *bookings,
            "pagination": PaginationDTO{
                Page:       input.Page,
                PageSize:   input.PageSize,
                TotalCount: totalCount,
                TotalPages: totalPages,
            },
        },
    }
}