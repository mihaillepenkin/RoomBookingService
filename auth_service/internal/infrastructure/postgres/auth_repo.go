package postgres

import (
	"auth_service/internal/domain"
	"context"
	"database/sql"
)

type UserRepoPG struct {
	DB *sql.DB
}

func (r *UserRepoPG) CreateUser(ctx context.Context, email, password, role string) (*domain.User, error) {
	row := r.DB.QueryRowContext(ctx, `INSERT INTO users (password, email, role)
            VALUES ($1, $2, $3) RETURNING id`, password, email, role)
	var user domain.User
	if err := row.Scan(&user.ID); err != nil {
		return nil, err
	}
	user.Email = email
	user.Password = password
	user.Role = role
	return &user, nil		
}

func (r *UserRepoPG) GetUser(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	row := r.DB.QueryRowContext(ctx, `SELECT id, password, email, role FROM users WHERE email = $1`, email)
	if err := row.Scan(&user.ID, &user.Password, &user.Email, &user.Role); err != nil {
		return nil, err
	}
	return user, nil
}