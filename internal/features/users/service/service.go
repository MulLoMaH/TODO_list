package user_service

import (
	"context"

	"github.com/MulLoMaH/TODO_list.git/internal/core/domain"
)

type UserService struct {
	userRepository UserRepository
}

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}
