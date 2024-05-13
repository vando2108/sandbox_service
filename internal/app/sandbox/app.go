package sandbox

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vando2108/sandbox_service/internal/app/sandbox/handler"
	"github.com/vando2108/sandbox_service/pb"
)

type SandboxServer struct {
	redisDB *redis.Client
}

func NewServer() *SandboxServer {
	return &SandboxServer{
		redisDB: redis.NewClient(&redis.Options{}),
	}
}

func (s *SandboxServer) Start(ctx context.Context) error {
	err := s.redisDB.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("failed to create listener: ", err)
	}
	log.Println("server is running at 127.0.0.1:9090")

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pb.RegisterAuthenticatorServer(grpcServer, handler.NewAuthenticatorHandler(s.redisDB))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln("failed to serve: ", err)
	}

	return nil
}
