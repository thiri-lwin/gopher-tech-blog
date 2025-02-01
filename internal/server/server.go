package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/thiri-lwin/gopher-tech-blog/internal/handlers"
	"github.com/thiri-lwin/gopher-tech-blog/internal/repo"
)

func New(dbURI string) *gin.Engine {
	router := gin.New()

	tmpl, err := loadTemplates("templates")
	if err != nil {
		log.Fatal(err)
	}

	router.Static("/static", "./static")
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.OPTIONS("/*any", responseOK())

	db := repo.New(dbURI)
	postHandler := handlers.NewHandler(db, tmpl)
	router.GET("/", postHandler.GetPosts)
	router.GET("/index", postHandler.GetPosts)       // Home page route
	router.GET("/about", postHandler.ServeAbout)     // About page route
	router.GET("/contact", postHandler.ServeContact) // Contact page route
	router.GET("/post/:id", postHandler.GetPost)     // Post page route

	return router
}

func responseOK() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

// Load all templates
func loadTemplates(templateDir string) (*template.Template, error) {
	// Glob to match all .html files under the template directory
	templates, err := filepath.Glob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to load main templates: %w", err)
	}

	// Parse all templates
	tmpl, err := template.New("").ParseFiles(templates...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return tmpl, nil
}
