package save

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"restapi/internal/domain"
	"restapi/internal/handlers"
	"restapi/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name=URLSaverService
type URLSaverService interface {
	SaveURL(url, alias string, userID uint32) (string, error)
}

func New(log *slog.Logger, s URLSaverService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: взять id пользователя из контекста запроса
		b, err := io.ReadAll(r.Body)
		if err != nil {
			handlers.JSONError(w, http.StatusBadRequest, err.Error(), log)
			return
		}

		url := domain.URL{}
		err = json.Unmarshal(b, &url)
		if err != nil {
			handlers.JSONError(w, http.StatusBadRequest, err.Error(), log)
			return
		}

		alias, err := s.SaveURL(url.URL, url.Alias, url.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrURLExists) {
				handlers.JSONError(w, http.StatusBadRequest, "Alias URL must be unique", log)
				return
			}
			handlers.JSONError(w, http.StatusBadRequest, "Unknown error", log)
			return
		}

		w.Write([]byte(fmt.Sprintf(`{"alias":"%s"}`, alias)))
	}
}