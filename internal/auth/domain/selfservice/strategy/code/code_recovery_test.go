// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package code_test

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"my.com/secrets/internal/auth/domain/external"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/code"
	"my.com/secrets/internal/auth/domain/x"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/urlx"

	"my.com/secrets/internal/auth/domain/selfservice/flow"
	"my.com/secrets/internal/auth/domain/selfservice/flow/recovery"
)

func TestRecoveryCode(t *testing.T) {
	conf, _ := external.NewFastRegistryWithMocks(t)

	newCode := func(expiresIn time.Duration, f *recovery.Flow) *code.RecoveryCode {
		return &code.RecoveryCode{
			ID:        x.NewUUID(),
			FlowID:    f.ID,
			ExpiresAt: time.Now().Add(expiresIn),
		}
	}

	req := &http.Request{URL: urlx.ParseOrPanic("https://www.ory.sh/")}
	t.Run("method=Validate", func(t *testing.T) {
		t.Parallel()

		t.Run("case=returns error if flow is expired", func(t *testing.T) {
			f, err := recovery.NewFlow(conf, -time.Hour, "", req, nil, flow.TypeBrowser)
			require.NoError(t, err)

			c := newCode(-time.Hour, f)
			expected := new(flow.ExpiredError)
			require.ErrorAs(t, c.Validate(), &expected)
		})
		t.Run("case=returns no error if flow is not expired", func(t *testing.T) {
			f, err := recovery.NewFlow(conf, time.Hour, "", req, nil, flow.TypeBrowser)
			require.NoError(t, err)

			c := newCode(time.Hour, f)
			require.NoError(t, c.Validate())
		})

		t.Run("case=returns error if flow has been used", func(t *testing.T) {
			f, err := recovery.NewFlow(conf, -time.Hour, "", req, nil, flow.TypeBrowser)
			require.NoError(t, err)

			c := newCode(time.Hour, f)
			c.UsedAt = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
			require.ErrorIs(t, c.Validate(), code.ErrCodeAlreadyUsed)
		})

		t.Run("case=returns no error if flow has not been used", func(t *testing.T) {
			f, err := recovery.NewFlow(conf, -time.Hour, "", req, nil, flow.TypeBrowser)
			require.NoError(t, err)

			c := newCode(time.Hour, f)
			c.UsedAt = sql.NullTime{
				Valid: false,
			}
			require.NoError(t, c.Validate())
		})

		t.Run("case=returns error if flow is nil", func(t *testing.T) {
			var c *code.RecoveryCode
			require.ErrorIs(t, c.Validate(), code.ErrCodeNotFound)
		})
	})
}

func TestRecoveryCodeType(t *testing.T) {
	assert.Equal(t, 1, int(code.RecoveryCodeTypeAdmin))
	assert.Equal(t, 2, int(code.RecoveryCodeTypeSelfService))
}
