package interfaces

type CreateScheduleDTO struct {
	RoomID    string     `json:"room_id,omitempty"`
	DaysOfWeek []int64        `json:"days_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type OutputDTO struct {
	Status int `json:"status"`
	Data map[string]interface{} `json:"data"`
}

