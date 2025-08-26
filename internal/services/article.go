package services

import (
	"errors"
	"fmt"
	"regexp"
	"rwa/internal/models"
	"strings"
	"time"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

type ArticleRepository interface {
	GetBySlug(string) (*models.Article, error)
	GetAllByUser(user models.User, tags []string) ([]*models.Article, error)
	GetAll(tags []string) ([]*models.Article, error)
	Save(models.Article) error
	Delete(models.Article) error
	Update(string, models.Article) error
}

type ArticleService struct {
	articleRepo ArticleRepository
}

func NewArticleService(articleRepo ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
	}
}

func (as *ArticleService) CreateArticle(user models.User, articleInfo models.ArticleInfo) (*models.Article, error) {
	if articleInfo.Slug != "" {
		articleBySlug, _ := as.articleRepo.GetBySlug(articleInfo.Slug)
		if articleBySlug != nil {
			return nil, errors.New("slug must be unique")
		}
	} else {
		articleInfo.Slug = as.generateSlug(articleInfo)
	}

	createdAt := time.Now()
	article := models.Article{
		Author:         user,
		Body:           articleInfo.Body,
		Title:          articleInfo.Title,
		Description:    articleInfo.Description,
		Favorited:      false,
		FavoritesCount: 0,
		Slug:           articleInfo.Slug,
		TagList:        articleInfo.TagList,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
	}

	err := as.articleRepo.Save(article)
	if err != nil {
		fmt.Println(err)

		return nil, errors.New("error: cannot save article")
	}

	return &article, err
}

func (as *ArticleService) UpdateArticle(user models.User, article models.Article, articleInfo models.ArticleInfo) (*models.Article, error) {
	if article.Author.ID != user.ID {
		return nil, errors.New("error: not permission")
	}
	oldSlug := article.Slug

	if articleInfo.Slug != "" {
		article.Slug = articleInfo.Slug
	}

	article.Title = articleInfo.Title
	article.Body = articleInfo.Body
	article.Description = articleInfo.Description
	article.TagList = articleInfo.TagList
	article.UpdatedAt = time.Now()

	err := as.articleRepo.Update(oldSlug, article)
	if err != nil {
		fmt.Println(err)

		return nil, errors.New("error: cannot update article")
	}

	return &article, err
}

func (as *ArticleService) DeleteArticle(user models.User, article models.Article) error {
	if article.Author.ID != user.ID {
		return errors.New("error: not permission")
	}

	err := as.articleRepo.Delete(article)
	if err != nil {
		fmt.Println(err)

		return errors.New("error: cannot delete article")
	}

	return nil
}

func (as *ArticleService) GetAll(tags []string) ([]*models.Article, error) {
	return as.articleRepo.GetAll(tags)
}

func (as *ArticleService) GetByTitle(title string) (*models.Article, error) {
	article, err := as.articleRepo.GetBySlug(title)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (as *ArticleService) GetAllByUser(user models.User, tags []string) ([]*models.Article, error) {
	articles, err := as.articleRepo.GetAllByUser(user, tags)
	if err != nil {
		return articles, err
	}

	return articles, nil
}

func (as *ArticleService) generateSlug(articleInfo models.ArticleInfo) string {
	title := articleInfo.Title

	return strings.Join(strings.Split(strings.ToLower(nonAlphanumericRegex.ReplaceAllString(title, "")), " "), "-")
}
