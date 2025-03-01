package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *repo) CreateUser(ctx context.Context, user User) (int, error) {
	query := `INSERT INTO users (first_name, last_name, email, password)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id`
	err := r.db.QueryRow(ctx, query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	return user.ID, nil
}

func (r *repo) GetUser(ctx context.Context, email string) (User, error) {
	query := `SELECT id, first_name, last_name, email, password, created_at FROM users WHERE email = $1`
	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, err
	}
	return user, nil
}

func (r *repo) GetUserByID(ctx context.Context, id int) (User, error) {
	query := `SELECT id, first_name, last_name, email, password, created_at FROM users WHERE id = $1`
	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, err
	}
	return user, nil
}
