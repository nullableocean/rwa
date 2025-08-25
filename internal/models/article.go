package models

import (
	"errors"
	"time"
)

type Article struct {
	Author         User      `json:"user"`
	Body           string    `json:"body"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Slug           string    `json:"slug"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ArticleInfo struct {
	Body        string   `json:"body"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	TagList     []string `json:"tagList"`
	Slug        string   `json:"slug"`
}

func (i *ArticleInfo) Validate() error {
	switch "" {
	case i.Body:
		return errors.New("Body is empty")
	case i.Title:
		return errors.New("Title is empty")
	case i.Slug:
		return errors.New("Slug is empty")
	}

	return nil
}
