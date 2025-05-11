package store

import "context"

type StorageGet interface {
	Get(ctx context.Context, key string) (value string, ok bool)
}
type StorageList interface {
	List(ctx context.Context) (map[string]string, error)
}
type StoragePing interface {
	Ping(ctx context.Context) error
}
type StorageSet interface {
	Set(ctx context.Context, key string, value string) error
}
type StorageSetMultiple interface {
	SetMultiple(ctx context.Context, l map[string]string) error
}
type StorageDelete interface {
	Delete(ctx context.Context, l []string) error
}
type Storage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
	StorageDelete
}
