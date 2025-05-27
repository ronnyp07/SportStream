package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/ronnyp07/SportStream/api/docs"
	"github.com/ronnyp07/SportStream/api/internal/app/httpserver/handler"
	portsHandler "github.com/ronnyp07/SportStream/api/internal/domain/ports/handler"
	"github.com/ronnyp07/SportStream/api/internal/pkg/config"
	"github.com/ronnyp07/SportStream/api/internal/pkg/infaestructure/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Metrics middleware wrapper
type metricsMiddleware struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

// Response writer wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewMetricsMiddleware() *metricsMiddleware {
	return &metricsMiddleware{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),
	}
}

func (m *metricsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		status := strconv.Itoa(rw.statusCode)

		m.requestsTotal.WithLabelValues(r.Method, path, status).Inc()
		m.requestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

func NewServerBuilder(services Services) *Server {

	return &Server{
		Services:   services,
		httpServer: &http.Server{},
	}
}

func (s *Server) WithAddr(addr string) *Server {
	s.httpServer.Addr = addr
	return s
}

func (s *Server) WithReadTimeout(readTimeOut time.Duration) *Server {
	s.httpServer.ReadTimeout = readTimeOut
	return s
}

func (s *Server) WithWriteTimeout(writeTimeOut time.Duration) *Server {
	s.httpServer.WriteTimeout = writeTimeOut
	return s
}

func (s *Server) Build() *Server {
	s.setupHandler()
	s.httpServer.Handler = s.Routes
	return s
}

func (s *Server) Start(ctx context.Context, shutDownCall func()) {
	go func() {
		log.Logger().Info(ctx, fmt.Sprintf("Server starting on %s", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger().Error(ctx, fmt.Sprintf("Application services started successfully in %s", config.App().Env.Name))
			shutDownCall()
		}
	}()
}

func (s *Server) setupHandler() {

	// Create metrics middleware
	metricsMiddleware := NewMetricsMiddleware()

	// Register metrics with Prometheus
	prometheus.MustRegister(metricsMiddleware.requestsTotal)
	prometheus.MustRegister(metricsMiddleware.requestDuration)

	articleHandler := handler.NewArticleHandler(s.Services.ArticleService)
	router := NewRouter(articleHandler, metricsMiddleware)
	s.Routes = router

}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("Application services started successf)ully in %v", err))
		return err
	}
	return nil
}

func NewRouter(articleHandler portsHandler.IHandler, metrics *metricsMiddleware) *mux.Router {
	r := mux.NewRouter()

	// Apply metrics middleware to all routes
	r.Use(metrics.Handler)

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// API versioning
	api := r.PathPrefix("/api/v1").Subrouter()

	// Article routes
	api.HandleFunc("/articles/{id:[0-9]+}", articleHandler.GetArticleByID).Methods("GET")
	api.HandleFunc("/articles/external/{externalID:[0-9]+}", articleHandler.GetArticleByExternalID).Methods("GET")
	api.HandleFunc("/articles", articleHandler.GetPaginatedArticles).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return r
}
