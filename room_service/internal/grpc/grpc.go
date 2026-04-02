package grpc

import (
	"context"
	"room_service/internal/interfaces"
	roomv1 "room_service/proto/room/v1"
)

type RoomServer struct {
    roomv1.UnimplementedRoomServiceServer
    repo interfaces.RoomRepository
}

func NewRoomServer(repo interfaces.RoomRepository) *RoomServer {
    return &RoomServer{repo: repo}
}


func (s *RoomServer) CheckRoomExists(ctx context.Context, req *roomv1.CheckRoomExistsRequest) (*roomv1.CheckRoomExistsResponse, error) {
    room, _ := s.repo.GetRoom(ctx, req.RoomId)
    return &roomv1.CheckRoomExistsResponse{Exists: room != nil}, nil
}