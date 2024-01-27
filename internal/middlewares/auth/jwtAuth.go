package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"restapi/internal/config"
	"restapi/internal/domains"
	"restapi/internal/middlewares"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func New(log *slog.Logger, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authParts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(authParts) < 2 {
				middlewares.JSONError(w, http.StatusUnauthorized, "unauthorized", log)
				return
			}

			inToken := authParts[1]
			token, err := jwt.Parse(inToken, func(t *jwt.Token) (interface{}, error) {
				if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || method.Alg() != "HS256" {
					return nil, fmt.Errorf("bad sign method")
				}
				return []byte(cfg.Server.Secret), nil
			})
			if err != nil || !token.Valid {
				middlewares.JSONError(w, http.StatusUnauthorized, "bad token", log)
				return
			}

			payload, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				middlewares.JSONError(w, http.StatusUnauthorized, "no payload", log)
				return
			}

			userStr, err := json.Marshal(payload)
			if err != nil {
				middlewares.JSONError(w, http.StatusInternalServerError, "unknown error", log)
				return
			}

			var userStruct domains.User
			err = json.Unmarshal(userStr, &userStruct)
			if err != nil {
				middlewares.JSONError(w, http.StatusInternalServerError, "unknown error", log)
				return
			}

			ctx := context.WithValue(r.Context(), "user", userStruct)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
