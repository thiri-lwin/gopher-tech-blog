package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
)

type Blog struct {
	ID          int
	Title       string
	Summary     string
	Content     string
	Author      string
	PostedAt    time.Time
	LikedByUser bool // optional field
	Date        string
	Likes       int
	Comments    []Comment
}

type Comment struct {
	ID      int
	PostID  int
	UserID  int
	Content string

	UserName string
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

	// Get comments
	comments, err := r.GetComments(ctx, id)
	if err != nil {
		return Blog{}, fmt.Errorf("failed to get comments: %w", err)
	}
	blog.Comments = comments

	// Get likes
	count, _ := r.GetLikes(ctx, id)
	blog.Likes = count
	return blog, nil
}

func (r *repo) GetLikes(ctx context.Context, postID int) (int, error) {
	query := `SELECT count(*) FROM likes WHERE post_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, postID).Scan(&count)
	if err != nil {
		log.Println("get post like count failed: ", err)
	}
	return count, nil
}

func (r *repo) GetComments(ctx context.Context, postID int) ([]Comment, error) {
	query := `SELECT id, post_id, user_id, content FROM comments WHERE post_id = $1`
	rows, err := r.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// get user
		user, _ := r.GetUserByID(ctx, comment.UserID)
		comment.UserName = user.FirstName + " " + user.LastName
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *repo) LikeToggleBlog(ctx context.Context, userID int, postID int) (bool, int, error) {
	var exists bool
	checkErr := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id=$1 AND post_id=$2)", userID, postID).Scan(&exists)
	if checkErr != nil {
		return false, 0, fmt.Errorf("failed to check like: %w", checkErr)
	}
	var liked bool
	if exists {
		query := `DELETE FROM likes WHERE user_id=$1 AND post_id=$2`
		_, deleteErr := r.db.Exec(ctx, query, userID, postID)
		if deleteErr != nil {
			return false, 0, fmt.Errorf("failed to unlike blog: %w", deleteErr)
		}
	} else {
		query := `INSERT INTO likes (post_id, user_id) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING id;`
		var id int
		insertErr := r.db.QueryRow(ctx, query, postID, userID).Scan(&id)
		if insertErr != nil {
			return false, 0, fmt.Errorf("failed to like blog: %w", insertErr)
		}
		liked = true
	}

	count, err := r.GetLikes(ctx, postID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get likes: %w", err)
	}
	return liked, count, nil
}

func (r *repo) CommentBlog(ctx context.Context, comment Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}
	return nil
}

func (r *repo) GetBlogWithUserLikeStatus(ctx context.Context, userID int, id int) (Blog, error) {
	query := `
			SELECT 
				p.id, p.title, p.summary, p.content, p.author, p.posted_at,
				CASE 
					WHEN l.user_id IS NOT NULL THEN TRUE 
					ELSE FALSE 
				END AS is_liked_by_user
			FROM posts as p 
			LEFT JOIN likes as l 
				ON p.id = l.post_id 
				AND l.user_id = $1
			WHERE p.id= $2;`
	var blog Blog
	err := r.db.QueryRow(ctx, query, userID, id).Scan(&blog.ID, &blog.Title, &blog.Summary, &blog.Content, &blog.Author, &blog.PostedAt, &blog.LikedByUser)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Blog{}, fmt.Errorf("blog not found")
		}
		return Blog{}, fmt.Errorf("failed to query blog: %w", err)
	}
	blog.Date = blog.PostedAt.Format("January 2, 2006")

	// Get comments
	comments, err := r.GetComments(ctx, id)
	if err != nil {
		return Blog{}, fmt.Errorf("failed to get comments: %w", err)
	}
	blog.Comments = comments

	// Get likes
	count, _ := r.GetLikes(ctx, id)
	blog.Likes = count
	return blog, nil
}
