// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package lookup_test

import (
	"fmt"
	"testing"

	"my.com/secrets/internal/auth/domain/selfservice/strategy/lookup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/identity"
)

func TestCountActiveFirstFactorCredentials(t *testing.T) {
	_, reg := external.NewFastRegistryWithMocks(t)
	strategy := lookup.NewStrategy(reg)

	t.Run("first factor", func(t *testing.T) {
		actual, err := strategy.CountActiveFirstFactorCredentials(nil)
		require.NoError(t, err)
		assert.Equal(t, 0, actual)
	})

	t.Run("multi factor", func(t *testing.T) {
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
					Config: []byte(`{"recovery_codes": []}`),
				}},
				expected: 0,
			},
			{
				in: map[identity.CredentialsType]identity.Credentials{strategy.ID(): {
					Type:        strategy.ID(),
					Identifiers: []string{"foo"},
					Config:      []byte(`{"recovery_codes": [{}]}`),
				}},
				expected: 1,
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
				actual, err := strategy.CountActiveMultiFactorCredentials(tc.in)
				require.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			})
		}
	})
}
