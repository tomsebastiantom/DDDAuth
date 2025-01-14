// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"my.com/secrets/config"
	"my.com/secrets/internal/others/application"
	"my.com/secrets/internal/others/domain/translation/entity"
	"my.com/secrets/internal/others/domain/translation/service"
	"my.com/secrets/internal/others/infrastructure/googleapi"
	"my.com/secrets/internal/others/infrastructure/repository"
	"my.com/secrets/internal/others/interfaces/amqp_rpc"
	"my.com/secrets/internal/others/interfaces/rest/v1/go"
	"my.com/secrets/pkg/httpserver"
	"my.com/secrets/pkg/logger"
	"my.com/secrets/pkg/postgres"
	"my.com/secrets/pkg/rabbitmq/rmq_rpc/server"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeConfig() *config.Config {
	configConfig := config.NewConfig()
	return configConfig
}

func InitializePostgresConnection() *postgres.Postgres {
	configConfig := config.NewConfig()
	postgresPostgres := postgres.NewOrGetSingleton(configConfig)
	return postgresPostgres
}

func InitializeTranslationRepository() *repository.TranslationRepository {
	configConfig := config.NewConfig()
	postgresPostgres := postgres.NewOrGetSingleton(configConfig)
	translationRepository := repository.New(postgresPostgres)
	return translationRepository
}

func InitializeTranslationWebAPI() *googleapi.GoogleTranslator {
	googleTranslator := googleapi.New()
	return googleTranslator
}

func InitializeTranslationUseCase() *application.TranslationUseCase {
	configConfig := config.NewConfig()
	postgresPostgres := postgres.NewOrGetSingleton(configConfig)
	translationRepository := repository.New(postgresPostgres)
	googleTranslator := googleapi.New()
	translationUseCase := application.NewWithDependencies(translationRepository, googleTranslator)
	return translationUseCase
}

func InitializeLogger() *logger.Logger {
	configConfig := config.NewConfig()
	loggerLogger := logger.New(configConfig)
	return loggerLogger
}

func InitializeNewRmqRpcServer() *server.Server {
	configConfig := config.NewConfig()
	loggerLogger := logger.New(configConfig)
	postgresPostgres := postgres.NewOrGetSingleton(configConfig)
	translationRepository := repository.New(postgresPostgres)
	googleTranslator := googleapi.New()
	translationUseCase := application.NewWithDependencies(translationRepository, googleTranslator)
	v := amqprpc.NewRouter(translationUseCase)
	serverServer := server.New(configConfig, loggerLogger, v)
	return serverServer
}

func InitializeNewRmqRpcServerForTesting(config2 *config.Config, translationRepository entity.TranslationRepository, translator service.Translator) *server.Server {
	loggerLogger := logger.New(config2)
	translationUseCase := application.NewWithDependencies(translationRepository, translator)
	v := amqprpc.NewRouter(translationUseCase)
	serverServer := server.New(config2, loggerLogger, v)
	return serverServer
}

func InitializeNewHttpServerForTesting(config2 *config.Config, translationRepository entity.TranslationRepository, translator service.Translator) *httpserver.Server {
	translationUseCase := application.NewWithDependencies(translationRepository, translator)
	loggerLogger := logger.New(config2)
	openapiTranslator := openapi.NewTranslator(translationUseCase, loggerLogger)
	engine := openapi.NewRouter(openapiTranslator)
	httpserverServer := httpserver.New(config2, engine)
	return httpserverServer
}

func InitializeNewTranslator() *openapi.Translator {
	configConfig := config.NewConfig()
	postgresPostgres := postgres.NewOrGetSingleton(configConfig)
	translationRepository := repository.New(postgresPostgres)
	googleTranslator := googleapi.New()
	translationUseCase := application.NewWithDependencies(translationRepository, googleTranslator)
	loggerLogger := logger.New(configConfig)
	translator := openapi.NewTranslator(translationUseCase, loggerLogger)
	return translator
}

func InitializeNewRouter() *gin.Engine {
	configConfig := config.NewConfig()
	postgresPostgres := postgres.NewOrGetSingleton(configConfig)
	translationRepository := repository.New(postgresPostgres)
	googleTranslator := googleapi.New()
	translationUseCase := application.NewWithDependencies(translationRepository, googleTranslator)
	loggerLogger := logger.New(configConfig)
	translator := openapi.NewTranslator(translationUseCase, loggerLogger)
	engine := openapi.NewRouter(translator)
	return engine
}

func InitializeNewHttpServer() *httpserver.Server {
	configConfig := config.NewConfig()
	postgresPostgres := postgres.NewOrGetSingleton(configConfig)
	translationRepository := repository.New(postgresPostgres)
	googleTranslator := googleapi.New()
	translationUseCase := application.NewWithDependencies(translationRepository, googleTranslator)
	loggerLogger := logger.New(configConfig)
	translator := openapi.NewTranslator(translationUseCase, loggerLogger)
	engine := openapi.NewRouter(translator)
	httpserverServer := httpserver.New(configConfig, engine)
	return httpserverServer
}

// wire.go:

var deps = []interface{}{}

var providerSet wire.ProviderSet = wire.NewSet(postgres.NewOrGetSingleton, repository.New, googleapi.New, logger.New, amqprpc.NewRouter, server.New, httpserver.New, openapi.NewTranslator, openapi.NewRouter, application.NewWithDependencies, wire.Bind(new(entity.TranslationRepository), new(*repository.TranslationRepository)), wire.Bind(new(service.Translator), new(*googleapi.GoogleTranslator)))

var providerSetSystemTests wire.ProviderSet = wire.NewSet(postgres.NewOrGetSingleton, application.NewWithDependencies, logger.New, amqprpc.NewRouter, server.New, httpserver.New, openapi.NewTranslator, openapi.NewRouter)
