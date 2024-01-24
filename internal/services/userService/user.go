package userservice

import (
	"errors"
	"fmt"
	"log/slog"
	"restapi/internal/config"
	"restapi/internal/domains"
	"restapi/internal/repository"
	service "restapi/internal/services"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	repo UserRepo
	cfg  *config.Config
	log  *slog.Logger
}

type UserRepo interface {
	CreateUser(login, password, email string) error
	GetByEmail(email string) (*domains.User, error)
	GetByLogin(login string) (*domains.User, error)
}

func NewUserService(repo UserRepo, cfg *config.Config, log *slog.Logger) *UserService {
	return &UserService{
		repo: repo,
		cfg:  cfg,
		log:  log,
	}
}

func (s *UserService) CreateUser(user *domains.User) (string, error) {
	fn := `services.UserService.CreateUser`
	err := s.repo.CreateUser(user.Login, user.Password, user.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return "", service.ErrUserAlredyExists
		}
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return generateToken(user, s.cfg.Server.Secret)
}

func (s *UserService) Login(login, password, email string) (string, error) {
	fn := `services.UserService.Login`
	var user *domains.User
	var err error
	if login != "" {
		user, err = s.repo.GetByLogin(login)
	} else {
		user, err = s.repo.GetByEmail(email)
	}

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", fn, service.ErrUserNotFound)
		}
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if user.Password != password {
		return "", fmt.Errorf("%s: %w", fn, service.ErrUserNotFound)
	}

	token, err := generateToken(user, s.cfg.Server.Secret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return token, nil
}

func generateToken(user *domains.User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"login": user.Login,
	})
	return token.SignedString([]byte(secret))
}
