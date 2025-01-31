package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetPosts(c *gin.Context) {
	ctx := c.Request.Context()
	blogs, err := h.repo.GetBlogs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Render the index template and pass the posts to it
	h.tmpl.ExecuteTemplate(c.Writer, "index.html", gin.H{
		"posts": blogs, // Assuming you have a list of posts
	})
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
		"post": post,
	})
}

// ServeAbout serves the about page
func (h Handler) ServeAbout(c *gin.Context) {
	h.tmpl.ExecuteTemplate(c.Writer, "about.html", nil)
}

// ServeContact serves the contact page
func (h Handler) ServeContact(c *gin.Context) {
	h.tmpl.ExecuteTemplate(c.Writer, "contact.html", nil)
}
