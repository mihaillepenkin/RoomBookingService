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

