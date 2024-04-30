package amqprpc

import (
	"sync"

	"my.com/secrets/internal/others/application"
	"my.com/secrets/pkg/rabbitmq/rmq_rpc/server"
)

var hdlOnce sync.Once
var amqpRpcRouter map[string]server.CallHandler

// NewRouter -.
func NewRouter(t *application.TranslationUseCase) map[string]server.CallHandler {

	hdlOnce.Do(func() {
		amqpRpcRouter = make(map[string]server.CallHandler)
		{
			newTranslationRoutes(amqpRpcRouter, t)
		}
	})

	return amqpRpcRouter
}
