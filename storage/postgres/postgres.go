package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/untrik/url-shortener/internal/config"
	"github.com/untrik/url-shortener/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storage config.DB) (*Storage, error) {
	const op = "storage.postgres.New"
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		storage.Host, storage.Port, storage.User, storage.Password, storage.NameDB, storage.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s,%w", op, err)
	}
	stmt := (`
		CREATE TABLE IF NOT EXISTS url(
		id BIGSERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
		`)
	_, err = db.Exec(stmt)
	if err != nil {
		return nil, fmt.Errorf("%s,%w", op, err)
	}
	return &Storage{db: db}, nil
}
func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"
	stmt := (`INSERT INTO url(url,alias) VALUES ($1,$2) RETURNING id`)
	var id int64
	if err := s.db.QueryRow(stmt, urlToSave, alias).Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s,%w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s,%w", op, err)
	}
	return id, nil
}
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	var url string
	stmt := (`SELECT url FROM url WHERE alias = $1`)
	if err := s.db.QueryRow(stmt, alias).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}
func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.GetURL"
	stmt := (`DELETE FROM url WHERE alias = $1`)
	res, err := s.db.Exec(stmt, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	return nil
}
