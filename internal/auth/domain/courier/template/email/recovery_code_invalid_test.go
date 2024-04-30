// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package email_test

import (
	"context"
	"testing"

	"my.com/secrets/internal/auth/domain/courier/template"
	"my.com/secrets/internal/auth/domain/courier/template/email"
	"my.com/secrets/internal/auth/domain/courier/template/testhelpers"
	"my.com/secrets/internal/auth/domain/external"
)

func TestRecoveryCodeInvalid(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("test=with courier templates directory", func(t *testing.T) {
		_, reg := external.NewFastRegistryWithMocks(t)
		tpl := email.NewRecoveryCodeInvalid(reg, &email.RecoveryCodeInvalidModel{})

		testhelpers.TestRendered(t, ctx, tpl)
	})

	t.Run("case=test remote resources", func(t *testing.T) {
		testhelpers.TestRemoteTemplates(t, "../courier/builtin/templates/recovery_code/invalid", template.TypeRecoveryCodeInvalid)
	})
}