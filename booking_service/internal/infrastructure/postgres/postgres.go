package postgres

import (
	"booking_service/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type ScheduleRepository struct {
	DB *sql.DB
}

func (s *ScheduleRepository) CreateSchedule(ctx context.Context, RoomID string, DayOfWeek []int64, startTime *time.Time, endTime *time.Time) (*domain.Schedule, error) {
	var sch domain.Schedule
	row := s.DB.QueryRowContext(ctx, `INSERT INTO schedules (room_id, days_of_week, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING id`, RoomID, pq.Array(DayOfWeek), startTime.UTC(), endTime.UTC())
	err := row.Scan(&sch.ID)
	if (err != nil) {
		return nil, err
	}
	sTime, _ := time.Parse("15:04", startTime.String())
	sch.StartTime = &sTime
	eTime, _ := time.Parse("15:04", endTime.String())
	sch.EndTime = &eTime
	sch.RoomID = RoomID
	sch.DaysOfWeek = DayOfWeek
	return &sch, nil
}

func (s *ScheduleRepository) GetScheduleByRoomID(ctx context.Context, RoomID string) (*domain.Schedule, error) {
	var sch domain.Schedule
	row := s.DB.QueryRowContext(ctx, `SELECT id, room_id, days_of_week, start_time, end_time FROM schedules WHERE room_id = $1`, RoomID)
	var startStr, endStr string
	sch.DaysOfWeek = make([]int64, 0)
	err := row.Scan(&sch.ID, &sch.RoomID, pq.Array(&sch.DaysOfWeek), &startStr, &endStr)
	if (err != nil) {
		fmt.Println(err)
		return nil, err
	}
	sTime, _ := time.Parse("15:04", startStr)
	sch.StartTime = &sTime
	eTime, _ := time.Parse("15:04", endStr)
	sch.EndTime = &eTime
	return &sch, nil
}



type SlotRepository struct {
	DB *sql.DB
}

func (s *SlotRepository) GetAllSlotsByRoomIDAndDate(ctx context.Context, RoomID string, Date string) (*[]domain.Slot, error) {
	listOfBooking := make([]domain.Slot, 0)
	rows, err := s.DB.QueryContext(ctx, `SELECT start_time FROM bookings WHERE room_id = $1 AND date = $2 AND status = 'active'`, RoomID, Date)
	if (err != nil) {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
        var startTime time.Time
        err = rows.Scan(&startTime)
        if err != nil {
            return nil, err
        }
        slot := domain.Slot{
            StartTime: &startTime,
        }
        listOfBooking = append(listOfBooking, slot)
    }
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &listOfBooking, nil
}

type BookingRepository struct {
	DB *sql.DB
}

func (b *BookingRepository) CreateBook(ctx context.Context, slot domain.Slot, userId string, slotId string) (*domain.Booking, error) {
	row := b.DB.QueryRowContext(ctx, `INSERT INTO bookings (room_id, date, start_time, end_time, status, user_id) VALUES ($1, $2, $3, $4, 'active', $5) RETURNING id`, slot.RoomID, slot.Date, slot.StartTime, slot.EndTime, userId)
	book := domain.Booking{Status: "active", SlotID: slotId, UserID: userId}
	err := row.Scan(&book.ID)
	if (err != nil) {
		return nil, err
	}
	return &book, err
}

func (b *BookingRepository) GetAllBooksByUserID(ctx context.Context, userId string) (*[]domain.Booking, *[]time.Time, error) {
	rows, err := b.DB.QueryContext(ctx, `SELECT id, room_id, date, start_time FROM bookings WHERE status = 'active' AND user_id = $1`, userId)
	if (err != nil) {
		return nil, nil, err
	}
	defer rows.Close()
	bookings := make([]domain.Booking, 0)
	times := make([]time.Time, 0)
	for rows.Next() {
		var id string
		var roomId string
		var data time.Time
		var startTime time.Time
		err = rows.Scan(&id, &roomId, &data, &startTime)
		if (err != nil) {
			return nil, nil, err
		}
		slotId := roomId + "|" + data.Format("2006-01-02") + "|" + startTime.Format("15:04")
		book := domain.Booking{SlotID: slotId, UserID: userId, Status: "active", ID: id}
		bookings = append(bookings, book)
		times = append(times, startTime)
	}
	return &bookings, &times, nil
}

func (b *BookingRepository) UpdateBookStatusByID(ctx context.Context, id string, userID string) (*domain.Booking, error) {
    row := b.DB.QueryRowContext(ctx, `SELECT status, user_id, room_id, date, start_time FROM bookings WHERE id = $1 AND user_id = $2`, id, userID)
    var status string
    book := domain.Booking{ID: id}
    var rId string
    var date time.Time
    var startTime time.Time
    err := row.Scan(&status, &book.UserID, &rId, &date, &startTime)
    if err != nil {
        return nil, err
    }
    _, err = b.DB.ExecContext(ctx, `UPDATE bookings SET status = 'cancelled', updated_at = NOW() WHERE id = $1 AND user_id = $2`, id, userID)
    if err != nil {
        return nil, err
    }
    book.Status = "cancelled"
    book.SlotID = rId + "|" + date.Format("2006-01-02") + "|" + startTime.Format("15:04")
    return &book, nil
}

func (b *BookingRepository) GetAllBookings(ctx context.Context, page, pageSize int) (*[]domain.Booking, int, error) {
    var totalCount int
    err := b.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM bookings`).Scan(&totalCount)
    if err != nil {
        return nil, 0, err
    }
    offset := (page - 1) * pageSize
    rows, err := b.DB.QueryContext(ctx, `SELECT id, room_id, date, start_time, status, user_id FROM bookings ORDER BY created_at DESC LIMIT $1 OFFSET $2`, pageSize, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    bookings := make([]domain.Booking, 0)
    for rows.Next() {
        var book domain.Booking
		var rId string
		var date time.Time
		var startTime time.Time
        err := rows.Scan(&book.ID, &rId, &date, &startTime, &book.Status, &book.UserID)
        if err != nil {
            return nil, 0, err
        }
		book.SlotID = rId + "|" + date.Format("2006-01-02") + "|" + startTime.Format("15:04")
        bookings = append(bookings, book)
    }
    return &bookings, totalCount, nil
}