package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"reflect"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/aube/url-shortener/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DBStorage interface {
	Get(ctx context.Context, key string) (value string, ok bool)
	List(ctx context.Context) map[string]string
	Ping() error
	Set(ctx context.Context, key string, value string) error
	SetMultiple(ctx context.Context, l map[string]string) error
}

type DBStore struct{}

var db *sql.DB

func (s *DBStore) Get(ctx context.Context, key string) (value string, ok bool) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, postgre.selectURL, key)
	var originalURL string
	err := row.Scan(&originalURL)

	if err != nil {
		logger.Println("SQL error", err)
	}

	return originalURL, err == nil
}

func (s *DBStore) Set(ctx context.Context, key string, value string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, postgre.insertURL, key, value)

	if err != nil {
		// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
		logger.Println("SQL error", err)
	}

	return err
}

func (s *DBStore) List(ctx context.Context) map[string]string {
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

func (s *DBStore) SetMultiple(ctx context.Context, items map[string]string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for k, v := range items {
		_, err := tx.ExecContext(ctx, postgre.insertURLIgnoreConflicts, k, v)

		if err != nil {
			// если ошибка, то откатываем
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func NewDBStore(dsn string) DBStorage {
	var err error
	db, err = sql.Open("pgx", dsn)

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	logger.Println("DB connection success:", dsn)

	return &DBStore{}
}
