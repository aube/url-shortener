// Package grpcserver implements the gRPC service for URL shortening operations.
// It provides both the service implementation and request validation middleware.
package grpcserver

import (
	"context"
	"strings"

	"github.com/aube/url-shortener/internal/api/grpc/proto"
	"github.com/aube/url-shortener/internal/app/usecases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationInterceptor is a gRPC unary interceptor that validates incoming requests.
// It checks for required fields and proper formatting before passing requests to handlers.
// Returns gRPC errors for invalid requests.
func ValidationInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	switch r := req.(type) {
	case *proto.SaveURLRequest:
		if strings.TrimSpace(r.OriginalUrl) == "" {
			return nil, status.Error(codes.InvalidArgument, "original_url cannot be empty")
		}
		if strings.TrimSpace(r.BaseUrl) == "" {
			return nil, status.Error(codes.InvalidArgument, "base_url cannot be empty")
		}
		if !isValidURL(r.OriginalUrl) {
			return nil, status.Error(codes.InvalidArgument, "original_url must be a valid URL")
		}
	case *proto.GetURLRequest:
		if strings.TrimSpace(r.Id) == "" {
			return nil, status.Error(codes.InvalidArgument, "id cannot be empty")
		}
	}

	return handler(ctx, req)
}

// isValidURL performs basic URL validation.
// Checks if the URL starts with http:// or https://.
// Returns true if valid, false otherwise.
func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// URLShortenerServer implements the gRPC service
type URLShortenerServer struct {
	proto.UnimplementedUrlShortenerServer
}

// NewURLShortenerServer creates a new instance of the URLShortenerServer.
// Returns a pointer to the initialized server.
func NewURLShortenerServer() *URLShortenerServer {
	return &URLShortenerServer{}
}

// SaveURL implements the gRPC method to shorten a URL.
// Takes a SaveURLRequest containing the original URL and base URL.
// Returns a SaveURLResponse with the shortened URL or an error.
func (s *URLShortenerServer) SaveURL(ctx context.Context, req *proto.SaveURLRequest) (*proto.SaveURLResponse, error) {
	shortURL, err := usecases.SaveURL(ctx, []byte(req.OriginalUrl), req.BaseUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save URL: %v", err)
	}

	return &proto.SaveURLResponse{
		ShortUrl: shortURL,
	}, nil
}

// GetURL implements the gRPC method to retrieve the original URL.
// Takes a GetURLRequest containing the shortened URL ID.
// Returns a GetURLResponse with the original URL or an error if not found.
func (s *URLShortenerServer) GetURL(ctx context.Context, req *proto.GetURLRequest) (*proto.GetURLResponse, error) {
	originalURL, err := usecases.GetURL(ctx, req.Id)
	if err != nil {
		if err.Error() == "not found" {
			return nil, status.Error(codes.NotFound, "URL not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get URL: %v", err)
	}

	return &proto.GetURLResponse{
		OriginalUrl: originalURL,
	}, nil
}
