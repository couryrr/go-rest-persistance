package internal

import (
	"log/slog"
	"net/http"
)

func Logger(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		slog.Info("Just logging")
		next.ServeHTTP(w, r)
	})
}
