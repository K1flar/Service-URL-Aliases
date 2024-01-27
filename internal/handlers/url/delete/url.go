package deleteurl

import (
	"errors"
	"log/slog"
	"net/http"
	"restapi/internal/domains"
	"restapi/internal/handlers"
	service "restapi/internal/services"

	"github.com/go-chi/chi/v5"
)

type URLDeleterService interface {
	DeleteURL(alias string, userID uint32) error
}

func New(log *slog.Logger, s URLDeleterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(domains.User)
		if !ok {
			handlers.JSONError(w, http.StatusForbidden, "forbidden", log)
			return
		}
		alias := chi.URLParam(r, "alias")

		err := s.DeleteURL(alias, user.ID)
		if err != nil {
			if errors.Is(err, service.ErrURLNotFound) {
				handlers.JSONError(w, http.StatusBadRequest, "url not found", log)
				return
			}

			if errors.Is(err, service.ErrURLForbiddenToDelete) {
				handlers.JSONError(w, http.StatusForbidden, "forbidden to delete url", log)
				return
			}

			handlers.JSONError(w, http.StatusInternalServerError, "unknown error", log)
			return
		}
	}
}
