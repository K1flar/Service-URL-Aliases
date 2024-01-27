package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"restapi/internal/config"
	"restapi/internal/handlers/redirect"
	deleteurl "restapi/internal/handlers/url/delete"
	"restapi/internal/handlers/url/save"
	"restapi/internal/handlers/user/create"
	"restapi/internal/handlers/user/login"
	"restapi/internal/logger"
	"restapi/internal/middlewares/auth"
	"restapi/internal/repository/sqlite"
	urlservice "restapi/internal/services/urlService"
	userservice "restapi/internal/services/userService"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.New("../configs/config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log := logger.New(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	fmt.Println(cfg)
	repository, err := sqlite.New(cfg.Storage.Path)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	urlService := urlservice.NewURLService(repository, cfg, log)
	userService := userservice.NewUserService(repository, cfg, log)

	router := chi.NewRouter()

	router.Post("/login", login.New(log, userService))
	router.Post("/register", create.New(log, userService))

	router.Group(func(r chi.Router) {
		r.Use(auth.New(log, cfg))

		r.Post("/url", save.New(log, urlService))
		r.Delete("/{alias}", deleteurl.New(log, urlService))
	})

	router.Get("/{alias}", redirect.New(log, urlService))

	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
