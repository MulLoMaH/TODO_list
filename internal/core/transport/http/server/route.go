package core_HTTP_server

import (
	"net/http"

	core_http_middleware "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

// func NewRoute(
// 	method string,
// 	path string,
// 	handler http.HandlerFunc,
// ) *Route {
// 	return &Route{
// 		Method:  method,
// 		Path:    path,
// 		Handler: handler,
// 	}
// }

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.CheinMiddleware(
		r.Handler,
		r.Middleware...,
	)

}
