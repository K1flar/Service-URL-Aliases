package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"restapi/internal/config"
	"restapi/internal/handlers/redirect"
	"restapi/internal/handlers/url/save"
	"restapi/internal/logger"
	"restapi/internal/repository/sqlite"
	urlservice "restapi/internal/services/urlService"

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

	router := chi.NewRouter()

	router.Post("/url", save.New(log, urlService))
	router.Get("/{alias}", redirect.New(log, urlService))

	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
