package handlers

import (
	"html/template"

	repo "github.com/thiri-lwin/gopher-tech-blog/internal/repo"
)

type Handler struct {
	repo repo.Repo
	tmpl *template.Template
}

func NewHandler(repo repo.Repo, tmpl *template.Template) Handler {
	return Handler{
		repo: repo,
		tmpl: tmpl,
	}
}
