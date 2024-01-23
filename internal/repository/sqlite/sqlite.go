package sqlite

import (
	"database/sql"
	"fmt"
	urlrepo "restapi/internal/repository/sqlite/urlRepo"
	userrepo "restapi/internal/repository/sqlite/userRepo"
)

type Repository struct {
	*urlrepo.URLRepository
	*userrepo.UserRepository
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
		urlrepo.NewURLRepository(db),
		userrepo.NewUserRepository(db),
	}, nil
}
