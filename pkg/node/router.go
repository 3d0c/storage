package node

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/3d0c/storage/pkg/log"
	"github.com/3d0c/storage/pkg/node/handlers"
	"github.com/3d0c/storage/pkg/node/middlewares"
)

// SetupRouter sets up endpoints
func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Put(
		"/file/{ID}",
		middlewares.Chain(
			handlers.FileHandler().Put,
		),
	)

	r.Get(
		"/file/{ID}",
		middlewares.Chain(
			handlers.FileHandler().Get,
		),
	)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.TheLogger().Debug("registered", zap.String("method", method), zap.String("route", route))
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.TheLogger().Debug("logging error", zap.Error(err))
	}

	return r
}
