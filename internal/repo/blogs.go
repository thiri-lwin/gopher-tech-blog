package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

type Blog struct {
	ID       int
	Title    string
	Summary  string
	Content  string
	Author   string
	PostedAt time.Time
	Date     string
}

func (r *repo) GetBlogs(ctx context.Context, limit, offset int) ([]Blog, error) {
	query := `SELECT id, title, summary, content, author, posted_at FROM posts ORDER BY posted_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query blogs: %w", err)
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		var blog Blog
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Summary, &blog.Content, &blog.Author, &blog.PostedAt); err != nil {
			return nil, fmt.Errorf("failed to scan blog: %w", err)
		}
		blog.Date = blog.PostedAt.Format("January 2, 2006")
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (r *repo) GetBlogsCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM posts`
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count blogs: %w", err)
	}
	return count, nil
}

func (r *repo) GetBlog(ctx context.Context, id string) (Blog, error) {
	query := `SELECT id, title, summary, content, author, posted_at FROM posts WHERE id = $1`
	var blog Blog
	err := r.db.QueryRow(ctx, query, id).Scan(&blog.ID, &blog.Title, &blog.Summary, &blog.Content, &blog.Author, &blog.PostedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Blog{}, fmt.Errorf("blog not found")
		}
		return Blog{}, fmt.Errorf("failed to query blog: %w", err)
	}
	blog.Date = blog.PostedAt.Format("January 2, 2006")
	return blog, nil
}
