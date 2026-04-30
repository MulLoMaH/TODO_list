package user_transport_http

type UserHTTPHandler struct {
	UsersService UsersService
}

type UsersService interface{}

func NewUserHTTPHandler(UsersService UsersService) *UserHTTPHandler {
	return &UserHTTPHandler{
		UsersService: UsersService,
	}
}
