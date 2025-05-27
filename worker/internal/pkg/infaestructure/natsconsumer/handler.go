package subcriptions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ronnyp07/SportStream/worker/internal/domain/models"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/log"
	"github.com/sts-solutions/base-code/ccmsgqueue"
)

type Handler struct {
	services *Service
}

func NewHandler(services *Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h Handler) HandleMessage(ctx context.Context, msg ccmsgqueue.ConsumeMessage) {
	log.Logger().Info(ctx, fmt.Sprintf("received message on subject: '%s': %s", msg.Subject(), string(msg.Data())))

	var articles []models.UpsertArticle

	err := json.Unmarshal(msg.Data(), &articles)
	if err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("error copying message: %v", err))
	}
	_, err = h.services.ArticleServ.UpsertByExternalID(ctx, articles)
	if err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("patching the message: %v", err))
	}
}
