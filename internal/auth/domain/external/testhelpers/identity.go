// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/driver"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/session"
)

func CreateSession(t *testing.T, reg driver.Registry) *session.Session {
	req := NewTestHTTPRequest(t, "GET", "/sessions/whoami", nil)
	i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
	require.NoError(t, reg.PrivilegedIdentityPool().CreateIdentity(req.Context(), i))
	sess, err := session.NewActiveSession(req, i, reg.Config(), time.Now().UTC(), identity.CredentialsTypePassword, identity.AuthenticatorAssuranceLevel1)
	require.NoError(t, err)
	require.NoError(t, reg.SessionPersister().UpsertSession(req.Context(), sess))
	return sess
}
