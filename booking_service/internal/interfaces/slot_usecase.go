package interfaces

import (
	"booking_service/internal/domain"
	"booking_service/internal/grpc"
	"context"
	"time"
)

type SlotUsecase struct {
	repo       SlotRepository
	jwt        JWTService
	sRepo      ScheduleRepository
	roomClient grpc.RoomServiceClient
}

func NewSlotService(repo SlotRepository, jwt JWTService, sRepo ScheduleRepository, roomClient grpc.RoomServiceClient) *SlotUsecase {
    return &SlotUsecase{
        repo:       repo,
		jwt:        jwt,
		sRepo:      sRepo,
		roomClient: roomClient,
    }
}

func (s *SlotUsecase) GetAllSlots(ctx context.Context, input GetAllSlotsDTO, token string) OutputDTO {
	today, err := time.Parse("2006-01-02", input.Date)
	if (err != nil) {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": "invalid format of date"}} 
	}
	_, err = s.jwt.ValidateToken(token)
	if (err != nil) {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "you are not authtorized"}}
	}
	exist, err := s.roomClient.CheckRoomExists(ctx, input.RoomID)
	if (err != nil) {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	if (!exist) {
		return OutputDTO{Status: 404, Data: map[string]interface{}{"error": "room not found"}}
	}
	schedule, _ := s.sRepo.GetScheduleByRoomID(ctx, input.RoomID)
	if (schedule == nil) {
		return OutputDTO{Status: 200, Data: map[string]interface{}{"list": []domain.Slot{}}}
	}
	flag := true
	for _, el := range schedule.DaysOfWeek {
		if ((el % 7) == int64(int(today.Weekday()))) {
			flag = false
			break
		}
	}
	if (flag) {
		return OutputDTO{Status: 200, Data: map[string]interface{}{"list": []domain.Slot{}}}
	}
	checkMap := make(map[string]bool)
	bookSlot, err := s.repo.GetAllSlotsByRoomIDAndDate(ctx, input.RoomID, input.Date)
	if (err != nil) {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}} 
	}
	for _, el := range *bookSlot {
		checkMap[el.StartTime.Format("15:04")] = true
	}
	listOfSlots := make([]domain.Slot, 0)
	start := schedule.StartTime
	end := schedule.EndTime
	for i := *start; i.Add(30*time.Minute).Compare(*end) <= 0; i = i.Add(30*time.Minute) {
		if (checkMap[i.Format("15:04")]) {
			continue
		}
		startTime := i
		endTime := i.Add(30 * time.Minute)
		slot := domain.Slot{
			Date:      today.Format("2006-01-02"),
			RoomID:    input.RoomID,
			StartTime: &startTime,
			EndTime:   &endTime,
		}
		slot.GenerateID()
		listOfSlots = append(listOfSlots, slot)
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"list": listOfSlots}}
}