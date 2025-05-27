package services

import (
	"context"

	"github.com/ronnyp07/SportStream/worker/internal/domain/models"
)

type IArticlesService interface {
	UpsertByExternalID(ctx context.Context,
		articles []models.UpsertArticle) (models.Article, error)
}
