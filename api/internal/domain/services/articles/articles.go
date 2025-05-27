package services

import (
	"context"
	"errors"

	"github.com/ronnyp07/SportStream/api/internal/domain/models"
	"github.com/ronnyp07/SportStream/api/internal/domain/ports/repos"
)

type ArticleService struct {
	repo repos.IArticlesRepos
}

func NewArticleService(repo repos.IArticlesRepos) *ArticleService {
	return &ArticleService{
		repo: repo,
	}
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id int) (*models.Article, error) {
	if id <= 0 {
		return nil, errors.New("invalid article ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ArticleService) GetArticleByExternalID(ctx context.Context, externalID int) (*models.Article, error) {
	if externalID <= 0 {
		return nil, errors.New("invalid external ID")
	}
	return s.repo.GetByExternalID(ctx, externalID)
}

func (s *ArticleService) GetPaginatedArticles(ctx context.Context, page, pageSize int) (*models.PaginatedArticles, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.GetPaginatedArticles(ctx, page, pageSize)
}
