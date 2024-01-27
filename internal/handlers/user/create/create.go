package create

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"restapi/internal/domains"
	"restapi/internal/handlers"
	service "restapi/internal/services"
)

type UserCreatorService interface {
	CreateUser(user *domains.User) (string, error)
}

func New(log *slog.Logger, s UserCreatorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			handlers.JSONError(w, http.StatusInternalServerError, "unknown error", log)
			return
		}
		defer r.Body.Close()

		var user domains.User
		err = json.Unmarshal(b, &user)
		if err != nil {
			handlers.JSONError(w, http.StatusBadRequest, "bad request", log)
			return
		}
		if user.Login == "" || user.Password == "" || user.Email == "" {
			handlers.JSONError(w, http.StatusBadRequest, "invalid fields", log)
			return
		}

		token, err := s.CreateUser(&user)
		if err != nil {
			if errors.Is(err, service.ErrUserExists) {
				handlers.JSONError(w, http.StatusConflict, "user alredy exists", log)
				return
			}
			handlers.JSONError(w, http.StatusInternalServerError, "unknown error", log)
			return
		}

		b, err = json.Marshal(token)
		if err != nil {
			handlers.JSONError(w, http.StatusInternalServerError, "unknown error", log)
			return
		}
		w.Write(b)
	}
}
