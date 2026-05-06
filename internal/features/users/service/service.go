package users_service

import (
	"context"

	"github.com/MulLoMaH/TODO_list.git/internal/core/domain"
)

type UsersService struct {
	userRepository UserRepository
}

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offser *int,
	) ([]domain.User, error)
}

func NewUserService(userRepository UserRepository) *UsersService {
	return &UsersService{
		userRepository: userRepository,
	}
}
