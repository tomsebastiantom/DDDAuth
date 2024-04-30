// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package link

import (
	"github.com/ory/x/decoderx"
	"my.com/secrets/internal/auth/domain/courier"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/schema"
	"my.com/secrets/internal/auth/domain/selfservice/errorx"
	"my.com/secrets/internal/auth/domain/selfservice/flow/recovery"
	"my.com/secrets/internal/auth/domain/selfservice/flow/settings"
	"my.com/secrets/internal/auth/domain/selfservice/flow/verification"
	"my.com/secrets/internal/auth/domain/session"
	"my.com/secrets/internal/auth/domain/ui/container"
	"my.com/secrets/internal/auth/domain/ui/node"
	"my.com/secrets/internal/auth/domain/x"
)

var (
	_ recovery.Strategy      = new(Strategy)
	_ recovery.AdminHandler  = new(Strategy)
	_ recovery.PublicHandler = new(Strategy)
)

var (
	_ verification.Strategy      = new(Strategy)
	_ verification.AdminHandler  = new(Strategy)
	_ verification.PublicHandler = new(Strategy)
)

type (
	// FlowMethod contains the configuration for this selfservice strategy.
	FlowMethod struct {
		*container.Container
	}

	strategyDependencies interface {
		x.CSRFProvider
		x.CSRFTokenGeneratorProvider
		x.WriterProvider
		x.LoggingProvider
		x.TracingProvider

		config.Provider

		session.HandlerProvider
		session.ManagementProvider
		settings.HandlerProvider
		settings.FlowPersistenceProvider

		identity.ValidationProvider
		identity.ManagementProvider
		identity.PoolProvider
		identity.PrivilegedPoolProvider

		courier.Provider

		errorx.ManagementProvider

		recovery.ErrorHandlerProvider
		recovery.FlowPersistenceProvider
		recovery.StrategyProvider
		recovery.HookExecutorProvider

		verification.ErrorHandlerProvider
		verification.FlowPersistenceProvider
		verification.StrategyProvider
		verification.HookExecutorProvider
		verification.HandlerProvider

		RecoveryTokenPersistenceProvider
		VerificationTokenPersistenceProvider
		SenderProvider

		schema.IdentityTraitsProvider
	}

	Strategy struct {
		d  strategyDependencies
		dx *decoderx.HTTP
	}
)

func NewStrategy(d any) *Strategy {
	return &Strategy{d: d.(strategyDependencies), dx: decoderx.NewHTTP()}
}

func (s *Strategy) NodeGroup() node.UiNodeGroup {
	return node.LinkGroup
}
