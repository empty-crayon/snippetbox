package main

import "github.com/empty-crayon/snippetbox/internal/models"

type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
}