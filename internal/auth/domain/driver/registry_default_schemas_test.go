// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/ory/x/urlx"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/schema"
)

func TestRegistryDefault_IdentityTraitsSchemas(t *testing.T) {
	ctx := context.Background()

	conf, reg := external.NewFastRegistryWithMocks(t)
	defaultSchema := schema.Schema{
		ID:     "default",
		URL:    urlx.ParseOrPanic("file://default.schema.json"),
		RawURL: "file://default.schema.json",
	}
	altSchema := schema.Schema{
		ID:     "alt",
		URL:    urlx.ParseOrPanic("file://other.schema.json"),
		RawURL: "file://other.schema.json",
	}

	conf.MustSet(ctx, config.ViperKeyIdentitySchemas, []config.Schema{
		{ID: altSchema.ID, URL: altSchema.RawURL},
		{ID: defaultSchema.ID, URL: defaultSchema.RawURL},
	})

	ss, err := reg.IdentityTraitsSchemas(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, len(ss))
	assert.Contains(t, ss, defaultSchema)
	assert.Contains(t, ss, altSchema)
}
