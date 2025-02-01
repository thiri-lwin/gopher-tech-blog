package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	limit = 10
)

func (h Handler) GetPosts(c *gin.Context) {
	ctx := c.Request.Context()
	blogs, err := h.repo.GetBlogs(ctx, limit, 0) // TODO: Implement pagination
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := gin.H{
		"BackgroundImage": "static/img/home-bg.jpg",
		"Heading":         "Gopher Blog",
		"Subheading":      "Tech Journal by A Gopher",
		"posts":           blogs,
	}
	// Render the index template and pass the posts to it
	h.tmpl.ExecuteTemplate(c.Writer, "index.html", data)
}

// GetPost handles displaying a single post based on its ID
func (h Handler) GetPost(c *gin.Context) {
	// Get the post ID from the URL parameter
	postID := c.Param("id")
	post, err := h.repo.GetBlog(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	// Render the post template with the post data
	h.tmpl.ExecuteTemplate(c.Writer, "post.html", gin.H{
		"BackgroundImage": "/static/img/post-bg.jpg",
		"post":            post,
	})
}

// ServeAbout serves the about page
func (h Handler) ServeAbout(c *gin.Context) {
	data := gin.H{
		"BackgroundImage": "static/img/about-bg.jpg",
		"Heading":         "About Me",
		"Subheading":      "This is what I do.",
	}
	h.tmpl.ExecuteTemplate(c.Writer, "about.html", data)
}

// ServeContact serves the contact page
func (h Handler) ServeContact(c *gin.Context) {
	data := gin.H{
		"BackgroundImage": "static/img/contact-bg.jpg",
		"Heading":         "Contact Me",
		"Subheading":      "Have questions? I have answers (maybe).",
	}
	h.tmpl.ExecuteTemplate(c.Writer, "contact.html", data)
}
