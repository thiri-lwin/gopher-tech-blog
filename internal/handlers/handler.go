package handlers

import (
	"html/template"

	"github.com/thiri-lwin/gopher-tech-blog/internal/config"
	repo "github.com/thiri-lwin/gopher-tech-blog/internal/repo"
)

type Handler struct {
	cfg  *config.Config
	repo repo.Repo
	tmpl *template.Template
}

func NewHandler(cfg *config.Config, repo repo.Repo, tmpl *template.Template) Handler {
	return Handler{
		cfg:  cfg,
		repo: repo,
		tmpl: tmpl,
	}
}
