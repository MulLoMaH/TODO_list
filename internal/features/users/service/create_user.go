package user_service

import (
	"context"
	"fmt"

	"github.com/MulLoMaH/TODO_list.git/internal/core/domain"
)

func (s *UserService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("user domain: %w", err)
	}

	return user, nil
}
