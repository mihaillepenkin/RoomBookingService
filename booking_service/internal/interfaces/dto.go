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


type GetAllSlotsDTO struct {
	RoomID    string     `json:"room_id,omitempty"`
	Date      string     `json:"date,omitempty"`
}

type CreateBookingDTO struct {
	SlotID    string     `json:"slot_id"`
}

type CancelBookingDTO struct {
	BookingID    string     `json:"booking_id,omitempty"`
}

type PaginationDTO struct {
    Page       int `json:"page"`
    PageSize   int `json:"pageSize"`
    TotalCount int `json:"totalCount"`
    TotalPages int `json:"totalPages"`
}

type GetAllBooksDTO struct {
    Page       int `json:"page"`
    PageSize   int `json:"pageSize"`
}