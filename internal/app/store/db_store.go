package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/app/ctxkeys"
	"github.com/aube/url-shortener/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DBStore struct{}

var db *sql.DB

func (s *DBStore) Get(ctx context.Context, key string) (value string, ok bool) {
	log := logger.WithContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, postgre.selectURL, key)
	var originalURL string
	var deleted bool

	err := row.Scan(&originalURL, &deleted)

	if err != nil {
		log.Error("Get", "err", err)
		return "", false
	}
	if deleted {
		return "", true
	}

	return originalURL, true
}

func (s *DBStore) Set(ctx context.Context, key string, value string) error {
	log := logger.WithContext(ctx)

	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, postgre.insertURLWithUser, key, value, userID)

	if err != nil {
		// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = appErrors.NewHTTPError(409, "conflict")
		}
		log.Error("Set", "err", err)
	}

	return err
}

func (s *DBStore) List(ctx context.Context) (map[string]string, error) {
	log := logger.WithContext(ctx)

	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	if userID == "" {
		return nil, appErrors.NewHTTPError(401, "user unauthorised")
	}
	log.Warn("List", "userID", userID)

	rows, err := db.QueryContext(ctx, postgre.selectURLsByUserID, userID)
	if err != nil {
		log.Error("List", "err", err)
		return nil, err
	}

	if err := rows.Err(); err != nil {
		log.Error("List", "rows.Err", err)
		panic(err)
	}

	m := make(map[string]string)

	// пробегаем по всем записям
	for rows.Next() {
		var hash string
		var URL string
		err = rows.Scan(&hash, &URL)
		if err != nil {
			log.Error("List", "err", err)
			return nil, err
		}
		m[hash] = URL
	}

	return m, nil
}

func (s *DBStore) Ping(ctx context.Context) error {
	log := logger.WithContext(ctx)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	log.Debug("Ping", "db", reflect.TypeOf(db))

	if err := db.PingContext(ctx); err != nil {
		log.Error("Ping", "err", err)
		return err
	}
	return nil
}

func (s *DBStore) SetMultiple(ctx context.Context, items map[string]string) error {
	log := logger.WithContext(ctx)

	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	log.Info("SetMultiple", "userID", userID)

	for k, v := range items {

		_, err := tx.ExecContext(ctx, postgre.insertURLIgnoreConflicts, k, v, userID)

		if err != nil {
			log.Error("SetMultiple", "err", err)
			// если ошибка, то откатываем транзакцию
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *DBStore) Delete(ctx context.Context, hashes []string) error {
	log := logger.WithContext(ctx)

	values := make([]interface{}, len(hashes)+1) // array of query values
	valuesKeys := make([]string, len(hashes))    // "$2,$3...$n"

	// first value in query sets for: user_id=$1
	values[0] = ctx.Value(ctxkeys.UserIDKey).(string)

	for i := 0; i < len(hashes); i++ {
		values[i+1] = hashes[i]
		valuesKeys[i] = "$" + strconv.Itoa(i+2)
	}

	r := strings.NewReplacer("$$$", strings.Join(valuesKeys, ","))
	query := r.Replace(postgre.setDeletedRows)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, values...)

	if err != nil {
		log.Error("Delete", "err", err)
		log.Error("Delete", "query", query, "values", values)
		return err
	}

	return nil
}

func NewDBStore(dsn string) Storage {
	log := logger.Get()

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

	log.Debug("NewDBStore", "dsn", dsn)

	return &DBStore{}
}
