package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	repoInterface "github.com/vando2108/sandbox_service/internal/app/sandbox/repository/interface"
	workqueue "github.com/vando2108/sandbox_service/internal/work_queue"
	"github.com/vando2108/sandbox_service/pb"
	"github.com/vando2108/sandbox_service/utils"
)

type SandboxHandler struct {
	userRepo  *repoInterface.UserRepository
	workQueue *workqueue.WorkQueue
	pb.UnimplementedSandboxServer
}

func NewSandboxHandler(userRepo *repoInterface.UserRepository, workQueue *workqueue.WorkQueue) *SandboxHandler {
	return &SandboxHandler{
		userRepo:  userRepo,
		workQueue: workQueue,
	}
}

func (h *SandboxHandler) CreateNewEnvironment(ctx context.Context, req *pb.CreateNewEnvironmentRequest) (*pb.CreateNewEnvironmentResponse, error) {
	if utils.ValidatePublickey(req.Publickey) != nil {
		return &pb.CreateNewEnvironmentResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID,
			ErrorMessage: ErrorMessages[pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID],
		}, nil
	}

	if _, err := (*h.userRepo).FindOneByPublickey(ctx, req.Publickey); err != nil && err.Error() == "user not found" {
		return &pb.CreateNewEnvironmentResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_USER_NOT_EXISTED,
			ErrorMessage: ErrorMessages[pb.ErrorCode_USER_NOT_EXISTED],
		}, nil
	}

	if verifySignature(req) != nil {
		return &pb.CreateNewEnvironmentResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_INVALID_SIGNATURE,
			ErrorMessage: ErrorMessages[pb.ErrorCode_INVALID_SIGNATURE],
		}, nil
	}

	// Insert new entry to environment table
	// Insert new entry to environment_status table with status = pending
	// enqueue one request to message_queue
	task := func() error {
		return createNewEnvironment(int(time.Now().UnixNano()))
	}
	h.workQueue.AddTask(task)

	return &pb.CreateNewEnvironmentResponse{
		Success: true,
	}, nil
}

func verifySignature(req *pb.CreateNewEnvironmentRequest) error {
	// for test
	return nil

	messageBytes, err := json.Marshal(req.Requirements)
	if err != nil {
		return err
	}

	signature, err := hex.DecodeString(req.Signature[2:])
	if err != nil {
		return err
	}

	messageHash := sha256.Sum256(messageBytes)

	publickey, err := crypto.SigToPub(messageHash[:], signature)
	if err != nil {
		return nil
	}

	if utils.PublickeyToString(publickey) != req.Publickey {
		return fmt.Errorf("Signature is not valid")
	}

	return nil
}

func createNewEnvironment(envID int) error {
	// pickup one request from message queue and create new env
	log.Println("Create new env: ", envID)
	return nil
}
