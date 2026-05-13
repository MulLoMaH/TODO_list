package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/MulLoMaH/TODO_list.git/internal/core/config"
	core_logger "github.com/MulLoMaH/TODO_list.git/internal/core/logger"
	core_pgx_pool "github.com/MulLoMaH/TODO_list.git/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/middleware"
	core_HTTP_server "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/server"

	statistics_postgres_repository "github.com/MulLoMaH/TODO_list.git/internal/features/statistics/repository/postgres"
	statistics_service "github.com/MulLoMaH/TODO_list.git/internal/features/statistics/service"
	statistics_transport_http "github.com/MulLoMaH/TODO_list.git/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/MulLoMaH/TODO_list.git/internal/features/tasks/repository/postgres"
	tasks_service "github.com/MulLoMaH/TODO_list.git/internal/features/tasks/service"
	tasks_transport_http "github.com/MulLoMaH/TODO_list.git/internal/features/tasks/transport/http"
	user_postgres_repository "github.com/MulLoMaH/TODO_list.git/internal/features/users/repository/postgres"
	user_service "github.com/MulLoMaH/TODO_list.git/internal/features/users/service"
	user_transport_http "github.com/MulLoMaH/TODO_list.git/internal/features/users/transport/http"
	web_fs_repository "github.com/MulLoMaH/TODO_list.git/internal/features/web/repository/html/file_system"
	web_service "github.com/MulLoMaH/TODO_list.git/internal/features/web/service"
	web_transport_http "github.com/MulLoMaH/TODO_list.git/internal/features/web/transport/http"
	"go.uber.org/zap"

	_ "github.com/MulLoMaH/TODO_list.git/docs"
)

// @title 			Golang Todo API
// @version 		1.0
// @description 	Todo Application REST-API scheme
// @host 			127.0.0.1:5050
// @basePath 		/api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

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

	logger.Debug("application time zone", zap.Any("zone", time.Local))
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

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTaskRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing feature", zap.String("feature", "web"))
	webRepoitory := web_fs_repository.NewWebRepoitory()
	webService := web_service.NewWebService(webRepoitory)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing HTTP server")
	httpConfig := core_HTTP_server.NewConfigMust()

	httpServer := core_HTTP_server.NewHTTPServer(
		logger,
		httpConfig,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.LoggerMiddleware(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_HTTP_server.NewAPIVersionRouter(core_HTTP_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
