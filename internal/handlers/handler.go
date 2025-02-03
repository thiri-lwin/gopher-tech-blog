package handlers

import (
	"html/template"

	"github.com/go-redis/redis/v8"
	"github.com/thiri-lwin/gopher-tech-blog/internal/config"
	emailsender "github.com/thiri-lwin/gopher-tech-blog/internal/pkg/mailsender"
	repo "github.com/thiri-lwin/gopher-tech-blog/internal/repo"
	rdb "github.com/thiri-lwin/gopher-tech-blog/internal/repo/redis"
)

type Handler struct {
	cfg         *config.Config
	repo        repo.Repo
	redisClient *redis.Client
	emailSender *emailsender.EmailSender
	tmpl        *template.Template
}

func NewHandler(cfg *config.Config, repo repo.Repo, emailSender *emailsender.EmailSender, tmpl *template.Template) Handler {
	return Handler{
		cfg:         cfg,
		repo:        repo,
		redisClient: rdb.RDB,
		emailSender: emailSender,
		tmpl:        tmpl,
	}
}
