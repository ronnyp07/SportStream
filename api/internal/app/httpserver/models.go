package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
	portsServices "github.com/ronnyp07/SportStream/api/internal/domain/ports/services"
)

type Services struct {
	ArticleService portsServices.IArticlesService
}

type Server struct {
	Services   Services
	httpServer *http.Server
	Routes     *mux.Router
}
