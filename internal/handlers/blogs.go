package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Phone   string `json:"phone"`
}

// GetPosts handles displaying all posts
func (h Handler) GetPosts(c *gin.Context) {
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
	// Get the post ID from the URL parameter
	postID := c.Param("id")
	post, err := h.repo.GetBlog(c.Request.Context(), postID)
	if err != nil {
		log.Println("Failed to get post:", err)
		h.renderError(c, "Post not found.")
		return
	}
	// Render the post template with the post data
	h.tmpl.ExecuteTemplate(c.Writer, "post.html", gin.H{
		"BackgroundImage": fmt.Sprintf("%s/post-bg.jpg", h.cfg.ImageURL),
		"post":            post,
	})
}

// ServeAbout serves the about page
func (h Handler) ServeAbout(c *gin.Context) {
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/about-bg.jpg", h.cfg.ImageURL),
		"Heading":         "About Me",
		"Subheading":      "This is what I do.",
	}
	h.tmpl.ExecuteTemplate(c.Writer, "about.html", data)
}

// ServeContact serves the contact page
func (h Handler) ServeContact(c *gin.Context) {
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/contact-bg.jpg", h.cfg.ImageURL),
		"Heading":         "Contact Me",
		"Subheading":      "Have questions? I have answers (maybe).",
	}
	h.tmpl.ExecuteTemplate(c.Writer, "contact.html", data)
}

func (h Handler) renderError(c *gin.Context, errorMessage string) {
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/error-bg.jpg", h.cfg.ImageURL),
		"Heading":         "Error",
		"Subheading":      errorMessage,
		"Status":          "Our team is working to resolve the issue. Please try again later.",
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
