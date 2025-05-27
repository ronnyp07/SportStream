package services

import (
	"context"

	"github.com/ronnyp07/SportStream/api/internal/domain/models"
)

type IArticlesService interface {
	GetArticleByID(ctx context.Context, id int) (*models.Article, error)
	GetArticleByExternalID(ctx context.Context, externalID int) (*models.Article, error)
	GetPaginatedArticles(ctx context.Context, page, pageSize int) (*models.PaginatedArticles, error)
}
