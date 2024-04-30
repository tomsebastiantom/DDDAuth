//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"my.com/secrets/config"
	"my.com/secrets/internal/others/application"
	"my.com/secrets/internal/others/domain/translation/entity"
	"my.com/secrets/internal/others/domain/translation/service"
	"my.com/secrets/internal/others/infrastructure/googleapi"
	"my.com/secrets/internal/others/infrastructure/repository"
	amqprpc "my.com/secrets/internal/others/interfaces/amqp_rpc"
	openapi "my.com/secrets/internal/others/interfaces/rest/v1/go"
	"my.com/secrets/pkg/httpserver"
	"my.com/secrets/pkg/logger"
	"my.com/secrets/pkg/postgres"
	"my.com/secrets/pkg/rabbitmq/rmq_rpc/server"
)

var deps = []interface{}{}

var providerSet wire.ProviderSet = wire.NewSet(
	postgres.NewOrGetSingleton,
	repository.New,
	googleapi.New,
	logger.New,
	amqprpc.NewRouter,
	server.New,
	httpserver.New,
	openapi.NewTranslator,
	openapi.NewRouter,

	application.NewWithDependencies,
	wire.Bind(new(entity.TranslationRepository), new(*repository.TranslationRepository)),
	wire.Bind(new(service.Translator), new(*googleapi.GoogleTranslator)),
)

var providerSetSystemTests wire.ProviderSet = wire.NewSet(
	postgres.NewOrGetSingleton,
	application.NewWithDependencies,
	logger.New,
	amqprpc.NewRouter,
	server.New,
	httpserver.New,
	openapi.NewTranslator,
	openapi.NewRouter,

	//wire.Bind(new(entity.TranslationRepository), new(*repository.TranslationRepository)),
	//wire.Bind(new(service.Translator), new(*googleapi.GoogleTranslator)),
)

func InitializeConfig() *config.Config {
	wire.Build(config.NewConfig)
	return &config.Config{}
}

func InitializePostgresConnection() *postgres.Postgres {
	wire.Build(providerSet, config.NewConfig)
	return &postgres.Postgres{}
}

func InitializeTranslationRepository() *repository.TranslationRepository {
	wire.Build(providerSet, config.NewConfig)
	return &repository.TranslationRepository{}
}

func InitializeTranslationWebAPI() *googleapi.GoogleTranslator {
	wire.Build(providerSet)
	return &googleapi.GoogleTranslator{}
}

func InitializeTranslationUseCase() *application.TranslationUseCase {
	wire.Build(providerSet, config.NewConfig)
	return &application.TranslationUseCase{}
}

func InitializeLogger() *logger.Logger {
	wire.Build(providerSet, config.NewConfig)
	return &logger.Logger{}
}

func InitializeNewRmqRpcServer() *server.Server {
	wire.Build(providerSet, config.NewConfig)
	return &server.Server{}
}

func InitializeNewRmqRpcServerForTesting(
	config *config.Config,
	translationRepository entity.TranslationRepository,
	translator service.Translator,
) *server.Server {
	wire.Build(providerSetSystemTests)
	return &server.Server{}
}

func InitializeNewHttpServerForTesting(
	config *config.Config,
	translationRepository entity.TranslationRepository,
	translator service.Translator,
) *httpserver.Server {
	wire.Build(providerSetSystemTests)
	return &httpserver.Server{}
}

func InitializeNewTranslator() *openapi.Translator {
	wire.Build(providerSet, config.NewConfig)
	return &openapi.Translator{}
}

func InitializeNewRouter() *gin.Engine {
	wire.Build(providerSet, config.NewConfig)
	return &gin.Engine{}
}

func InitializeNewHttpServer() *httpserver.Server {
	wire.Build(providerSet, config.NewConfig)
	return &httpserver.Server{}
}
