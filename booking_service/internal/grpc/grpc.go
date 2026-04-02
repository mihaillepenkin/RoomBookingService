package grpc

import (
    "context"
    "fmt"

    roomv1 "booking_service/proto/room/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type RoomServiceClient interface {
    CheckRoomExists(ctx context.Context, roomID string) (bool, error)
    Close() error
}

type roomGRPCClient struct {
    client roomv1.RoomServiceClient
    conn   *grpc.ClientConn
}

func NewRoomServiceClient(addr string) (RoomServiceClient, error) {
    conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, fmt.Errorf("failed to dial room service: %w", err)
    }
    return &roomGRPCClient{
        client: roomv1.NewRoomServiceClient(conn),
        conn:   conn,
    }, nil
}

func (c *roomGRPCClient) CheckRoomExists(ctx context.Context, roomID string) (bool, error) {
    resp, err := c.client.CheckRoomExists(ctx, &roomv1.CheckRoomExistsRequest{
        RoomId: roomID,
    })
    if err != nil {
        return false, fmt.Errorf("grpc call failed: %w", err)
    }
    return resp.Exists, nil
}

func (c *roomGRPCClient) Close() error {
    if c.conn != nil {
        return c.conn.Close()
    }
    return nil
}