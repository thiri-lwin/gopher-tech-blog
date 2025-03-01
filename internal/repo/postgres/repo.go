package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repo interface {
	GetBlogs(ctx context.Context, limit, page int) ([]Blog, error)
	GetBlog(ctx context.Context, id int) (Blog, error)
	GetBlogsCount(ctx context.Context) (int64, error)
	LikeBlog(ctx context.Context, userID, id int) (int, error)
	CommentBlog(ctx context.Context, comment Comment) error

	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, email string) (User, error)
}

type repo struct {
	db *pgxpool.Pool
}

func New(dbURI string) Repo {
	return &repo{
		db: connectDB(dbURI),
	}
}

func connectDB(dbURI string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		log.Fatalf("Failed to parse PostgreSQL config: %v", err)
	}

	// Create a new PostgreSQL connection pool
	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Set a timeout for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ensure the connection is established
	err = db.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Return the database pool instance
	return db
}
