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

	row := db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE short_url=$1", key)
	var original_url string
	err := row.Scan(&original_url)

	if err != nil {
		logger.Println("SQL error", err)
	}

	return original_url, err == nil
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
	// _, err = db.ExecContext(ctx, "CREATE INDEX video_id ON videos (video_id)")
}
