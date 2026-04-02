package postgres

import (
	"context"
	"database/sql"
	"room_service/internal/domain"
)

type RoomRepository struct {
	DB *sql.DB
}

func (r *RoomRepository) GetAllRooms(ctx context.Context) (*[]domain.Room, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, description, capacity FROM rooms ORDER BY name ASC`)
	if (err != nil) {
		return nil, err
	}
	defer rows.Close()
	var rooms []domain.Room
	for rows.Next() {
		var room domain.Room
		err = rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity)
		if (err != nil) {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &rooms, nil
}

func (r *RoomRepository) CreateRoom(ctx context.Context, name, description string, capacity int) (*domain.Room, error) {
	rows := r.DB.QueryRowContext(ctx, `INSERT INTO rooms (name, description, capacity) VALUES ($1, $2, $3) RETURNING id`, name, description, capacity)
	var room domain.Room
	err := rows.Scan(&room.ID)
	if (err != nil) {
		return nil, err
	}
	room.Capacity = capacity
	room.Description = description
	room.Name = name
	return &room, nil
}

func (r *RoomRepository) GetRoom(ctx context.Context, id string) (*domain.Room, error) {
	rows := r.DB.QueryRowContext(ctx, `SELECT id, name, description, capacity FROM rooms WHERE id = $1`, id)
	var room domain.Room
	err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity)
	if (err != nil) {
		return nil, err
	}
	return &room, nil
}