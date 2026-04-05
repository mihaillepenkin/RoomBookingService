package interfaces_test

import (
	"context"
	"testing"
	"time"

	"booking_service/internal/domain"
	"booking_service/internal/interfaces"
)

func TestSlotUsecase_GetAllSlots_Success(t *testing.T){
	date:="2024-06-10"
	today,_:=time.Parse("2006-01-02",date)
	repo:=&mockBookingRepo{
		getBookedFunc:func(ctx context.Context,roomID,date string)(*[]domain.Slot,error){
			s:=[]domain.Slot{}
			return &s,nil
		},
	}
	sRepo:=&mockScheduleRepo{
		getFunc:func(ctx context.Context,roomID string)(*domain.Schedule,error){
			st,_:=time.Parse("15:04","09:00")
			et,_:=time.Parse("15:04","18:00")
			return &domain.Schedule{RoomID:roomID,DaysOfWeek:[]int64{int64(today.Weekday())},StartTime:&st,EndTime:&et},nil
		},
	}
	roomClient:=&mockRoomClient{checkFunc:func(ctx context.Context,id string)(bool,error){return true,nil}}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{},nil}}
	uc:=interfaces.NewSlotService(repo,jwt,sRepo,roomClient)
	out:=uc.GetAllSlots(context.Background(),interfaces.GetAllSlotsDTO{RoomID:"r1",Date:date},"token")
	if out.Status!=200{t.Errorf("expected 200,got %d",out.Status)}
	slots:=out.Data["list"].([]domain.Slot)
	if len(slots)==0{t.Errorf("expected slots,got empty")}
}

func TestSlotUsecase_GetAllSlots_NoSchedule(t *testing.T){
	repo:=&mockBookingRepo{}
	sRepo:=&mockScheduleRepo{getFunc:func(ctx context.Context,roomID string)(*domain.Schedule,error){return nil,nil}}
	roomClient:=&mockRoomClient{checkFunc:func(ctx context.Context,id string)(bool,error){return true,nil}}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{},nil}}
	uc:=interfaces.NewSlotService(repo,jwt,sRepo,roomClient)
	out:=uc.GetAllSlots(context.Background(),interfaces.GetAllSlotsDTO{RoomID:"r1",Date:"2024-06-10"},"token")
	if out.Status!=200{t.Errorf("expected 200,got %d",out.Status)}
	if len(out.Data["list"].([]domain.Slot))!=0{t.Errorf("expected empty list")}
}

func TestSlotUsecase_GetAllSlots_WrongDay(t *testing.T){
	date:="2024-06-10"
	repo:=&mockBookingRepo{}
	sRepo:=&mockScheduleRepo{
		getFunc:func(ctx context.Context,roomID string)(*domain.Schedule,error){
			st,_:=time.Parse("15:04","09:00")
			et,_:=time.Parse("15:04","18:00")
			return &domain.Schedule{RoomID:roomID,DaysOfWeek:[]int64{7},StartTime:&st,EndTime:&et},nil
		},
	}
	roomClient:=&mockRoomClient{checkFunc:func(ctx context.Context,id string)(bool,error){return true,nil}}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{},nil}}
	uc:=interfaces.NewSlotService(repo,jwt,sRepo,roomClient)
	out:=uc.GetAllSlots(context.Background(),interfaces.GetAllSlotsDTO{RoomID:"r1",Date:date},"token")
	if out.Status!=200{t.Errorf("expected 200,got %d",out.Status)}
	if len(out.Data["list"].([]domain.Slot))!=0{t.Errorf("expected empty list for wrong day")}
}