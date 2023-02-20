package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/3d0c/storage/pkg/log"
)

// Middlewares type
type Middlewares func(res http.ResponseWriter, request *http.Request) (int, error)

func Chain(m ...Middlewares) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			status int
		)

		for _, middleware := range m {
			if status, err = middleware(w, r); err != nil {
				w.WriteHeader(status)
				log.TheLogger().Error("HTTP Request",
					zap.Error(err), zap.String("method", r.Method), zap.String("source", r.RemoteAddr), zap.String("URI", r.RequestURI))
				break
			}
		}
	}
}
