// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package hook_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/selfservice/flow"
	"my.com/secrets/internal/auth/domain/selfservice/flow/registration"
	"my.com/secrets/internal/auth/domain/selfservice/flow/verification"
	"my.com/secrets/internal/auth/domain/selfservice/hook"
)

func TestExecutePostRegistrationPostPersistHook(t *testing.T) {
	t.Run("case=no continue with items returns 200 OK", func(t *testing.T) {
		_, reg := external.NewVeryFastRegistryWithoutDB(t)
		h := hook.NewShowVerificationUIHook(reg)
		browserRequest := httptest.NewRequest("GET", "/", nil)
		f := &registration.Flow{}
		rec := httptest.NewRecorder()
		require.NoError(t, h.ExecutePostRegistrationPostPersistHook(rec, browserRequest, f, nil))
		require.Equal(t, 200, rec.Code)
	})

	t.Run("case=not a browser request returns 200 OK", func(t *testing.T) {
		_, reg := external.NewVeryFastRegistryWithoutDB(t)
		h := hook.NewShowVerificationUIHook(reg)
		browserRequest := httptest.NewRequest("GET", "/", nil)
		browserRequest.Header.Add("Accept", "application/json")
		f := &registration.Flow{}
		rec := httptest.NewRecorder()
		require.NoError(t, h.ExecutePostRegistrationPostPersistHook(rec, browserRequest, f, nil))
		require.Equal(t, 200, rec.Code)
	})

	t.Run("case=verification ui in continue with item returns redirect", func(t *testing.T) {
		conf, reg := external.NewVeryFastRegistryWithoutDB(t)
		conf.Set(context.Background(), config.ViperKeySelfServiceVerificationUI, "/verification")
		h := hook.NewShowVerificationUIHook(reg)
		browserRequest := httptest.NewRequest("GET", "/", nil)
		vf := &verification.Flow{
			ID: uuid.Must(uuid.NewV4()),
		}
		rf := &registration.Flow{}
		rf.ContinueWithItems = []flow.ContinueWith{
			flow.NewContinueWithVerificationUI(vf, "some@ory.sh", ""),
		}
		rec := httptest.NewRecorder()
		require.NoError(t, h.ExecutePostRegistrationPostPersistHook(rec, browserRequest, rf, nil))
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "/verification?flow="+vf.ID.String(), rf.ReturnToVerification)
	})

	t.Run("case=no verification ui in continue with item returns 200 OK", func(t *testing.T) {
		conf, reg := external.NewVeryFastRegistryWithoutDB(t)
		conf.Set(context.Background(), config.ViperKeySelfServiceVerificationUI, "/verification")
		h := hook.NewShowVerificationUIHook(reg)
		browserRequest := httptest.NewRequest("GET", "/", nil)
		rf := &registration.Flow{}
		rf.ContinueWithItems = []flow.ContinueWith{
			flow.NewContinueWithSetToken("token"),
		}
		rec := httptest.NewRecorder()
		require.NoError(t, h.ExecutePostRegistrationPostPersistHook(rec, browserRequest, rf, nil))
		assert.Equal(t, 200, rec.Code)
	})
}
