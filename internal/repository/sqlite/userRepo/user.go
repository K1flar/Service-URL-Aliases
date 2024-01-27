package userrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"restapi/internal/domains"
	"restapi/internal/repository"

	"github.com/mattn/go-sqlite3"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(login, password, email string) (*domains.User, error) {
	fn := `repository.sqlite.UserRepository.CreateUser`
	stmt, err := r.db.Prepare(`
		INSERT INTO users(login, password, email)
		VALUES (?, ?, ?);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec(login, password, email)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", fn, repository.ErrUserExists)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	user, err := r.GetByLogin(login)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domains.User, error) {
	fn := `repository.sqlite.UserRepository.GetByEmail`

	stmt, err := r.db.Prepare(`
		SELECT id, login, password, email FROM users
		WHERE email = ?;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var user domains.User
	err = stmt.QueryRow(email).Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &user, nil
}

func (r *UserRepository) GetByLogin(login string) (*domains.User, error) {
	fn := `repository.sqlite.UserRepository.GetByLogin`

	stmt, err := r.db.Prepare(`
		SELECT id, login, password, email FROM users
		WHERE login = ?;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var user domains.User
	err = stmt.QueryRow(login).Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &user, nil
}
