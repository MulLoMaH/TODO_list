package users_transport_http

import (
	"context"
	"net/http"

	"github.com/MulLoMaH/TODO_list.git/internal/core/domain"
	core_HTTP_server "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/server"
)

type UserHTTPHandler struct {
	UsersService UsersService
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
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
		patch domain.UserPatch,
	) (domain.User, error)
}

func NewUserHTTPHandler(UsersService UsersService) *UserHTTPHandler {
	return &UserHTTPHandler{
		UsersService: UsersService,
	}
}

func (h *UserHTTPHandler) Routes() []core_HTTP_server.Route {
	return []core_HTTP_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
