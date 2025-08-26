package ram

import (
	"errors"
	"rwa/internal/models"
	"sync"
)

type ArticleRepository struct {
	store         map[string]models.Article
	usersArticles map[int64][]string
	tags          map[string][]string

	mu *sync.Mutex
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{
		store:         make(map[string]models.Article),
		usersArticles: make(map[int64][]string),
		tags:          make(map[string][]string),
		mu:            &sync.Mutex{},
	}
}

func (r *ArticleRepository) GetAll(tags []string) ([]*models.Article, error) {
	if len(tags) != 0 {
		return r.getAllByTags(tags), nil
	} else {
		return r.getAll(), nil
	}
}

func (r *ArticleRepository) getAllByTags(tags []string) []*models.Article {
	slugs := r.getSlugsMapByTags(tags)

	articles := make([]*models.Article, 0, len(slugs))
	for s := range slugs {
		articles = append(articles, r.getBySlug(s))
	}

	return articles
}

func (r *ArticleRepository) getAll() []*models.Article {
	articles := make([]*models.Article, 0, len(r.store))

	for _, a := range r.store {
		articles = append(articles, &a)
	}

	return articles
}

func (r *ArticleRepository) GetBySlug(slug string) (*models.Article, error) {
	if !r.isStored(slug) {
		return nil, models.ErrNotFound
	}

	return r.getBySlug(slug), nil
}

func (r *ArticleRepository) GetAllByUser(user models.User, tags []string) ([]*models.Article, error) {
	slugs, exist := r.usersArticles[user.ID]
	if !exist {
		return []*models.Article{}, models.ErrNotFound
	}

	useTagsFilter := len(tags) != 0
	var tagsSlugs map[string]bool
	if useTagsFilter {
		tagsSlugs = r.getSlugsMapByTags(tags)
	}

	articles := make([]*models.Article, 0, len(slugs))
	for _, s := range slugs {
		if useTagsFilter && !tagsSlugs[s] {
			continue
		}
		articles = append(articles, r.getBySlug(s))
	}

	return articles, nil
}

func (r *ArticleRepository) Save(article models.Article) error {
	r.mu.Lock()
	r.mu.Unlock()

	slug := article.Slug
	if r.isStored(slug) {
		return errors.New("slug must be unique: " + slug)
	}

	r.store[slug] = article
	r.saveSlugToUsersArticles(article.Author.ID, slug)
	r.saveTagsAndSlug(article.TagList, article.Slug)

	return nil
}

func (r *ArticleRepository) Delete(article models.Article) error {
	r.mu.Lock()
	r.mu.Unlock()

	slug := article.Slug
	userId := article.Author.ID

	if !r.isStored(slug) {
		return models.ErrNotFound
	}

	delete(r.store, slug)
	r.deleteSlugFromUsersArticles(userId, slug)

	return nil
}

func (r *ArticleRepository) Update(oldSlug string, article models.Article) error {
	r.mu.Lock()
	r.mu.Unlock()

	if !r.isStored(oldSlug) {
		return models.ErrNotFound
	}

	if article.Slug != oldSlug {
		if r.isStored(article.Slug) {
			return errors.New("slug must be unique: " + article.Slug)
		}

		r.Delete(*r.getBySlug(oldSlug))
		r.deleteSlugFromUsersArticles(article.Author.ID, oldSlug)

		r.Save(article)
	} else {
		r.store[oldSlug] = article
	}

	return nil
}

func (r *ArticleRepository) isStored(slug string) bool {
	_, ex := r.store[slug]

	return ex
}

func (r *ArticleRepository) getBySlug(slug string) *models.Article {
	article := r.store[slug]
	return &article
}

func (r *ArticleRepository) saveSlugToUsersArticles(userId int64, slug string) {
	_, ex := r.usersArticles[userId]
	if !ex {
		r.usersArticles[userId] = make([]string, 0)
	}

	r.usersArticles[userId] = append(r.usersArticles[userId], slug)
}

func (r *ArticleRepository) saveTagsAndSlug(tags []string, slug string) {
	for _, t := range tags {
		if _, ex := r.tags[t]; !ex {
			r.tags[t] = make([]string, 0)
		}

		r.tags[t] = append(r.tags[t], slug)
	}
}

func (r *ArticleRepository) deleteSlugFromUsersArticles(userId int64, slug string) {
	deletingInd := -1
	for i, s := range r.usersArticles[userId] {
		if s == slug {
			deletingInd = i
			break
		}
	}

	if deletingInd != -1 {
		if deletingInd == len(r.usersArticles[userId])-1 {
			r.usersArticles[userId] = r.usersArticles[userId][:deletingInd]
		} else if deletingInd == 0 {
			r.usersArticles[userId] = r.usersArticles[userId][deletingInd:]
		} else {
			r.usersArticles[userId] = append(r.usersArticles[userId][:deletingInd], r.usersArticles[userId][deletingInd+1:]...)
		}
	}
}

func (r *ArticleRepository) getSlugsMapByTags(tags []string) map[string]bool {
	slugs := make(map[string]bool)

	for _, t := range tags {
		for _, s := range r.tags[t] {
			slugs[s] = true
		}
	}

	return slugs
}
