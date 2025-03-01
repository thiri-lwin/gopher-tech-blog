package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	mw "github.com/thiri-lwin/gopher-tech-blog/internal/middleware"
	repo "github.com/thiri-lwin/gopher-tech-blog/internal/repo/postgres"
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Phone   string `json:"phone"`
}

type CommentResp struct {
	Content  string `json:"content"`
	UserName string `json:"user_name"`
}

// GetPosts handles displaying all posts
func (h Handler) GetPosts(c *gin.Context) {
	auth := mw.GetRequestMeta(c)
	ctx := c.Request.Context()
	page, offset := h.getPaginationParams(c)

	blogs, err := h.repo.GetBlogs(ctx, h.cfg.PostLimit, offset)
	if err != nil {
		log.Println("Failed to get posts:", err)
		h.renderError(c, "Something went wrong. Please try again later.")
		return
	}

	prevPage, nextPage := h.getPaginationLinks(ctx, page, offset)

	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/home-bg.jpg", h.cfg.ImageURL),
		"Heading":         "Gopher Blog",
		"Subheading":      "Tech Journal by A Gopher",
		"posts":           blogs,
		"PrevPage":        prevPage,
		"NextPage":        nextPage,
		"IsAuthenticated": auth.UserID != 0,
	}

	// Render the index template and pass the posts to it
	h.tmpl.ExecuteTemplate(c.Writer, "index.html", data)
}

func (h Handler) getPaginationParams(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * h.cfg.PostLimit
	return page, offset
}

func (h Handler) getPaginationLinks(ctx context.Context, page, offset int) (int, int) {
	var prevPage, nextPage int
	if page > 1 {
		prevPage = page - 1
	}

	count, err := h.repo.GetBlogsCount(ctx)
	if err != nil {
		return prevPage, nextPage
	}
	if count > int64(offset+h.cfg.PostLimit) {
		nextPage = page + 1
	}

	return prevPage, nextPage
}

// GetPost handles displaying a single post based on its ID
func (h Handler) GetPost(c *gin.Context) {
	auth := mw.GetRequestMeta(c)
	// Get the post ID from the URL parameter
	postID := c.Param("id")
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		log.Println("Failed to convert post ID to integer:", err)
		h.renderError(c, "Post not found.")
		return
	}
	post, err := h.repo.GetBlog(c.Request.Context(), postIDInt)
	if err != nil {
		log.Println("Failed to get post:", err)
		h.renderError(c, "Post not found.")
		return
	}
	// Render the post template with the post data
	h.tmpl.ExecuteTemplate(c.Writer, "post.html", gin.H{
		"BackgroundImage": fmt.Sprintf("%s/post-bg.jpg", h.cfg.ImageURL),
		"post":            post,
		"UserName":        auth.FirstName + " " + auth.LastName,
		"IsAuthenticated": auth.UserID != 0,
	})
}

// ServeAbout serves the about page
func (h Handler) ServeAbout(c *gin.Context) {
	auth := mw.GetRequestMeta(c)
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/about-bg.jpg", h.cfg.ImageURL),
		"Heading":         "About Me",
		"Subheading":      "This is what I do.",
		"IsAuthenticated": auth.UserID != 0,
	}
	h.tmpl.ExecuteTemplate(c.Writer, "about.html", data)
}

// ServeContact serves the contact page
func (h Handler) ServeContact(c *gin.Context) {
	auth := mw.GetRequestMeta(c)
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/contact-bg.jpg", h.cfg.ImageURL),
		"Heading":         "Contact Me",
		"Subheading":      "Have questions? I have answers (maybe).",
		"IsAuthenticated": auth.UserID != 0,
	}
	h.tmpl.ExecuteTemplate(c.Writer, "contact.html", data)
}

func (h Handler) renderError(c *gin.Context, errorMessage string) {
	auth := mw.GetRequestMeta(c)
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/error-bg.jpg", h.cfg.ImageURL),
		"Heading":         "Error",
		"Subheading":      errorMessage,
		"Status":          "Our team is working to resolve the issue. Please try again later.",
		"IsAuthenticated": auth.UserID != 0,
	}
	h.tmpl.ExecuteTemplate(c.Writer, "error.html", data)
}

func (h Handler) SendMessage(c *gin.Context) {
	var form ContactForm
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println("Failed to bind form data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Send the message to the admin
	if err := h.emailSender.SendEmail(h.cfg.AdminEmail, "Gopher Tech Blog Contact Form Submission", fmt.Sprintf("Name: %s\nEmail: %s\nPhone: %s\nMessage: %s", form.Name, form.Email, form.Phone, form.Message)); err != nil {
		log.Println("Failed to send email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.Redirect(http.StatusFound, "/contact")
}

// LikePost handles liking a post
func (h Handler) LikePost(c *gin.Context) {
	auth := mw.GetRequestMeta(c)
	if auth.UserID <= 0 {
		log.Println("unauthenticated user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please Signin to like the post"})
	}

	postID := c.Param("id")
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		log.Println("Failed to convert post ID to integer:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	likes, err := h.repo.LikeBlog(c.Request.Context(), auth.UserID, postIDInt)
	if err != nil {
		log.Println("Failed to like post:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"likes": likes})
}

// CommentPost handles commenting on a post
func (h Handler) CommentPost(c *gin.Context) {
	auth := mw.GetRequestMeta(c)
	if auth.UserID <= 0 {
		log.Println("unauthenticated user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please Signin to comment under the post"})
	}

	postID := c.Param("id")
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var comment repo.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Println("invalid json: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	comment.PostID = postIDInt
	comment.UserID = auth.UserID
	err = h.repo.CommentBlog(c.Request.Context(), comment)
	if err != nil {
		log.Println("Failed to add comment: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}
	cmtRes := CommentResp{
		Content:  comment.Content,
		UserName: auth.FirstName + " " + auth.LastName,
	}
	c.JSON(http.StatusOK, cmtRes)
}
