// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oidc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/oidc"
)

func TestConfig(t *testing.T) {
	ctx := context.Background()
	conf, reg := external.NewFastRegistryWithMocks(t)

	var c map[string]interface{}
	require.NoError(t, json.NewDecoder(
		bytes.NewBufferString(`{"config":{"providers": [{"provider": "generic"}]}}`)).Decode(&c))
	conf.MustSet(ctx, config.ViperKeySelfServiceStrategyConfig+"."+string(identity.CredentialsTypeOIDC), c)

	s := oidc.NewStrategy(reg)
	collection, err := s.Config(ctx)
	require.NoError(t, err)

	require.Len(t, collection.Providers, 1)
	assert.Equal(t, "generic", collection.Providers[0].Provider)
}
