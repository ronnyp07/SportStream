package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ronnyp07/SportStream/api/internal/domain/models"
	services "github.com/ronnyp07/SportStream/api/internal/domain/services/articles"
	repomocks "github.com/ronnyp07/SportStream/api/tests/mocks/repos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticleService_GetArticleByID(t *testing.T) {
	t.Parallel()

	// Test cases
	tests := []struct {
		name          string
		id            int
		mockSetup     func(*repomocks.MockIArticlesRepos)
		expectedError string
	}{
		{
			name: "success - valid ID",
			id:   123,
			mockSetup: func(m *repomocks.MockIArticlesRepos) {
				m.EXPECT().GetByID(gomock.Any(), 123).
					Return(&models.Article{ID: 123}, nil)
			},
		},
		{
			name:          "error - invalid ID",
			id:            0,
			expectedError: "invalid article ID",
		},
		{
			name: "error - repository error",
			id:   123,
			mockSetup: func(m *repomocks.MockIArticlesRepos) {
				m.EXPECT().GetByID(gomock.Any(), 123).
					Return(nil, assert.AnError)
			},
			expectedError: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repomocks.NewMockIArticlesRepos(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := services.NewArticleService(mockRepo)

			// Execute
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			article, err := service.GetArticleByID(ctx, tt.id)

			// Verify
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, article)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, tt.id, article.ID)
			}
		})
	}
}

func TestArticleService_GetArticleByExternalID(t *testing.T) {
	t.Parallel()

	// Test cases
	tests := []struct {
		name          string
		externalID    int
		mockSetup     func(*repomocks.MockIArticlesRepos)
		expectedError string
	}{
		{
			name:       "success - valid external ID",
			externalID: 456,
			mockSetup: func(m *repomocks.MockIArticlesRepos) {
				m.EXPECT().GetByExternalID(gomock.Any(), 456).
					Return(&models.Article{ID: 1, ExternalID: 456}, nil)
			},
		},
		{
			name:          "error - invalid external ID",
			externalID:    0,
			expectedError: "invalid external ID",
		},
		{
			name:       "error - repository error",
			externalID: 456,
			mockSetup: func(m *repomocks.MockIArticlesRepos) {
				m.EXPECT().GetByExternalID(gomock.Any(), 456).
					Return(nil, assert.AnError)
			},
			expectedError: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repomocks.NewMockIArticlesRepos(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := services.NewArticleService(mockRepo)

			// Execute
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			article, err := service.GetArticleByExternalID(ctx, tt.externalID)

			// Verify
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, article)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, tt.externalID, article.ExternalID)
			}
		})
	}
}

func TestArticleService_GetPaginatedArticles(t *testing.T) {
	t.Parallel()

	// Test data
	mockArticles := &models.PaginatedArticles{
		PageInfo: models.PageInfo{
			Page:       1,
			NumPages:   1,
			PageSize:   20,
			NumEntries: 5,
		},
		Content: []models.Article{
			{ID: 1}, {ID: 2}, {ID: 3},
		},
	}

	// Test cases
	tests := []struct {
		name          string
		page          int
		pageSize      int
		mockSetup     func(*repomocks.MockIArticlesRepos)
		expectedPage  int
		expectedSize  int
		expectedError string
	}{
		{
			name:     "success - default pagination",
			page:     0,
			pageSize: 0,
			mockSetup: func(m *repomocks.MockIArticlesRepos) {
				m.EXPECT().GetPaginatedArticles(gomock.Any(), 1, 20).
					Return(mockArticles, nil)
			},
			expectedPage: 1,
			expectedSize: 20,
		},
		{
			name:     "error - repository error",
			page:     1,
			pageSize: 20,
			mockSetup: func(m *repomocks.MockIArticlesRepos) {
				m.EXPECT().GetPaginatedArticles(gomock.Any(), 1, 20).
					Return(nil, assert.AnError)
			},
			expectedError: assert.AnError.Error(),
		},
		{
			name:         "invalid - negative page",
			page:         -1,
			pageSize:     20,
			expectedPage: 1,
			expectedSize: 20,
			mockSetup: func(m *repomocks.MockIArticlesRepos) {
				m.EXPECT().GetPaginatedArticles(gomock.Any(), 1, 20).
					Return(mockArticles, nil)
			},
		},
		{
			name:         "invalid - page size too large",
			page:         1,
			pageSize:     150,
			expectedPage: 1,
			expectedSize: 20,
			mockSetup: func(m *repomocks.MockIArticlesRepos) {
				m.EXPECT().GetPaginatedArticles(gomock.Any(), 1, 20).
					Return(mockArticles, nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repomocks.NewMockIArticlesRepos(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := services.NewArticleService(mockRepo)

			// Execute
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := service.GetPaginatedArticles(ctx, tt.page, tt.pageSize)

			// Verify
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedPage, result.PageInfo.Page)
				assert.Equal(t, tt.expectedSize, result.PageInfo.PageSize)
				assert.Len(t, result.Content, len(mockArticles.Content))
			}
		})
	}
}
