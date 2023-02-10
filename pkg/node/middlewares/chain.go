package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/3d0c/storage/pkg/log"
)

// Middlewares type
type Middlewares func(res http.ResponseWriter, request *http.Request) (interface{}, int, error)

func Chain(m ...Middlewares) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			status int
			result interface{}
		)

		for _, middleware := range m {
			if result, status, err = middleware(w, r); err != nil {
				break
			}
		}

		w.WriteHeader(status)

		if err != nil {
			log.TheLogger().Error("HTTP Request",
				zap.Error(err), zap.String("method", r.Method), zap.String("source", r.RemoteAddr), zap.String("URI", r.RequestURI))
			return
		}

		if result != nil {
			if _, err := w.Write(result.([]byte)); err != nil {
				log.TheLogger().Error("HTTP Response error", zap.Error(err))
			}
		}
	}
}
