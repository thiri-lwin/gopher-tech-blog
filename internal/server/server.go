package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/thiri-lwin/gopher-tech-blog/internal/handlers"
	"github.com/thiri-lwin/gopher-tech-blog/internal/repo"
)

func New() *gin.Engine {
	router := gin.New()

	tmpl, err := loadTemplates("templates")
	if err != nil {
		log.Fatal(err)
	}

	router.Static("/static", "./static")
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.OPTIONS("/*any", responseOK())

	//public := router.Group("/public/v1")

	db := repo.New()
	postHandler := handlers.NewHandler(db, tmpl)
	router.GET("/", postHandler.GetPosts)
	router.GET("/post/:id", postHandler.GetPost) // Post page route

	return router
}

func responseOK() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

// Load all templates (including partials)
func loadTemplates(templateDir string) (*template.Template, error) {
	// Glob to match all .html files under the template directory
	templates, err := filepath.Glob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}

	// Load the main template and all partials
	tmpl := template.Must(template.New("").ParseFiles(templates...))

	// If you want to load partials separately (e.g., header, footer), you can do so.
	partialFiles, err := filepath.Glob(filepath.Join(templateDir, "partials", "*.html"))
	if err != nil {
		return nil, err
	}

	// Parse the partials and add them to the main template
	if len(partialFiles) > 0 {
		_, err = tmpl.ParseFiles(partialFiles...)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}
