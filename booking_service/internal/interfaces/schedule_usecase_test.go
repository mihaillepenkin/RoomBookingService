package interfaces_test

import (
	"context"
	"testing"
	"time"

	"booking_service/internal/domain"
	"booking_service/internal/interfaces"
)

func TestScheduleUsecase_CreateSchedule_Admin_Success(t *testing.T){
	repo:=&mockScheduleRepo{
		getFunc:func(ctx context.Context,roomID string)(*domain.Schedule,error){return nil,nil},
		createFunc:func(ctx context.Context,roomID string,days[]int64,start,end*time.Time)(*domain.Schedule,error){
			return &domain.Schedule{RoomID:roomID},nil
		},
	}
	roomClient:=&mockRoomClient{checkFunc:func(ctx context.Context,id string)(bool,error){return true,nil}}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{Role:"admin"},nil}}
	uc:=interfaces.NewScheduleService(repo,roomClient,jwt)
	out:=uc.CreateSchedule(context.Background(),interfaces.CreateScheduleDTO{RoomID:"r1",DaysOfWeek:[]int64{1,2},StartTime:"09:00",EndTime:"18:00"},"token")
	if out.Status!=201{t.Errorf("expected 201,got %d",out.Status)}
}

func TestScheduleUsecase_CreateSchedule_RoomNotFound(t *testing.T){
	repo:=&mockScheduleRepo{}
	roomClient:=&mockRoomClient{checkFunc:func(ctx context.Context,id string)(bool,error){return false,nil}}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{Role:"admin"},nil}}
	uc:=interfaces.NewScheduleService(repo, roomClient, jwt)
	out:=uc.CreateSchedule(context.Background(),interfaces.CreateScheduleDTO{RoomID:"notfound",StartTime: "09:00",EndTime: "18:00"}, "token")
	if out.Status!=404{t.Errorf("expected 404,got %d",out.Status)}
}

func TestScheduleUsecase_CreateSchedule_AlreadyExists(t *testing.T){
	repo:=&mockScheduleRepo{
		getFunc:func(ctx context.Context,roomID string)(*domain.Schedule,error){
			return &domain.Schedule{RoomID:roomID},nil
		},
	}
	roomClient:=&mockRoomClient{checkFunc:func(ctx context.Context,id string)(bool,error){return true,nil}}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{Role:"admin"},nil}}
	uc:=interfaces.NewScheduleService(repo,roomClient,jwt)
	out:=uc.CreateSchedule(context.Background(),interfaces.CreateScheduleDTO{RoomID:"r1",StartTime:"09:00",EndTime:"18:00"},"token")
	if out.Status!=409{t.Errorf("expected 409,got %d",out.Status)}
}