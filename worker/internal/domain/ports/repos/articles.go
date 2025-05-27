package repos

import (
	"context"

	"github.com/ronnyp07/SportStream/worker/internal/domain/models"
)

type IArticlesRepos interface {
	UpsertByExternalID(ctx context.Context, article models.UpsertArticle) (models.Article, error)
}
