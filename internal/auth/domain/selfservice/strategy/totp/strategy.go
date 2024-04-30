// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package totp

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pquerna/otp"

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
	"my.com/secrets/internal/auth/domain/ui/node"
	"my.com/secrets/internal/auth/domain/x"
)

var (
	_ login.Strategy                    = new(Strategy)
	_ settings.Strategy                 = new(Strategy)
	_ identity.ActiveCredentialsCounter = new(Strategy)
)

type totpStrategyDependencies interface {
	x.LoggingProvider
	x.WriterProvider
	x.CSRFTokenGeneratorProvider
	x.CSRFProvider

	config.Provider

	continuity.ManagementProvider

	errorx.ManagementProvider
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
	session.PersistenceProvider
}

type Strategy struct {
	d  totpStrategyDependencies
	hd *decoderx.HTTP
}

func NewStrategy(d any) *Strategy {
	return &Strategy{
		d:  d.(totpStrategyDependencies),
		hd: decoderx.NewHTTP(),
	}
}

func (s *Strategy) CountActiveFirstFactorCredentials(cc map[identity.CredentialsType]identity.Credentials) (count int, err error) {
	return 0, nil
}

func (s *Strategy) CountActiveMultiFactorCredentials(cc map[identity.CredentialsType]identity.Credentials) (count int, err error) {
	for _, c := range cc {
		if c.Type == s.ID() && len(c.Config) > 0 {
			var conf identity.CredentialsTOTPConfig
			if err = json.Unmarshal(c.Config, &conf); err != nil {
				return 0, errors.WithStack(err)
			}

			_, err := otp.NewKeyFromURL(conf.TOTPURL)
			if len(c.Identifiers) > 0 && len(c.Identifiers[0]) > 0 && len(conf.TOTPURL) > 0 && err == nil {
				count++
			}
		}
	}
	return
}

func (s *Strategy) ID() identity.CredentialsType {
	return identity.CredentialsTypeTOTP
}

func (s *Strategy) NodeGroup() node.UiNodeGroup {
	return node.TOTPGroup
}

func (s *Strategy) CompletedAuthenticationMethod(ctx context.Context, _ session.AuthenticationMethods) session.AuthenticationMethod {
	return session.AuthenticationMethod{
		Method: s.ID(),
		AAL:    identity.AuthenticatorAssuranceLevel2,
	}
}
