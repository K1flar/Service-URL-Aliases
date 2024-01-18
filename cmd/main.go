package main

import (
	"fmt"
	"log/slog"
	"os"
	"restapi/internal/config"
	"restapi/internal/logger"
	"restapi/internal/repository/sqlite"
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

	err = repository.SaveURL("https://www.youtube.com/watch?v=rCJvW2xgnk0&t=2482s", "test2", 1)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	urlTest, err := repository.GetURL("test2")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	fmt.Println(urlTest)

}
