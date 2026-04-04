package domain

import (
	"fmt"
	"time"
)

type Slot struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	Date string `json:"date"`
	StartTime *time.Time `json:"start"`
	EndTime *time.Time `json:"end"`
}

func (s *Slot)GenerateID() {
    s.ID = fmt.Sprintf("%s|%s|%s", s.RoomID, s.Date, s.StartTime.Format("15:04"))
}
