package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo interface {
	GetBlogs(ctx context.Context, limit, page int) ([]Blog, error)
	GetBlog(ctx context.Context, id string) (Blog, error)
	GetBlogsCount(ctx context.Context) (int64, error)
}

type repo struct {
	db *mongo.Database
}

func New(dbURI string) Repo {
	return &repo{
		db: connectDB(dbURI),
	}

}

func connectDB(dbURI string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(dbURI)
	//clientOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: false})

	// Create a new MongoDB client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Set a timeout for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ensure the client is connected
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Connect to the desired database
	db := client.Database("tech_blog_db")
	fmt.Println("Successfully connected to MongoDB!")

	// Return the database instance
	return db
}

type mongoPaginate struct {
	limit int64
	page  int64
}

func newMongoPaginate(limit, page int) *mongoPaginate {
	return &mongoPaginate{
		limit: int64(limit),
		page:  int64(page),
	}
}

func (m *mongoPaginate) getPaginatedOpts() *options.FindOptions {
	opts := options.Find()
	opts.SetLimit(m.limit)
	opts.SetSkip(m.page)
	return opts
}
