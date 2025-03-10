package store

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/aube/url-shortener/internal/logger"
)

type StorageDB interface {
	Storage
	Ping() error
}

type DBStore struct{}

var db *sql.DB

func (s *DBStore) Get(key string) (value string, ok bool) {
	// value, ok = dbData.s[key]
	// logger.Infoln("Get key:", key, value)
	return value, ok
}

func (s *DBStore) Set(key string, value string) error {
	// if key == "" || value == "" {
	// 	return fmt.Errorf("invalid input")
	// }

	// logger.Infoln("Set key:", key, value)
	// dbData.s[key] = value

	return nil
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

func NewDBStore(dsn string) StorageDB {
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dsn, `videouser`, `videopass`, `videodb`)

	var err error
	db, err = sql.Open("pgx", ps)

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
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS videos (
        "video_id" TEXT,
        "title" TEXT,
        "publish_time" TEXT,
        "tags" TEXT,
        "views" INTEGER
      )`)
	if err != nil {
		logger.Println("createDB error:", err)
	}
	// _, err = db.ExecContext(ctx, "CREATE INDEX video_id ON videos (video_id)")
}
