package app

import (
	"github.com/ronnyp07/SportStream/api/internal/domain/ports/services"
)

type Services struct {
	ArticleServ services.IArticlesService
}
