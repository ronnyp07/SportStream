package services

import "context"

type MessageQueueService interface {
	GetArticlesList(ctx context.Context, page int, limit int)
}

type MessageQueueProcessor interface {
	Start(ctx context.Context)
}
