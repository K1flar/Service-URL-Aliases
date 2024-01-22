package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"restapi/internal/handlers"
	"restapi/internal/repository"

	"github.com/go-chi/chi/v5"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name=URLGetterService
type URLGetterService interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, s URLGetterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		log.Info("request")
		if alias == "" {
			handlers.JSONError(w, http.StatusBadRequest, "invalid request", log)
			return
		}
		url, err := s.GetURL(alias)
		if err != nil {
			if errors.Is(err, repository.ErrURLNotFound) {
				handlers.JSONError(w, http.StatusBadRequest, "URL not found", log)
				return
			}
			handlers.JSONError(w, http.StatusBadRequest, "Unknown error", log)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}
