package interfaces_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"booking_service/internal/domain"
	"booking_service/internal/interfaces"
)


func TestBookingUsecase_CancelBook_NotFound(t *testing.T){
	repo:=&mockBookingRepo{
		updateFunc:func(ctx context.Context,id,userID string)(*domain.Booking,error){
			return nil,sql.ErrNoRows
		},
	}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{Role:"user"},nil}}
	uc:=interfaces.NewBookingService(repo,jwt,nil)
	out:=uc.CancelBook(context.Background(),interfaces.CancelBookingDTO{BookingID:"notfound"},"token")
	if out.Status!=404{t.Errorf("expected 404,got %d",out.Status)}
}

func TestBookingUsecase_GetAllBookings_Admin_Pagination(t *testing.T){
	repo:=&mockBookingRepo{
		getAllFunc:func(ctx context.Context,page,pageSize int)(*[]domain.Booking,int,error){
			b:=[]domain.Booking{{ID:"b1"}}
			return &b,100,nil
		},
	}
	jwt:=&mockJWT{validateFunc:func(t string)(*domain.User,error){return &domain.User{Role:"admin"},nil}}
	uc:=interfaces.NewBookingService(repo,jwt,nil)
	out:=uc.GetAllBookings(context.Background(),interfaces.GetAllBooksDTO{Page:1,PageSize:20},"token")
	if out.Status!=200{t.Errorf("expected 200,got %d",out.Status)}
	pag:=out.Data["pagination"].(interfaces.PaginationDTO)
	fmt.Println(pag)
	if pag.TotalPages!=5{t.Errorf("expected 5 pages,got %d",pag.TotalPages)}
}