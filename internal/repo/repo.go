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
	GetBlogs(ctx context.Context) ([]Blog, error)
	GetBlog(ctx context.Context, id string) (Blog, error)
}

type repo struct {
	db *mongo.Database
}

func New() Repo {
	return &repo{
		db: connectDB(),
	}

}

func connectDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("tech_blog_db")
	fmt.Println("Connected to MongoDB!")
	return db
}
