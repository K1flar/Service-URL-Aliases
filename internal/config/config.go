package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Env         string  `yaml:"env" env-default:"local"`
	Server      Server  `yaml:"server"`
	Storage     Storage `yaml:"storage"`
	GenAliasLen int     `yaml:"genAliasLen" env-default:"6"`
}

type Server struct {
	Host   string `yaml:"host" env-default:"localhost"`
	Port   string `yaml:"port" env-default:"8080"`
	Secret string `env:"SECRET" env-required:"true"`
}

type Storage struct {
	Path string `yaml:"path" env-required:"true"`
}

func New(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
