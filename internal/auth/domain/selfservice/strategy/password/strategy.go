// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package password

import (
	"context"
	"encoding/json"

	"my.com/secrets/internal/auth/domain/ui/node"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/ory/x/decoderx"

	"my.com/secrets/internal/auth/domain/continuity"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/hash"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/errorx"
	"my.com/secrets/internal/auth/domain/selfservice/flow/login"
	"my.com/secrets/internal/auth/domain/selfservice/flow/registration"
	"my.com/secrets/internal/auth/domain/selfservice/flow/settings"
	"my.com/secrets/internal/auth/domain/session"
	"my.com/secrets/internal/auth/domain/x"
)

var (
	_ login.Strategy                    = new(Strategy)
	_ registration.Strategy             = new(Strategy)
	_ identity.ActiveCredentialsCounter = new(Strategy)
)

type registrationStrategyDependencies interface {
	x.LoggingProvider
	x.WriterProvider
	x.CSRFTokenGeneratorProvider
	x.CSRFProvider

	config.Provider

	continuity.ManagementProvider

	errorx.ManagementProvider
	ValidationProvider
	hash.HashProvider

	registration.HandlerProvider
	registration.HooksProvider
	registration.ErrorHandlerProvider
	registration.HookExecutorProvider
	registration.FlowPersistenceProvider

	login.HooksProvider
	login.ErrorHandlerProvider
	login.HookExecutorProvider
	login.FlowPersistenceProvider
	login.HandlerProvider

	settings.FlowPersistenceProvider
	settings.HookExecutorProvider
	settings.HooksProvider
	settings.ErrorHandlerProvider

	identity.PrivilegedPoolProvider
	identity.ValidationProvider

	session.HandlerProvider
	session.ManagementProvider
}

type Strategy struct {
	d  registrationStrategyDependencies
	v  *validator.Validate
	hd *decoderx.HTTP
}

func NewStrategy(d any) *Strategy {
	return &Strategy{
		d:  d.(registrationStrategyDependencies),
		v:  validator.New(),
		hd: decoderx.NewHTTP(),
	}
}

func (s *Strategy) CountActiveFirstFactorCredentials(cc map[identity.CredentialsType]identity.Credentials) (count int, err error) {
	for _, c := range cc {
		if c.Type == s.ID() && len(c.Config) > 0 {
			var conf identity.CredentialsPassword
			if err = json.Unmarshal(c.Config, &conf); err != nil {
				return 0, errors.WithStack(err)
			}

			if len(c.Identifiers) > 0 && len(c.Identifiers[0]) > 0 &&
				(hash.IsBcryptHash([]byte(conf.HashedPassword)) || hash.IsArgon2idHash([]byte(conf.HashedPassword))) {
				count++
			}
		}
	}
	return
}

func (s *Strategy) CountActiveMultiFactorCredentials(cc map[identity.CredentialsType]identity.Credentials) (count int, err error) {
	return 0, nil
}

func (s *Strategy) ID() identity.CredentialsType {
	return identity.CredentialsTypePassword
}

func (s *Strategy) CompletedAuthenticationMethod(ctx context.Context, _ session.AuthenticationMethods) session.AuthenticationMethod {
	return session.AuthenticationMethod{
		Method: s.ID(),
		AAL:    identity.AuthenticatorAssuranceLevel1,
	}
}

func (s *Strategy) NodeGroup() node.UiNodeGroup {
	return node.PasswordGroup
}
