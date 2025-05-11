package ctxkeys

type contextKey string

const (
	UserIDKey    contextKey = "userID"
	AuthTokenKey contextKey = "authToken"
	RequestIDKey contextKey = "requestID"
)
