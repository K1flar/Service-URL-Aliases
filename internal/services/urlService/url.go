package urlservice

import (
	"fmt"
	"log/slog"
	"restapi/internal/config"
	"restapi/internal/lib/random"
)

type URLService struct {
	repo URLRepo
	cfg  *config.Config
	log  *slog.Logger
}

type URLRepo interface {
	SaveURL(url, alias string, userID uint32) error
	GetURL(alias string) (string, error)
}

func NewURLService(repo URLRepo, cfg *config.Config, log *slog.Logger) *URLService {
	return &URLService{
		repo: repo,
		cfg:  cfg,
		log:  log,
	}
}

func (s *URLService) SaveURL(url, alias string, userID uint32) (string, error) {
	fn := `services.url.SaveURL`
	if alias == "" {
		alias = random.GenerateRandomString(s.cfg.GenAliasLen)
	}
	err := s.repo.SaveURL(url, alias, userID)
	if err != nil {
		s.log.Info(err.Error())
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return alias, nil
}

func (s *URLService) GetURL(alias string) (string, error) {
	fn := `services.url.GetURL`
	url, err := s.repo.GetURL(alias)
	if err != nil {
		s.log.Info(err.Error())
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return url, nil
}

//TODO: DeleteURLa
