package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"my.com/secrets/internal"
	"my.com/secrets/pkg/httpserver"
	"my.com/secrets/pkg/logger"
	"my.com/secrets/pkg/rabbitmq/rmq_rpc/server"
)

func main() {
	log := internal.InitializeLogger()

	rmqServer, httpServer := startServers()
	err := waitForSignals(log, httpServer, rmqServer)
	shutdown(err, httpServer, log, rmqServer)
}

func startServers() (*server.Server, *httpserver.Server) {
	rmqServer := internal.InitializeNewRmqRpcServer()
	httpServer := internal.InitializeNewHttpServer()
	return rmqServer, httpServer
}

func waitForSignals(log *logger.Logger, httpServer *httpserver.Server, rmqServer *server.Server) error {
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		log.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}
	return err
}

func shutdown(err error, httpServer *httpserver.Server, log *logger.Logger, rmqServer *server.Server) {
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}
}
