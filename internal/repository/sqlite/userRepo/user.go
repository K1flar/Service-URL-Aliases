package userrepo

import (
	"database/sql"
	"fmt"
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

func (r *UserRepository) CreateUser(login, password, email string) error {
	fn := `repository.sqlite.UserRepository.CreateUser`
	stmt, err := r.db.Prepare(`
		INSERT INTO users(login, password, email)
		VALUES (?, ?, ?);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec(login, password, email)
	if err != nil {
		if _, ok := err.(sqlite3.Error); ok && err == sqlite3.ErrConstraintUnique {
			return repository.ErrUserExists
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
