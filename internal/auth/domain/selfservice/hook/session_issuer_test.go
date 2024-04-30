// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package hook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"my.com/secrets/internal/auth/domain/external/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/randx"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/flow"
	"my.com/secrets/internal/auth/domain/selfservice/flow/registration"
	"my.com/secrets/internal/auth/domain/selfservice/hook"
	"my.com/secrets/internal/auth/domain/session"
	"my.com/secrets/internal/auth/domain/x"
)

func TestSessionIssuer(t *testing.T) {
	ctx := context.Background()
	conf, reg := external.NewFastRegistryWithMocks(t)
	conf.MustSet(ctx, config.ViperKeyPublicBaseURL, "http://localhost/")
	testhelpers.SetDefaultIdentitySchema(conf, "file://./stub/stub.schema.json")

	var r http.Request
	h := hook.NewSessionIssuer(reg)

	t.Run("method=sign-up", func(t *testing.T) {
		t.Run("flow=browser", func(t *testing.T) {
			w := httptest.NewRecorder()
			s := testhelpers.CreateSession(t, reg)
			f := &registration.Flow{Type: flow.TypeBrowser}

			require.NoError(t, h.ExecutePostRegistrationPostPersistHook(w, &r,
				f, &session.Session{ID: s.ID, Identity: s.Identity, Token: randx.MustString(12, randx.AlphaLowerNum)}))

			require.Empty(t, f.ContinueWithItems)

			got, err := reg.SessionPersister().GetSession(context.Background(), s.ID, session.ExpandNothing)
			require.NoError(t, err)
			assert.Equal(t, s.ID, got.ID)
			assert.True(t, got.AuthenticatedAt.After(time.Now().Add(-time.Minute)))

			assert.Contains(t, w.Header().Get("Set-Cookie"), config.DefaultSessionCookieName)
		})

		t.Run("flow=api", func(t *testing.T) {
			w := httptest.NewRecorder()

			i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
			s := &session.Session{
				ID:              x.NewUUID(),
				Identity:        i,
				Token:           randx.MustString(12, randx.AlphaLowerNum),
				LogoutToken:     randx.MustString(12, randx.AlphaLowerNum),
				AuthenticatedAt: time.Now().UTC(),
			}
			f := &registration.Flow{Type: flow.TypeAPI}

			require.NoError(t, reg.PrivilegedIdentityPool().CreateIdentity(context.Background(), i))
			require.NoError(t, reg.SessionPersister().UpsertSession(ctx, s))

			err := h.ExecutePostRegistrationPostPersistHook(w, &http.Request{Header: http.Header{"Accept": {"application/json"}}}, f, s)
			require.ErrorIs(t, err, registration.ErrHookAbortFlow, "%+v", err)
			require.Len(t, f.ContinueWithItems, 1)

			st := f.ContinueWithItems[0]
			require.IsType(t, &flow.ContinueWithSetOrySessionToken{}, st)
			assert.NotEmpty(t, st.(*flow.ContinueWithSetOrySessionToken).OrySessionToken)

			got, err := reg.SessionPersister().GetSession(context.Background(), s.ID, session.ExpandNothing)
			require.NoError(t, err)
			assert.Equal(t, s.ID.String(), got.ID.String())
			assert.True(t, got.AuthenticatedAt.After(time.Now().Add(-time.Minute)))

			assert.Empty(t, w.Header().Get("Set-Cookie"))
			body := w.Body.Bytes()
			assert.Equal(t, i.ID.String(), gjson.GetBytes(body, "identity.id").String())
			assert.Equal(t, s.ID.String(), gjson.GetBytes(body, "session.id").String())
			assert.Equal(t, got.Token, gjson.GetBytes(body, "session_token").String())
		})

		t.Run("flow=spa", func(t *testing.T) {
			w := httptest.NewRecorder()

			i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
			s := &session.Session{
				ID:              x.NewUUID(),
				Identity:        i,
				Token:           randx.MustString(12, randx.AlphaLowerNum),
				LogoutToken:     randx.MustString(12, randx.AlphaLowerNum),
				AuthenticatedAt: time.Now().UTC(),
			}
			f := &registration.Flow{Type: flow.TypeBrowser}

			require.NoError(t, reg.PrivilegedIdentityPool().CreateIdentity(context.Background(), i))
			require.NoError(t, reg.SessionPersister().UpsertSession(ctx, s))

			err := h.ExecutePostRegistrationPostPersistHook(w, &http.Request{Header: http.Header{"Accept": {"application/json"}}}, f, s)
			require.ErrorIs(t, err, registration.ErrHookAbortFlow, "%+v", err)
			require.Empty(t, f.ContinueWithItems)

			got, err := reg.SessionPersister().GetSession(context.Background(), s.ID, session.ExpandNothing)
			require.NoError(t, err)
			assert.Equal(t, s.ID.String(), got.ID.String())
			assert.True(t, got.AuthenticatedAt.After(time.Now().Add(-time.Minute)))

			assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
			body := w.Body.Bytes()
			assert.Equal(t, i.ID.String(), gjson.GetBytes(body, "identity.id").String())
			assert.Equal(t, s.ID.String(), gjson.GetBytes(body, "session.id").String())
			assert.Empty(t, gjson.GetBytes(body, "session_token").String())
		})
	})
}
