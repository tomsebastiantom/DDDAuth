// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package webauthn_test

import (
	"context"
	"fmt"
	"testing"

	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/webauthn"
	"my.com/secrets/internal/auth/domain/session"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/identity"
)

func TestCompletedAuthenticationMethod(t *testing.T) {
	conf, reg := external.NewFastRegistryWithMocks(t)
	strategy := webauthn.NewStrategy(reg)

	assert.Equal(t, session.AuthenticationMethod{
		Method: strategy.ID(),
		AAL:    identity.AuthenticatorAssuranceLevel2,
	}, strategy.CompletedAuthenticationMethod(context.Background(), session.AuthenticationMethods{}))

	conf.MustSet(ctx, config.ViperKeyWebAuthnPasswordless, true)
	assert.Equal(t, session.AuthenticationMethod{
		Method: strategy.ID(),
		AAL:    identity.AuthenticatorAssuranceLevel1,
	}, strategy.CompletedAuthenticationMethod(context.Background(), session.AuthenticationMethods{}))
}

func TestCountActiveFirstFactorCredentials(t *testing.T) {
	_, reg := external.NewFastRegistryWithMocks(t)
	strategy := webauthn.NewStrategy(reg)

	for k, tc := range []struct {
		in            map[identity.CredentialsType]identity.Credentials
		expectedFirst int
		expectedMulti int
	}{
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:   strategy.ID(),
				Config: []byte{},
			}},
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:   strategy.ID(),
				Config: []byte(`{"credentials": []}`),
			}},
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:        strategy.ID(),
				Identifiers: []string{"foo"},
				Config:      []byte(`{"credentials": [{}]}`),
			}},
			expectedMulti: 1,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:        strategy.ID(),
				Identifiers: []string{"foo"},
				Config:      []byte(`{"credentials": [{"is_passwordless": true}]}`),
			}},
			expectedFirst: 1,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:        strategy.ID(),
				Identifiers: []string{"foo"},
				Config:      []byte(`{"credentials": [{"is_passwordless": true}, {"is_passwordless": true}]}`),
			}},
			expectedFirst: 2,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:        strategy.ID(),
				Identifiers: []string{"foo"},
				Config:      []byte(`{"credentials": [{"is_passwordless": true}, {"is_passwordless": false}]}`),
			}},
			expectedFirst: 1,
			expectedMulti: 1,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:   strategy.ID(),
				Config: []byte(`{}`),
			}},
		},
		{
			in: nil,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			cc := map[identity.CredentialsType]identity.Credentials{}
			for _, c := range tc.in {
				cc[c.Type] = c
			}

			actual, err := strategy.CountActiveFirstFactorCredentials(cc)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFirst, actual)

			actual, err = strategy.CountActiveMultiFactorCredentials(cc)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedMulti, actual)
		})
	}
}
