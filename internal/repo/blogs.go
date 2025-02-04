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
	Likes    *int
	Comments []Comment
}

type Comment struct {
	ID      int
	PostID  int
	UserID  int
	Content string
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

func (r *repo) GetBlog(ctx context.Context, id int) (Blog, error) {
	query := `SELECT id, title, summary, content, author, likes, posted_at FROM posts WHERE id = $1`
	var blog Blog
	err := r.db.QueryRow(ctx, query, id).Scan(&blog.ID, &blog.Title, &blog.Summary, &blog.Content, &blog.Author, &blog.Likes, &blog.PostedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Blog{}, fmt.Errorf("blog not found")
		}
		return Blog{}, fmt.Errorf("failed to query blog: %w", err)
	}
	blog.Date = blog.PostedAt.Format("January 2, 2006")
	if blog.Likes == nil {
		zero := 0
		blog.Likes = &zero
	}

	comments, err := r.GetComments(ctx, id)
	if err != nil {
		return Blog{}, fmt.Errorf("failed to get comments: %w", err)
	}
	blog.Comments = comments
	return blog, nil
}

func (r *repo) GetComments(ctx context.Context, postID int) ([]Comment, error) {
	query := `SELECT id, post_id, content FROM comments WHERE post_id = $1`
	rows, err := r.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *repo) LikeBlog(ctx context.Context, id int) (int, error) {
	query := `UPDATE posts SET likes = likes + 1 WHERE id = $1 RETURNING likes`
	var likes int
	err := r.db.QueryRow(ctx, query, id).Scan(&likes)
	if err != nil {
		return 0, fmt.Errorf("failed to like blog: %w", err)
	}
	return likes, nil
}

func (r *repo) CommentBlog(ctx context.Context, comment Comment) error {
	query := `INSERT INTO comments (post_id, content) VALUES ($1, $2) RETURNING id, post_id`
	err := r.db.QueryRow(ctx, query, comment.PostID, comment.Content).Scan(&comment.ID, &comment.PostID)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}
	return nil
}
