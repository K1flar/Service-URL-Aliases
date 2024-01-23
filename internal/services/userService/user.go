package userservice

import (
	"fmt"
	"log/slog"
	"restapi/internal/config"
	"restapi/internal/domains"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	repo UserRepo
	cfg  *config.Config
	log  *slog.Logger
}

type UserRepo interface {
	CreateUser(login, password, email string) error
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
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return generateToken(user, s.cfg.Server.Secret)
}

func generateToken(user *domains.User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"login": user.Login,
	})
	return token.SignedString([]byte(secret))
}
