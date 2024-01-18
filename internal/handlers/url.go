package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"restapi/internal/domain"
	"restapi/internal/repository"
	"restapi/internal/services"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	service *services.URLService
	log     *slog.Logger
}

type Request struct {
}

func NewURLHandler(s *services.URLService, log *slog.Logger) *URLHandler {
	return &URLHandler{
		service: s,
		log:     log,
	}
}

func (h URLHandler) SaveURL(w http.ResponseWriter, r *http.Request) {
	//TODO: взять id пользователя из контекста запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	url := domain.URL{}
	err = json.Unmarshal(b, &url)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	alias, err := h.service.SaveURL(url.URL, url.Alias, url.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrURLExists) {
			jsonError(w, http.StatusBadRequest, "Alias URL must be unique", h.log)
			return
		}
		jsonError(w, http.StatusBadRequest, "Unknown error", h.log)
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"alias":"%s"}`, alias)))
}

func (h URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	h.log.Info("request")
	if alias == "" {
		jsonError(w, http.StatusBadRequest, "invalid request", h.log)
		return
	}
	url, err := h.service.GetURL(alias)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			jsonError(w, http.StatusBadRequest, "URL not found", h.log)
			return
		}
		jsonError(w, http.StatusBadRequest, "Unknown error", h.log)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func jsonError(w http.ResponseWriter, code int, msg string, log *slog.Logger) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, msg)))
	if err != nil {
		log.Error(err.Error())
	}
}
