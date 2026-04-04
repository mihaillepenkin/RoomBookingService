package interfaces

import (
	"booking_service/internal/domain"
	"booking_service/internal/grpc"
	"context"
	"fmt"
	"time"
)

type ScheduleUsecase struct {
	repo       ScheduleRepository
	jwt        JWTService
	roomClient grpc.RoomServiceClient
}

func NewScheduleService(repo ScheduleRepository, roomClient grpc.RoomServiceClient, jwt JWTService) *ScheduleUsecase {
    return &ScheduleUsecase{
        repo:       repo,
        roomClient: roomClient,
		jwt:        jwt,
    }
}

func (s *ScheduleUsecase) CreateSchedule(ctx context.Context, input CreateScheduleDTO, token string) OutputDTO {
	user, err := s.jwt.ValidateToken(token)
	if (err != nil) {
		return OutputDTO{Status: 401, Data: map[string]interface{}{"error": "you are not authtorized"}}
	}
	if (user.Role != "admin") {
		return OutputDTO{Status: 403, Data: map[string]interface{}{"error": "for this action need admin role"}}
	}
	startTime, err := time.Parse("15:04", input.StartTime)
    if err != nil {
        return OutputDTO{Status: 400, Data: map[string]interface{}{"error": fmt.Sprintf("invalid start_time format: %w", err)}}
    }
    
    endTime, err := time.Parse("15:04", input.EndTime)
    if err != nil {
        return OutputDTO{Status: 400, Data: map[string]interface{}{"error": fmt.Sprintf("invalid end_time format: %w", err)}}
    }
	sch := &domain.Schedule{
		DaysOfWeek: input.DaysOfWeek,
		RoomID:     input.RoomID,
		StartTime:  &startTime,
		EndTime:    &endTime,
	}
	err = sch.Validate()
	if (err != nil) {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": err.Error()}}
	}
	exist, err := s.roomClient.CheckRoomExists(ctx, input.RoomID)
	if (err != nil) {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	if (!exist) {
		return OutputDTO{Status: 404, Data: map[string]interface{}{"error": "room not found"}}
	}
	sch, err = s.repo.GetScheduleByRoomID(ctx, input.RoomID)
	if (sch != nil) {
		return OutputDTO{Status: 409, Data: map[string]interface{}{"error": "schedule for this room already exists and cannot be changed"}}
	}
	fmt.Println(sch)
	fmt.Println(err)
	sch, err = s.repo.CreateSchedule(ctx, input.RoomID, input.DaysOfWeek, &startTime, &endTime)
	if (err != nil) {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 201, Data: map[string]interface{}{"schedule": *sch}}
}