package handlers

type Storage interface {
	Get(key string) (value string, ok bool)
	Set(key string, value string) error
	List() map[string]string
}
