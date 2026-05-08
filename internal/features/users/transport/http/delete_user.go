package users_transport_http

import (
	"net/http"

	core_logger "github.com/MulLoMaH/TODO_list.git/internal/core/logger"
	core_http_response "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/response"
	core_http_utils "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/utils"
)

func (h *UserHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	if err := h.UsersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
	}

	responseHandler.NoContentResponse()
}
