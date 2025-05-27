package repositories

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/ronnyp07/SportStream/worker/internal/domain/models"
	"github.com/ronnyp07/SportStream/worker/internal/domain/ports/metrics"
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

const counterCollectionName = "article_counters"

func (r *ArticleRepository) getNextArticleID(ctx context.Context) (int, error) {
	r.metrics.DBCall("getNextArticleID")

	counterCollection := r.collection.Database().Collection(counterCollectionName)

	filter := bson.M{"_id": "article_id"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}

	err := counterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)

	if err != nil {
		r.metrics.DBErrorInc("getNextArticleID", err.Error())
		return 0, errors.Wrap(err, "failed to get next article ID")
	}

	return result.Seq, nil
}

func (r *ArticleRepository) UpsertByExternalID(ctx context.Context, article models.UpsertArticle) (models.Article, error) {
	r.metrics.DBCall("UpsertByExternalID")
	now := time.Now()
	var updatedArticle models.Article
	// Check if document exists
	var existing struct {
		ID int `bson:"id"`
	}
	err := r.collection.FindOne(
		ctx,
		bson.M{"externalID": article.ExternalID},
		options.FindOne().SetProjection(bson.M{"id": 1}),
	).Decode(&existing)

	var articleID int
	if err == mongo.ErrNoDocuments {
		articleID, err = r.getNextArticleID(ctx)
		if err != nil {
			r.metrics.DBErrorInc("UpsertByExternalID", "sequence_error")
			return updatedArticle, errors.Wrap(err, "failed to get next article ID")
		}
	} else if err != nil {
		r.metrics.DBErrorInc("UpsertByExternalID", "find_error")
		return updatedArticle, errors.Wrap(err, "failed to check existing article")
	} else {
		articleID = existing.ID
	}

	// Prepare full article
	fullArticle := models.Article{
		ID:          articleID,
		Title:       article.Title,
		Description: article.Description,
		Date:        article.Date,
		Body:        article.Body,
		Summary:     article.Summary,
		LeadMedia:   article.LeadMedia,
		Tags:        article.Tags,
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	filter := bson.M{"externalID": article.ExternalID}
	update := bson.M{
		"$set": fullArticle,
		"$setOnInsert": bson.M{
			"createdAt":  now,
			"externalID": article.ExternalID,
		},
	}

	err = r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedArticle)

	if err != nil {
		r.metrics.DBErrorInc("UpsertByExternalID", "upsert_error")
		if err == mongo.ErrNoDocuments {
			return updatedArticle, errors.New("failed to upsert article")
		}
		return updatedArticle, err
	}

	return updatedArticle, nil
}
