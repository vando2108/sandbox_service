package repository_interface

import (
	"context"
	"time"
)

type NonceCacheRepository interface {
	SetNonce(context.Context, string, string, time.Duration) error
	GetNonce(context.Context, string) (string, error)
}
