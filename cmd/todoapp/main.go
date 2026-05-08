package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/MulLoMaH/TODO_list.git/internal/core/logger"
	core_pgx_pool "github.com/MulLoMaH/TODO_list.git/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/middleware"
	core_HTTP_server "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/server"
	user_postgres_repository "github.com/MulLoMaH/TODO_list.git/internal/features/users/repository/postgres"
	user_service "github.com/MulLoMaH/TODO_list.git/internal/features/users/service"
	user_transport_http "github.com/MulLoMaH/TODO_list.git/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postges connection pool")
	// pool, err := core_postgres_pool.NewConnectionPool(
	// 	ctx,
	// 	core_postgres_pool.NewConfigMust(),
	// )

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init  postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := user_postgres_repository.NewUserRepository(pool)
	usersService := user_service.NewUserService(usersRepository)

	usersTransportHTTP := user_transport_http.NewUserHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := core_HTTP_server.NewHTTPServer(
		logger,
		core_HTTP_server.NewConfigMust(),
		core_http_middleware.RequestID(),
		core_http_middleware.LoggerMiddleware(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_HTTP_server.NewAPIVersionRouter(core_HTTP_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
