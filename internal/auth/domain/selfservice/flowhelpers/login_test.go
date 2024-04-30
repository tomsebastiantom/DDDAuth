// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flowhelpers_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/external/testhelpers"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/flow/login"
	"my.com/secrets/internal/auth/domain/selfservice/flowhelpers"
	"my.com/secrets/internal/auth/domain/session"
)

func TestGuessForcedLoginIdentifier(t *testing.T) {
	conf, reg := external.NewFastRegistryWithMocks(t)
	testhelpers.SetDefaultIdentitySchema(conf, "file://./stub/login.schema.json")

	i := identity.NewIdentity("")
	ic := identity.Credentials{
		Type:        identity.CredentialsTypePassword,
		Identifiers: []string{"foobar"},
	}
	i.Credentials[identity.CredentialsTypePassword] = ic
	require.NoError(t, reg.IdentityManager().Create(context.Background(), i))

	req := httptest.NewRequest("GET", "/sessions/whoami", nil)

	sess, err := session.NewActiveSession(req, i, conf, time.Now(), identity.CredentialsTypePassword, identity.AuthenticatorAssuranceLevel1)
	require.NoError(t, err)
	reg.SessionPersister().UpsertSession(context.Background(), sess)

	r := httptest.NewRequest("GET", "/login", nil)
	r.Header.Set("Authorization", "Bearer "+sess.Token)

	var f login.Flow
	f.Refresh = true

	identifier, id, creds := flowhelpers.GuessForcedLoginIdentifier(r, reg, &f, identity.CredentialsTypePassword)
	assert.Equal(t, "foobar", identifier)
	assert.EqualValues(t, ic.Type, creds.Type)
	assert.EqualValues(t, ic.Identifiers, creds.Identifiers)
	assert.EqualValues(t, id.ID, id.ID)
}
