// Package ctxkeys provides types to use as keys in a context.Context
package ctxkeys

type contextKey string

// UserIDKey is used to store and retrieve user IDs from a Context
// AuthTokenKey is used to store and retrieve authentication tokens from a Context
// RequestIDKey is used to store and retrieve request IDs from a Context
const (
	UserIDKey contextKey = "userID"

	AuthTokenKey contextKey = "authToken"

	RequestIDKey contextKey = "requestID"
)
