package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log"
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

type DBStorage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
	StorageDelete
}
type DBStore struct{}

var db *sql.DB

func (s *DBStore) Get(ctx context.Context, key string) (value string, ok bool) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, postgre.selectURL, key)
	var originalURL string
	var deleted bool

	err := row.Scan(&originalURL, &deleted)

	if err != nil {
		logger.Errorln("SQL error", err)
		return "", false
	}
	if deleted {
		return "", true
	}

	return originalURL, true
}

func (s *DBStore) Set(ctx context.Context, key string, value string) error {
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	// сделал context.Background(), т.к. после добавления auth middleware появилась ошибка "context deadline exceeded"
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, postgre.insertURLWithUser, key, value, userID)

	if err != nil {
		// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = appErrors.NewHTTPError(409, "conflict")
		}
		logger.Errorln("SQL error", err)
	}

	return err
}

func (s *DBStore) List(ctx context.Context) (map[string]string, error) {
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	if userID == "" {
		return nil, appErrors.NewHTTPError(401, "user unauthorised")
	}
	logger.Warnln("userID", userID)

	rows, err := db.QueryContext(ctx, postgre.selectURLsByUserID, userID)
	if err != nil {
		logger.Errorln("SQL error:", err)
		return nil, err
	}

	if err := rows.Err(); err != nil {
		logger.Errorln("SQL error:", err)
		panic(err)
	}

	m := make(map[string]string)

	// пробегаем по всем записям
	for rows.Next() {
		var hash string
		var URL string
		err = rows.Scan(&hash, &URL)
		if err != nil {
			logger.Errorln("SQL error:", err)
			return nil, err
		}
		m[hash] = URL
	}

	return m, nil
}

func (s *DBStore) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	logger.Println("DB", reflect.TypeOf(db))

	if err := db.PingContext(ctx); err != nil {
		logger.Errorln("err", err)
		return err
	}
	return nil
}

func (s *DBStore) SetMultiple(ctx context.Context, items map[string]string) error {
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	// сделал context.Background(), т.к. после добавления auth middleware появилась ошибка "context deadline exceeded"
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	logger.Infoln("userID", userID)

	for k, v := range items {

		_, err := tx.ExecContext(ctx, postgre.insertURLIgnoreConflicts, k, v, userID)

		if err != nil {
			log.Println(postgre.insertURLIgnoreConflicts, k, v, userID)
			logger.Errorln("SQL error", err)
			// если ошибка, то откатываем
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *DBStore) Delete(ctx context.Context, hashes []interface{}) error {
	values := make([]interface{}, len(hashes)+1)
	valuesKeys := make([]string, len(hashes))

	values[0] = ctx.Value(ctxkeys.UserIDKey).(string)

	for i := 0; i < len(hashes); i++ {
		values[i+1] = hashes[i]
		valuesKeys[i] = "$" + strconv.Itoa(i+2)
	}

	r := strings.NewReplacer("$$$", strings.Join(valuesKeys, ","))
	query := r.Replace(postgre.setDeletedRows)

	// userID := ctx.Value(ctxkeys.UserIDKey).(string)

	// сделал context.Background(), т.к. после добавления auth middleware появилась ошибка "context deadline exceeded"
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, values...)

	logger.Errorln("query", query)
	logger.Errorln("values", values)
	if err != nil {
		logger.Errorln("SQL error", err)
		return err
	}

	return nil
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

	logger.Println("DB connection success", dsn)

	return &DBStore{}
}
