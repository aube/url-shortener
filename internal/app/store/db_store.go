package store

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/aube/url-shortener/internal/logger"
)

type DBStore struct{}

var db *sql.DB

func (s *DBStore) Get(key string) (value string, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, "SELECT original_url as originalURL FROM urls WHERE short_url=$1", key)
	var originalURL string
	err := row.Scan(&originalURL)

	if err != nil {
		logger.Println("SQL error", err)
	}

	return originalURL, err == nil
}

func (s *DBStore) Set(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO urls (short_url, original_url) VALUES($1, $2)", key, value)

	if err != nil {
		logger.Println("SQL error", err)
	}

	return err
}

func (s *DBStore) List() map[string]string {
	m := make(map[string]string)
	return m
}

func (s *DBStore) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	logger.Println("DB", reflect.TypeOf(db))

	if err := db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) SetMultiple(items map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for k, v := range items {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING", k, v)

		if err != nil {
			// если ошибка, то откатываем
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func NewDBStore(dsn string) Storage {
	var err error
	db, err = sql.Open("pgx", dsn)

	if err != nil {
		panic(err)
	}
	// defer db.Close()

	createDB(db)

	logger.Println("DB connection success:", dsn)

	return &DBStore{}
}

func createDB(db *sql.DB) {
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS urls (
        id serial PRIMARY KEY,
        short_url CHAR(10) UNIQUE,
        original_url TEXT
      )`)
	if err != nil {
		logger.Println("createDB error:", err)
	}
}
