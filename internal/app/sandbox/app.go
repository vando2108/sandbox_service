package sandbox

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vando2108/sandbox_service/internal/app/sandbox/handler"
	memRepo "github.com/vando2108/sandbox_service/internal/app/sandbox/repository/implement/memory"
	workqueue "github.com/vando2108/sandbox_service/internal/work_queue"
	"github.com/vando2108/sandbox_service/pb"
)

type SandboxServer struct {
	wq *workqueue.WorkQueue
}

func NewServer() *SandboxServer {
	return &SandboxServer{}
}

func (s *SandboxServer) Start(ctx context.Context) error {
	// init db
	userRepo := memRepo.NewMemUserRepository()
	nonceCache := memRepo.NewMemNonceCache()
	workQueue := workqueue.NewWorkQueue(10, 10*time.Second, 2)
	workQueue.Start()

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("failed to create listener: ", err)
	}
	log.Println("server is running at 127.0.0.1:9090")

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pb.RegisterAuthenticatorServer(grpcServer, handler.NewAuthenticatorHandler(&nonceCache, &userRepo))
	pb.RegisterSandboxServer(grpcServer, handler.NewSandboxHandler(&userRepo, workQueue))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln("failed to serve: ", err)
	}

	return nil
}

func (s *SandboxServer) Stop(ctx context.Context) error {
	s.wq.Stop()

	return nil
}
