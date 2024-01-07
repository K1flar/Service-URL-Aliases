package main

import (
	"fmt"
	"log/slog"
	"os"
	"restapi/internal/config"
	"restapi/internal/logger"
)

func main() {
	cfg, err := config.New("../configs/config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_ = cfg

	log := logger.New(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

}
