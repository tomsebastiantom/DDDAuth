// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package password_test

import (
	"context"
	"fmt"
	"testing"

	hash2 "my.com/secrets/internal/auth/domain/hash"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/password"
)

func TestCountActiveFirstFactorCredentials(t *testing.T) {
	ctx := context.Background()
	_, reg := external.NewFastRegistryWithMocks(t)
	strategy := password.NewStrategy(reg)

	h1, err := hash2.NewHasherBcrypt(reg).Generate(context.Background(), []byte("a password"))
	require.NoError(t, err)
	h2, err := reg.Hasher(ctx).Generate(context.Background(), []byte("a password"))
	require.NoError(t, err)

	for k, tc := range []struct {
		in       map[identity.CredentialsType]identity.Credentials
		expected int
	}{
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:   strategy.ID(),
				Config: []byte{},
			}},
			expected: 0,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:   strategy.ID(),
				Config: []byte(`{"hashed_password": "` + string(h1) + `"}`),
			}},
			expected: 0,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:        strategy.ID(),
				Identifiers: []string{""},
				Config:      []byte(`{"hashed_password": "` + string(h1) + `"}`),
			}},
			expected: 0,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:        strategy.ID(),
				Identifiers: []string{"foo"},
				Config:      []byte(`{"hashed_password": "` + string(h1) + `"}`),
			}},
			expected: 1,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:        strategy.ID(),
				Identifiers: []string{"foo"},
				Config:      []byte(`{"hashed_password": "` + string(h2) + `"}`),
			}},
			expected: 1,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:   strategy.ID(),
				Config: []byte(`{"hashed_password": "asdf"}`),
			}},
			expected: 0,
		},
		{
			in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
				Type:   strategy.ID(),
				Config: []byte(`{}`),
			}},
			expected: 0,
		},
		{
			in:       nil,
			expected: 0,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			actual, err := strategy.CountActiveFirstFactorCredentials(tc.in)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)

			actual, err = strategy.CountActiveMultiFactorCredentials(tc.in)
			assert.NoError(t, err)
			assert.Equal(t, 0, actual)
		})
	}
}
