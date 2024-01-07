package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Env    string `yaml:"env" env-default:"local"`
	Server Server `yaml:"server"`
}

type Server struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"8080"`
}

func New(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
