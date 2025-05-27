package articles

import (
	"context"

	"emperror.dev/errors"
	"github.com/ronnyp07/SportStream/worker/internal/domain/models"
	"github.com/ronnyp07/SportStream/worker/internal/domain/ports/repos"
)

type articles struct {
	repo repos.IArticlesRepos
}

func NewArticlesService(repo repos.IArticlesRepos) *articles {
	return &articles{
		repo: repo,
	}
}

func (a articles) UpsertByExternalID(ctx context.Context,
	articles []models.UpsertArticle) (models.Article, error) {
	var result models.Article

	for _, article := range articles {
		article.ExternalID = article.ID
		result, err := a.repo.UpsertByExternalID(ctx, article)
		if err != nil {
			return result, errors.Wrap(err, "unable to upsert article")
		}
	}

	return result, nil
}
