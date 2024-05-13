package handler

import "github.com/vando2108/sandbox_service/pb"

var ErrorMessages = map[pb.ErrorCode]string{
	pb.ErrorCode_NONCE_MISMATCH:          "Provided nonce does not match",
	pb.ErrorCode_INTERNAL_SERVER_ERROR:   "Internal server error",
	pb.ErrorCode_NONCE_NOT_EXSIST:        "Nonce does not exist or has expired. Please register again.",
	pb.ErrorCode_PUBLIC_KEY_IS_NOT_VALID: "Publickey is not valid",
	pb.ErrorCode_USER_EXISTSED:           "User existed",
	pb.ErrorCode_INVALID_SIGNATURE:       "Invalid signature",
	pb.ErrorCode_USER_NOT_EXISTED:        "User not existed",
}
