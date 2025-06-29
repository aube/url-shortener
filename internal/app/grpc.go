package app

import (
	"log"
	"net"

	grpcserver "github.com/aube/url-shortener/internal/api/grpc"
	"github.com/aube/url-shortener/internal/api/grpc/proto"
	"google.golang.org/grpc"
)

func grpcServerStart() *grpc.Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create server with validation interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcserver.ValidationInterceptor),
	)

	proto.RegisterUrlShortenerServer(grpcServer, grpcserver.NewURLShortenerServer())

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// defer grpcServer.Stop()

	return grpcServer
}
