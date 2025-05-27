package repos

import (
	"context"

	"github.com/ronnyp07/SportStream/api/internal/domain/models"
)

type IArticlesRepos interface {
	GetByID(ctx context.Context, id int) (*models.Article, error)
	GetByExternalID(ctx context.Context, externalID int) (*models.Article, error)
	GetPaginatedArticles(ctx context.Context, page, pageSize int) (*models.PaginatedArticles, error)
}
