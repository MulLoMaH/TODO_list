package users_transport_http

import (
	"net/http"

	core_logger "github.com/MulLoMaH/TODO_list.git/internal/core/logger"
	core_http_request "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/request"
	core_http_response "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

func (h *UserHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	user, err := h.UsersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)

		return
	}

	response := GetUserResponse(userDtoFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
