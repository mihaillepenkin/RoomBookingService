package domain

import (
	"fmt"
	"time"
)

type Schedule struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	DaysOfWeek []int64 `json:"days_of_week"`
	StartTime *time.Time `json:"start_time"`
	EndTime *time.Time `json:"end_time"`
}

func (s *Schedule) Validate() error {
	for _, DayOfWeek := range s.DaysOfWeek {
		if !(DayOfWeek <= 7 && DayOfWeek >= 1) {
			return fmt.Errorf("value of day need be between 1 and 7")
		}
	}
	if (s.StartTime.Compare(*s.EndTime) != -1) {
		return fmt.Errorf("start time and end time need to be different and start < end")
	}
	return nil
}