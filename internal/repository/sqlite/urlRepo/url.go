package urlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"restapi/internal/domains"
	"restapi/internal/repository"

	"github.com/mattn/go-sqlite3"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) SaveURL(url, alias string, userID uint32) error {
	fn := `repository.sqlite.URLRepository.SaveURL`
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

func (r *URLRepository) GetURL(alias string) (*domains.URL, error) {
	fn := `repository.sqlite.URLRepository.GetURL`
	stmt, err := r.db.Prepare(`
		SELECT id, url, alias, user_id FROM url
		WHERE alias = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var url domains.URL
	err = stmt.QueryRow(alias).Scan(&url.ID, &url.URL, &url.Alias, &url.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, repository.ErrURLNotFound)
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &url, nil
}

func (r *URLRepository) DeleteURL(alias string) error {
	fn := `repository.sqlite.URLRepository.DeleteURL`

	stmt, err := r.db.Prepare(`
		DELETE FROM url
		WHERE alias = ?;
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	if rowsAffected == 0 {
		return repository.ErrURLNotFound
	}

	return nil
}
