package repo

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
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

func New(dbUser, dbPassword string) Repo {
	return &repo{
		db: connectDB(dbUser, dbPassword),
	}

}

func connectDB(dbUser, dbPassword string) *mongo.Database {
	// URL-encode password in case it contains special characters
	encodedPassword := url.QueryEscape(dbPassword)

	// Create MongoDB URI with encoded password
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@mongo-cluster01.fgq10.mongodb.net/tech_blog_db?retryWrites=true&w=majority", dbUser, encodedPassword))
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: false})

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
