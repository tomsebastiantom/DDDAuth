// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identities_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/tidwall/gjson"

	"my.com/secrets/internal/auth/domain/cmd/identities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlcon"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/x"
)

func TestDeleteCmd(t *testing.T) {
	reg, cmd := setup(t, identities.NewDeleteIdentityCmd)

	t.Run("case=deletes successfully", func(t *testing.T) {
		// create identity to delete
		i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
		require.NoError(t, reg.Persister().CreateIdentity(context.Background(), i))

		stdOut := cmd.ExecNoErr(t, i.ID.String())

		// expect ID and no error
		assert.Equal(t, i.ID.String(), gjson.Parse(stdOut).String())

		// expect identity to be deleted
		_, err := reg.Persister().GetIdentity(context.Background(), i.ID, identity.ExpandNothing)
		assert.True(t, errors.Is(err, sqlcon.ErrNoRows))
	})

	t.Run("case=deletes three identities", func(t *testing.T) {
		is, ids := makeIdentities(t, reg, 3)

		stdOut := cmd.ExecNoErr(t, ids...)

		assert.Equal(t, `["`+strings.Join(ids, "\",\"")+"\"]\n", stdOut)

		for _, i := range is {
			_, err := reg.Persister().GetIdentity(context.Background(), i.ID, identity.ExpandNothing)
			assert.Error(t, err)
		}
	})

	t.Run("case=fails with unknown ID", func(t *testing.T) {
		stdErr := cmd.ExecExpectedErr(t, x.NewUUID().String())

		assert.Contains(t, stdErr, "Unable to locate the resource", stdErr)
	})
}
