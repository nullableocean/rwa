package handlers

import (
	"encoding/json"
	"net/http"
	"rwa/internal/models"
	"rwa/internal/services"
)

type ArticleHandler struct {
	as *services.ArticleService
	us *services.UserService
}

func NewArticleHandler(as *services.ArticleService, us *services.UserService) *ArticleHandler {
	return &ArticleHandler{
		as: as,
		us: us,
	}
}

func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	articleInfo := models.ArticleInfo{}
	err := json.NewDecoder(r.Body).Decode(articleInfo)
	if err != nil {
		badJsonError(w)
		return
	}

	uId, err := GetUserIdFromRequestCtx(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := h.us.GetUserById(uId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	article, err := h.as.CreateArticle(*user, articleInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

func (h *ArticleHandler) Get(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	tags := vals["tag"]
	username := vals.Get("author")

	var articles []*models.Article
	var err error

	if username != "" {
		user, err := h.us.GetUserByUsername(username)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))
			return
		}

		articles, err = h.as.GetAllByUser(*user, tags)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		articles, err = h.as.GetAll(tags)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}
