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
	mw "github.com/thiri-lwin/gopher-tech-blog/internal/middleware"
	"github.com/thiri-lwin/gopher-tech-blog/internal/pkg/mailsender"
	repo "github.com/thiri-lwin/gopher-tech-blog/internal/repo/postgres"
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

	// router.Use(func(c *gin.Context) {
	// 	mw.RateLimitMW(c, tmpl)
	// })

	router.Use(mw.AuthMiddleware(cfg.JWTKey))

	router.OPTIONS("/*any", responseOK())

	// Initialize the database
	db := repo.New(cfg.DatabaseURI)

	// init email sender
	emailSender := mailsender.NewEmailSender(cfg.SMTPServer, cfg.SMTPPort, cfg.EmailFrom, cfg.EmailPass)

	handler := handlers.NewHandler(cfg, db, emailSender, tmpl)

	routerGroup := router.Group("/")
	routerGroup.Use(func(c *gin.Context) {
		mw.RateLimitMW(c, tmpl)
	})
	routerGroup.GET("/", handler.GetPosts)
	routerGroup.GET("/index", handler.GetPosts)       // Home page route
	routerGroup.GET("/about", handler.ServeAbout)     // About page route
	routerGroup.GET("/contact", handler.ServeContact) // Contact page route
	routerGroup.GET("/posts/:id", handler.GetPost)    // Post page route
	routerGroup.GET("/signin", handler.ServeSignIn)   // SignIn page
	routerGroup.GET("/signup", handler.ServeSignUp)   // SignUp page
	routerGroup.GET("/logout", handler.Logout)

	routerGroup.POST("/posts/:id/like-toggle", handler.TogglePostLike) // ToggleLike post route
	routerGroup.POST("/posts/:id/comment", handler.CommentPost)
	routerGroup.POST("/signin", handler.SignIn) // SignIn User
	routerGroup.POST("/signup", handler.SignUp)

	router.POST("/contact", mw.RateLimitSendMessageMW, handler.SendMessage) // Contact form submission route

	router.DELETE("comments/:id", handler.DeleteComment)

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
