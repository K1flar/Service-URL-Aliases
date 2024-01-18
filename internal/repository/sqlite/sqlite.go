package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"restapi/internal/repository"

	"github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func New(path string) (*Repository, error) {
	fn := `repository.sqlite.New`

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT NOT NULL,
			password TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS url (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			alias TEXT NOT NULL UNIQUE,
			user_id INTEGER REFERENCES users(id) ON DELETE SET NULL
		);
	`)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	tx.Commit()
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) SaveURL(url, alias string, userID uint32) error {
	fn := `repository.sqlite.SaveURL`
	stmt, err := r.db.Prepare(`
		INSERT INTO url(url, alias, user_id)
		VALUES (?, ?, ?);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec(url, alias, userID)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", fn, repository.ErrURLExists)
		}

		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *Repository) GetURL(alias string) (string, error) {
	fn := `repository.sqlite.GetURL`
	stmt, err := r.db.Prepare(`
		SELECT url FROM url
		WHERE alias = ?
	`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", fn, repository.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return url, nil
}
