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

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
}

func NewUserService(userRepository UserRepository) *UsersService {
	return &UsersService{
		userRepository: userRepository,
	}
}
