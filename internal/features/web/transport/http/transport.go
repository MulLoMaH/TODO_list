package web_transport_http

import (
	core_HTTP_server "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHTTPHandler(webService WebService) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: webService,
	}
}

func (h *WebHTTPHandler) Routes() []core_HTTP_server.Route {
	return []core_HTTP_server.Route{
		{
			Path:    "/",
			Handler: h.GetMainPage,
		},
	}
}
