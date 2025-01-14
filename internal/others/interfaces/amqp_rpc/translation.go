package amqprpc

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"

	"my.com/secrets/internal/others/application"
	"my.com/secrets/internal/others/domain/translation/entity"
	"my.com/secrets/pkg/rabbitmq/rmq_rpc/server"
)

type translationRoutes struct {
	translationUseCase *application.TranslationUseCase
}

func newTranslationRoutes(routes map[string]server.CallHandler, t *application.TranslationUseCase) {
	r := &translationRoutes{t}
	{
		routes["getHistory"] = r.getHistory()
	}
}

type historyResponse struct {
	History []entity.Translation `json:"history"`
}

func (r *translationRoutes) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		translations, err := r.translationUseCase.History(context.Background())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - translationRoutes - getHistory - r.translationUseCase.History: %w", err)
		}

		response := historyResponse{translations}

		return response, nil
	}
}
