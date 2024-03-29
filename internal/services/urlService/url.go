package urlservice

import (
	"errors"
	"fmt"
	"log/slog"
	"restapi/internal/config"
	"restapi/internal/domains"
	"restapi/internal/lib/random"
	"restapi/internal/repository"
	service "restapi/internal/services"
)

type URLService struct {
	repo URLRepo
	cfg  *config.Config
	log  *slog.Logger
}

type URLRepo interface {
	SaveURL(url, alias string, userID uint32) error
	GetURL(alias string) (*domains.URL, error)
	DeleteURL(alias string) error
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
		if errors.Is(err, repository.ErrURLExists) {
			return "", fmt.Errorf("%s: %w", fn, service.ErrURLExists)
		}
		s.log.Info(err.Error())
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return alias, nil
}

func (s *URLService) GetURL(alias string) (string, error) {
	fn := `services.url.GetURL`
	url, err := s.repo.GetURL(alias)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			return "", fmt.Errorf("%s: %w", fn, service.ErrURLNotFound)
		}
		s.log.Info(err.Error())
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return url.URL, nil
}

func (s *URLService) DeleteURL(alias string, userID uint32) error {
	fn := `services.url.DeleteURL`
	url, err := s.repo.GetURL(alias)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			return fmt.Errorf("%s: %w", fn, service.ErrURLNotFound)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}
	if url.UserID != userID {
		return fmt.Errorf("%s: %w", fn, service.ErrURLForbiddenToDelete)
	}

	err = s.repo.DeleteURL(url.Alias)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			return fmt.Errorf("%s: %w", fn, repository.ErrURLNotFound)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
