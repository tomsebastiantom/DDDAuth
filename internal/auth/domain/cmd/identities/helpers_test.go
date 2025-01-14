// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identities_test

import (
	"context"
	"testing"

	"github.com/ory/x/cmdx"

	"my.com/secrets/internal/auth/domain/identity"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/cmd/cliclient"
	"my.com/secrets/internal/auth/domain/driver"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/external/testhelpers"
)

func setup(t *testing.T, newCmd func() *cobra.Command) (driver.Registry, *cmdx.CommandExecuter) {
	conf, reg := external.NewFastRegistryWithMocks(t)
	_, admin := testhelpers.NewKratosServerWithCSRF(t, reg)
	testhelpers.SetDefaultIdentitySchema(conf, "file://./stubs/identity.schema.json")
	// setup command
	return reg, &cmdx.CommandExecuter{
		New: func() *cobra.Command {
			cmd := newCmd()
			cliclient.RegisterClientFlags(cmd.Flags())
			cmdx.RegisterFormatFlags(cmd.Flags())
			return cmd
		},
		PersistentArgs: []string{"--" + cliclient.FlagEndpoint, admin.URL, "--" + cmdx.FlagFormat, string(cmdx.FormatJSON)},
	}
}

func makeIdentities(t *testing.T, reg driver.Registry, n int) (is []*identity.Identity, ids []string) {
	for j := 0; j < n; j++ {
		i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
		i.MetadataPublic = []byte(`{"foo":"bar"}`)
		require.NoError(t, reg.Persister().CreateIdentity(context.Background(), i))
		is = append(is, i)
		ids = append(ids, i.ID.String())
	}
	return
}
