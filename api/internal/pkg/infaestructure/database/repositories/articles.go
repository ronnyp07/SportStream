package repositories

import (
	"context"

	"github.com/pkg/errors"
	"github.com/ronnyp07/SportStream/api/internal/domain/models"
	"github.com/ronnyp07/SportStream/api/internal/domain/ports/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "articles"
)

type ArticleRepository struct {
	collection *mongo.Collection
	metrics    metrics.MetricsHandler
}

func NewArticleRepository(db *mongo.Database,
	metrics metrics.MetricsHandler) *ArticleRepository {
	return &ArticleRepository{
		collection: db.Collection(collectionName),
		metrics:    metrics,
	}
}

func (r *ArticleRepository) GetByID(ctx context.Context, id int) (*models.Article, error) {
	r.metrics.DBCall("GetByID")

	var article models.Article
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&article)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.metrics.DBErrorInc("GetByID", "not_found")
			return nil, errors.New("article not found")
		}
		r.metrics.DBErrorInc("GetByID", err.Error())
		return nil, err
	}

	return &article, nil
}

func (r *ArticleRepository) GetByExternalID(ctx context.Context, externalID int) (*models.Article, error) {
	r.metrics.DBCall("GetByExternalID")

	var article models.Article
	err := r.collection.FindOne(ctx, bson.M{"externalID": externalID}).Decode(&article)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.metrics.DBErrorInc("GetByExternalID", "not_found")
			return nil, errors.New("article not found")
		}
		r.metrics.DBErrorInc("GetByExternalID", err.Error())
		return nil, err
	}

	return &article, nil
}

func (r *ArticleRepository) GetPaginatedArticles(ctx context.Context, page, pageSize int) (*models.PaginatedArticles, error) {
	r.metrics.DBCall("GetPaginatedArticles")

	skip := (page - 1) * pageSize

	// Get total count of articles
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.metrics.DBErrorInc("GetPaginatedArticles", "count_error")
		return nil, errors.Wrap(err, "failed to count articles")
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "date", Value: -1}})

	// Execute query
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		r.metrics.DBErrorInc("GetPaginatedArticles", "find_error")
		return nil, errors.Wrap(err, "failed to find articles")
	}
	defer cursor.Close(ctx)

	var articles []models.Article
	if err := cursor.All(ctx, &articles); err != nil {
		r.metrics.DBErrorInc("GetPaginatedArticles", "decode_error")
		return nil, errors.Wrap(err, "failed to decode articles")
	}

	return &models.PaginatedArticles{
		PageInfo: models.PageInfo{
			Page:       page,
			NumPages:   totalPages,
			PageSize:   pageSize,
			NumEntries: int(total),
		},
		Content: articles,
	}, nil
}
