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
	"github.com/aube/url-shortener/internal/app/workerpool"
	"github.com/aube/url-shortener/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// DBStore is a PostgreSQL implementation of the Storage interface.
type DBStore struct {
	dispatcher *workerpool.WorkDispatcher
}

var db *sql.DB

// Get retrieves a URL by its shortened key from the database.
// Returns the URL and true if found (even if deleted), empty string and false otherwise.
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

// Set stores a new URL mapping in the database.
// Returns an error if the operation fails, including a conflict error if the key exists.
func (s *DBStore) Set(ctx context.Context, key string, value string) error {
	log := logger.WithContext(ctx)

	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, postgre.insertURLWithUser, key, value, userID)

	if err != nil {
		// Check if error is a integrity constraint violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = appErrors.NewHTTPError(409, "conflict")
		}
		log.Error("Set", "err", err)
	}

	return err
}

// List returns all URL mappings for the current user from the database.
// Returns an unauthorized error if no user ID is present in context.
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

	// Iterate through all records
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

// Ping checks if the database connection is alive.
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

// SetMultiple stores multiple URL mappings in a single transaction.
// If any operation fails, the entire transaction is rolled back.
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
			// Rollback transaction on error
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Delete marks one or more URLs as deleted in the database.
// Only URLs belonging to the current user are affected.
func (s *DBStore) Delete(ctx context.Context, hashes []string) error {
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	for _, hash := range hashes {
		s.dispatcher.AddWork(ctx, hash, userID)
	}

	// s.delMultiple(ctx, hashes, userID)

	return nil
}

func (s *DBStore) delByRow(ctx context.Context, hash string, userID string) error {
	log := logger.WithContext(ctx)

	_, err := db.Exec(postgre.setDeleteOnceRow, userID, hash)

	if err != nil {
		log.Error("Delete", "query", postgre.setDeleteOnceRow, "userID", userID, "hash", hash)
		log.Error("Delete", "err", err)
		return err
	}
	return nil
}

func (s *DBStore) delMultiple(ctx context.Context, hashes []string, userID string) error {
	log := logger.WithContext(ctx)
	values := make([]any, len(hashes)+1)      // array of query values
	valuesKeys := make([]string, len(hashes)) // "$2,$3...$n"

	// first value in query sets for: user_id=$1
	values[0] = userID

	for i := range len(hashes) {
		values[i+1] = hashes[i]
		valuesKeys[i] = "$" + strconv.Itoa(i+2)
	}

	r := strings.NewReplacer("$$$", strings.Join(valuesKeys, ","))
	query := r.Replace(postgre.setDeletedRows)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, values...)

	if err != nil {
		log.Error("Delete", "query", query, "values", values)
		log.Error("Delete", "err", err)
		return err
	}
	return nil
}

// NewDBStore creates and initializes a new PostgreSQL storage instance.
// It establishes a database connection, sets connection pool parameters,
// and runs any pending migrations using Goose.
func NewDBStore(dsn string) Storage {
	log := logger.Get()

	var err error
	db, err = sql.Open("pgx", dsn)

	if err != nil {
		panic(err)
	}

	// Set connection pool parameters
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Configure Goose migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	log.Debug("NewDBStore", "dsn", dsn)

	store := &DBStore{}
	store.initWorkerPool()

	return store
}

func (s *DBStore) initWorkerPool() {

	s.dispatcher = workerpool.New(3, s.delByRow)
	// defer dispatcher.Close()

	// for _, id := range orders {
	// 	dispatcher.AddWork(id)
	// }
}
