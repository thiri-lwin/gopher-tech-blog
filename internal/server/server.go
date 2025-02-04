package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/thiri-lwin/gopher-tech-blog/internal/config"
	"github.com/thiri-lwin/gopher-tech-blog/internal/handlers"
	"github.com/thiri-lwin/gopher-tech-blog/internal/pkg/mailsender"
	"github.com/thiri-lwin/gopher-tech-blog/internal/repo"
	"github.com/thiri-lwin/gopher-tech-blog/internal/repo/redis"
)

var tmpl *template.Template

func New(cfg *config.Config) *gin.Engine {
	router := gin.New()

	err := loadTemplates("templates")
	if err != nil {
		log.Fatal(err)
	}

	router.Static("/static", "./static")
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// init redis
	redis.InitRedis(cfg.RedisAddr, cfg.RedisUser, cfg.RedisPass)

	router.Use(rateLimitMW)

	router.OPTIONS("/*any", responseOK())

	// Initialize the database
	db := repo.New(cfg.DatabaseURI)

	// init email sender
	emailSender := mailsender.NewEmailSender(cfg.SMTPServer, cfg.SMTPPort, cfg.EmailFrom, cfg.EmailPass)

	postHandler := handlers.NewHandler(cfg, db, emailSender, tmpl)
	router.GET("/", postHandler.GetPosts)
	router.GET("/index", postHandler.GetPosts)                               // Home page route
	router.GET("/about", postHandler.ServeAbout)                             // About page route
	router.GET("/contact", postHandler.ServeContact)                         // Contact page route
	router.GET("/posts/:id", postHandler.GetPost)                            // Post page route
	router.POST("/contact", rateLimitSendMessageMW, postHandler.SendMessage) // Contact form submission route
	router.POST("/posts/:id/like", postHandler.LikePost)                     // Like post route
	router.POST("/posts/:id/comment", postHandler.CommentPost)               // Comment post route

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
func loadTemplates(templateDir string) error {
	// Define a function map for Go templates
	funcMap := template.FuncMap{
		"safeHTML": func(text string) template.HTML {
			return template.HTML(text)
		},
	}

	// Route to serve the page

	var err error
	// Glob to match all .html files under the template directory
	templates, err := filepath.Glob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return fmt.Errorf("failed to load main templates: %w", err)
	}

	// Parse all templates
	tmpl, err = template.New("").Funcs(funcMap).ParseFiles(templates...)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	return nil
}
