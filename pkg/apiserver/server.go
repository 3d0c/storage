package apiserver

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/log"
)

// APIHTTPServer structure
type APIHTTPServer struct {
	cfg    config.ProxyConfig
	srv    *http.Server
	logger *zap.Logger
}

// NewAPIHTTPServer API server constructor
func NewAPIHTTPServer(cfg config.ProxyConfig) (*APIHTTPServer, error) {
	var (
		s = &APIHTTPServer{
			cfg:    cfg,
			logger: log.TheLogger().With(zap.String("component", "APIHTTPServer")),
		}
	)

	s.srv = &http.Server{
		Addr:         s.cfg.Address,
		Handler:      SetupRouter(s.cfg),
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTemout,
	}

	return s, nil
}

// Run starts HTTP server, ctx is used for server shutdown in case if ctx is closed
func (s *APIHTTPServer) Run(ctx context.Context) {
	loggerWithField := s.logger.With(zap.String("method", "Run"))

	go func() {
		for {
			<-ctx.Done()
			shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
			_ = s.srv.Shutdown(shutdownCtx)
			cancelFn()
			return
		}
	}()

	loggerWithField.With(zap.String("address", s.srv.Addr))

	if err := s.srv.ListenAndServe(); err != nil {
		loggerWithField.Warn("http server finished with error", zap.Error(err))
	}
}
