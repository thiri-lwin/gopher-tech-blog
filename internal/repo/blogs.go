package repo

import (
	"context"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Blog struct {
	ID      string `bson:"id" json:"id"`
	Title   string `bson:"title" json:"title"`
	Summary string `bson:"summary" json:"summary"`
	Content string `bson:"content" json:"content"`
	Author  string `bson:"author" json:"author"`
}

func (r *repo) GetBlogs(ctx context.Context, limit, page int) ([]Blog, error) {
	var blogs []Blog

	collection := r.db.Collection("posts")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	findOptions := newMongoPaginate(limit, page).getPaginatedOpts()

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var blog Blog
		if err := cursor.Decode(&blog); err != nil {
			log.Println("Error decoding post:", err)
			continue
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (r *repo) GetBlog(ctx context.Context, id string) (Blog, error) {
	collection := r.db.Collection("posts")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var blog Blog
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&blog)
	if err != nil {
		return Blog{}, err
	}
	return blog, nil
}
