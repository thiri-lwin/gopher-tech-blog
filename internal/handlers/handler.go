package handlers

import (
	"html/template"

	"github.com/go-redis/redis/v8"
	"github.com/thiri-lwin/gopher-tech-blog/internal/config"
	repo "github.com/thiri-lwin/gopher-tech-blog/internal/repo"
	rdb "github.com/thiri-lwin/gopher-tech-blog/internal/repo/redis"
)

type Handler struct {
	cfg         *config.Config
	repo        repo.Repo
	redisClient *redis.Client
	tmpl        *template.Template
}

func NewHandler(cfg *config.Config, repo repo.Repo, tmpl *template.Template) Handler {
	return Handler{
		cfg:         cfg,
		repo:        repo,
		redisClient: rdb.RDB,
		tmpl:        tmpl,
	}
}
