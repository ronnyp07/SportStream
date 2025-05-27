package subcriptions

import (
	"github.com/ronnyp07/SportStream/worker/internal/domain/ports/services"
)

type Service struct {
	ArticleServ services.IArticlesService
}
