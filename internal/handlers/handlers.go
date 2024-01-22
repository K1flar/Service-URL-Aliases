package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

func JSONError(w http.ResponseWriter, code int, msg string, log *slog.Logger) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, msg)))
	if err != nil {
		log.Error(err.Error())
	}
}
