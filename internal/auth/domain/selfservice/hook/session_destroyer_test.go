// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package hook_test

import (
	"context"
	"net/http"
	"testing"

	"my.com/secrets/internal/auth/domain/external/testhelpers"

	"my.com/secrets/internal/auth/domain/corpx"
	"my.com/secrets/internal/auth/domain/ui/node"

	"github.com/go-faker/faker/v4"
	"github.com/gobuffalo/httptest"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/hook"
	"my.com/secrets/internal/auth/domain/session"
)

func init() {
	corpx.RegisterFakes()
}

func TestSessionDestroyer(t *testing.T) {
	ctx := context.Background()
	conf, reg := external.NewFastRegistryWithMocks(t)

	conf.MustSet(ctx, config.ViperKeyPublicBaseURL, "http://localhost/")
	testhelpers.SetDefaultIdentitySchema(conf, "file://./stub/stub.schema.json")

	h := hook.NewSessionDestroyer(reg)

	for _, tc := range []struct {
		name string
		hook func(*identity.Identity) error
	}{
		{
			name: "ExecuteLoginPostHook",
			hook: func(i *identity.Identity) error {
				return h.ExecuteLoginPostHook(
					httptest.NewRecorder(),
					new(http.Request),
					node.DefaultGroup,
					nil,
					&session.Session{Identity: i},
				)
			},
		},
		{
			name: "ExecutePostRecoveryHook",
			hook: func(i *identity.Identity) error {
				return h.ExecutePostRecoveryHook(
					httptest.NewRecorder(),
					new(http.Request),
					nil,
					&session.Session{Identity: i},
				)
			},
		},
		{
			name: "ExecuteSettingsPostPersistHook",
			hook: func(i *identity.Identity) error {
				return h.ExecuteSettingsPostPersistHook(
					httptest.NewRecorder(),
					new(http.Request),
					nil,
					i,
					&session.Session{Identity: i},
				)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var i identity.Identity
			require.NoError(t, faker.FakeData(&i))
			require.NoError(t, reg.PrivilegedIdentityPool().CreateIdentity(context.Background(), &i))

			sessions := make([]session.Session, 5)
			for k := range sessions {
				s := sessions[k] // keep this for pointers' sake ;)
				require.NoError(t, faker.FakeData(&s))
				s.IdentityID = uuid.Nil
				s.Identity = &i
				s.Active = true

				require.NoError(t, reg.SessionPersister().UpsertSession(context.Background(), &s))
				sessions[k] = s
			}

			for k := range sessions {
				sess, err := reg.SessionPersister().GetSession(context.Background(), sessions[k].ID, session.ExpandNothing)
				require.NoError(t, err)
				assert.True(t, sess.IsActive())
			}

			// Should revoke all the sessions.
			require.NoError(t, tc.hook(&i))

			for k := range sessions {
				sess, err := reg.SessionPersister().GetSession(context.Background(), sessions[k].ID, session.ExpandNothing)
				require.NoError(t, err)
				assert.False(t, sess.IsActive())
			}
		})
	}
}
