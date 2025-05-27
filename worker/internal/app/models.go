package app

import (
	"github.com/ronnyp07/SportStream/worker/internal/domain/ports/services"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/msgqueue"
)

type Services struct {
	MsgQueueService msgqueue.MsgQueueService
	ArticleServ     services.IArticlesService
}
