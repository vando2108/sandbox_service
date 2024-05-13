package repository_interface

import (
	"context"

	"github.com/vando2108/sandbox_service/internal/app/sandbox/model"
)

type UserRepository interface {
	FindOneByPublickey(context.Context, string) (*model.User, error)
	Insert(context.Context, *model.User) error
}
