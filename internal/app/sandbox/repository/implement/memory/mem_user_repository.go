package repository_implement

import (
	"context"
	"fmt"

	"github.com/vando2108/sandbox_service/internal/app/sandbox/model"
	repoInterface "github.com/vando2108/sandbox_service/internal/app/sandbox/repository/interface"
)

type MemUserRepository struct {
	db map[string]*model.User
}

func NewMemUserRepository() repoInterface.UserRepository {
	return &MemUserRepository{
		db: make(map[string]*model.User),
	}
}

func (r *MemUserRepository) FindOneByPublickey(ctx context.Context, publickey string) (*model.User, error) {
	user, exists := r.db[publickey]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (r *MemUserRepository) Insert(ctx context.Context, user *model.User) error {
	if _, exists := r.db[user.Publickey]; exists {
		return fmt.Errorf("user already exists")
	}
	user.ID = len(r.db) + 1
	r.db[user.Publickey] = user

	return nil
}
