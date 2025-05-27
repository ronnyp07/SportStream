package handler

import "net/http"

type IHandler interface {
	GetArticleByID(w http.ResponseWriter, r *http.Request)
	GetArticleByExternalID(w http.ResponseWriter, r *http.Request)
	GetPaginatedArticles(w http.ResponseWriter, r *http.Request)
}
