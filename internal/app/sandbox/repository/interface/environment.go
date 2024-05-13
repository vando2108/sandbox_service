package repository_interface

import (
	"context"

	"github.com/vando2108/sandbox_service/internal/app/sandbox/model"
)

type EnvironmentRepository interface {
	Insert(context.Context, *model.Environment) error
}
