package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/ronnyp07/SportStream/api/docs"
	"github.com/ronnyp07/SportStream/api/internal/domain/ports/services"
)

type ArticleHandler struct {
	service services.IArticlesService
}

func NewArticleHandler(service services.IArticlesService) *ArticleHandler {
	return &ArticleHandler{
		service: service,
	}
}

// @title SportStream Articles API
// @version 1.0
// @description This is a sample server for managing sport articles.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// GetArticleByID godoc
// @Summary Get article by internal ID
// @Description Get article by internal auto-incremented ID
// @Tags articles
// @Accept  json
// @Produce  json
// @Param id path int true "Article ID"
// @Success 200 {object} models.Article
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := h.service.GetArticleByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// GetArticleByExternalID godoc
// @Summary Get article by external ID
// @Description Get article by external system ID
// @Tags articles
// @Accept  json
// @Produce  json
// @Param externalID path int true "External Article ID"
// @Success 200 {object} models.Article
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /articles/external/{externalID} [get]
func (h *ArticleHandler) GetArticleByExternalID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	externalID, err := strconv.Atoi(vars["externalID"])
	if err != nil {
		http.Error(w, "Invalid external ID", http.StatusBadRequest)
		return
	}

	article, err := h.service.GetArticleByExternalID(r.Context(), externalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// GetPaginatedArticles godoc
// @Summary Get paginated articles
// @Description Get articles with pagination support
// @Tags articles
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Items per page" default(20)
// @Success 200 {object} models.PaginatedArticles
// @Failure 500 {object} map[string]string
// @Router /articles [get]
func (h *ArticleHandler) GetPaginatedArticles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	result, err := h.service.GetPaginatedArticles(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
