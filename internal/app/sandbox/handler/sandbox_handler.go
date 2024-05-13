package handler

import (
	"context"

	"github.com/vando2108/sandbox_service/pb"
)

type SandboxHandler struct {
	pb.UnimplementedSandboxServer
}

func NewSandboxHandler() *SandboxHandler {
	return &SandboxHandler{}
}

func (h *SandboxHandler) CreateNewEnvironment(ctx context.Context, req *pb.CreateNewEnvironmentRequest) (*pb.RegisterResponse, error) {
	return nil, nil
}
