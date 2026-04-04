package domain

type Booking struct {
	ID        string `json:"id"`
	SlotID    string `json:"slot_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
}