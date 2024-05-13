package handler

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/redis/go-redis/v9"

	"github.com/vando2108/sandbox_service/internal/app/sandbox/model"
	memRepo "github.com/vando2108/sandbox_service/internal/app/sandbox/repository/implement/memory"
	redisRepo "github.com/vando2108/sandbox_service/internal/app/sandbox/repository/implement/redis"
	repoInterface "github.com/vando2108/sandbox_service/internal/app/sandbox/repository/interface"
	"github.com/vando2108/sandbox_service/pb"
	"github.com/vando2108/sandbox_service/utils"
)

var errorMessages = map[pb.ErrorCode]string{
	pb.ErrorCode_NONCE_MISMATCH:          "Provided nonce does not match",
	pb.ErrorCode_INTERNAL_SERVER_ERROR:   "Internal server error",
	pb.ErrorCode_NONCE_NOT_EXSIST:        "Nonce does not exist or has expired. Please register again.",
	pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID: "Publickey is not valid",
	pb.ErrorCode_USER_EXISTSED:           "User existed",
}

type AuthenticatorHandler struct {
	nonceCache repoInterface.NonceCacheRepository
	userRepo   repoInterface.UserRepository
	pb.UnimplementedAuthenticatorServer
}

func NewAuthenticatorHandler(redisClient *redis.Client) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		nonceCache: redisRepo.NewRedisNonceCache(redisClient),
		userRepo:   memRepo.NewMemUserRepository(),
	}
}

func (h *AuthenticatorHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := utils.ValidatePublickey(req.Publickey); err != nil {
		// Suppose that ValidatePublickey method cover all cases
		return &pb.RegisterResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID,
			ErrorMessage: errorMessages[pb.ErrorCode(pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID)],
		}, nil
	}

	if _, err := h.userRepo.FindOneByPublickey(ctx, req.Publickey); err == nil || err.Error() != "user not found" {
		return &pb.RegisterResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_USER_EXISTSED,
			ErrorMessage: errorMessages[pb.ErrorCode_USER_EXISTSED],
		}, nil
	}

	// Using server timestamps as a random factor is easily predictable and introduces a potential security vulnerability.
	// While it might be acceptable for minimal projects
	nonce := strconv.Itoa(int(time.Now().UnixNano()))
	log.Println(req.Publickey, nonce)

	publickeyBytes, err := hexutil.Decode(req.Publickey)
	if err != nil {
		log.Println(err)
		return &pb.RegisterResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID,
			ErrorMessage: errorMessages[pb.ErrorCode(pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID)],
		}, nil
	}

	x, y := new(big.Int).SetBytes(publickeyBytes[:32]), new(big.Int).SetBytes(publickeyBytes[32:])
	publicKey := &ecdsa.PublicKey{Curve: ecies.DefaultCurve, X: x, Y: y}
	eciesPublickey := ecies.ImportECDSAPublic(publicKey)

	encryptedNonce, err := ecies.Encrypt(rand.Reader, eciesPublickey, []byte(nonce), nil, nil)
	if err != nil {
		return &pb.RegisterResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_INTERNAL_SERVER_ERROR,
			ErrorMessage: errorMessages[pb.ErrorCode_INTERNAL_SERVER_ERROR],
		}, nil
	}

	h.nonceCache.SetNonce(ctx, req.Publickey, nonce, time.Duration(5*time.Minute))

	return &pb.RegisterResponse{
		Success:     true,
		HashedNonce: hexutil.Encode(encryptedNonce),
	}, nil
}

func (h *AuthenticatorHandler) NonceConfirm(ctx context.Context, req *pb.NonceConfirmRequest) (*pb.NonceConfirmResponse, error) {
	if err := utils.ValidatePublickey(req.Publickey); err != nil {
		return &pb.NonceConfirmResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID,
			ErrorMessage: errorMessages[pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID],
		}, nil
	}

	nonce, err := h.nonceCache.GetNonce(ctx, req.Publickey)
	if err != nil {
		if err == redis.Nil {
			return &pb.NonceConfirmResponse{
				Success:      false,
				ErrorCode:    pb.ErrorCode_NONCE_NOT_EXSIST,
				ErrorMessage: errorMessages[pb.ErrorCode_NONCE_NOT_EXSIST],
			}, nil
		} else {
			return &pb.NonceConfirmResponse{
				Success:      false,
				ErrorCode:    pb.ErrorCode_INTERNAL_SERVER_ERROR,
				ErrorMessage: errorMessages[pb.ErrorCode_INTERNAL_SERVER_ERROR],
			}, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	if nonce != req.Nonce {
		return &pb.NonceConfirmResponse{
			Success:      false,
			ErrorCode:    pb.ErrorCode_NONCE_MISMATCH,
			ErrorMessage: errorMessages[pb.ErrorCode_NONCE_MISMATCH],
		}, nil
	} else {
		user := &model.User{
			Publickey: req.Publickey,
		}

		if err := h.userRepo.Insert(ctx, user); err != nil {
			return &pb.NonceConfirmResponse{
				Success:      false,
				ErrorCode:    pb.ErrorCode_USER_EXISTSED,
				ErrorMessage: errorMessages[pb.ErrorCode_USER_EXISTSED],
			}, nil
		}

		return &pb.NonceConfirmResponse{
			Success: true,
		}, nil
	}
}
