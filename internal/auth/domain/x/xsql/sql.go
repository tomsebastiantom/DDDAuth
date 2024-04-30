// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package xsql

import (
	"context"
	"testing"

	"github.com/gobuffalo/pop/v6"

	"my.com/secrets/internal/auth/domain/selfservice/errorx"
	"my.com/secrets/internal/auth/domain/selfservice/sessiontokenexchange"

	"my.com/secrets/internal/auth/domain/continuity"
	"my.com/secrets/internal/auth/domain/courier"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/flow/login"
	"my.com/secrets/internal/auth/domain/selfservice/flow/recovery"
	"my.com/secrets/internal/auth/domain/selfservice/flow/registration"
	"my.com/secrets/internal/auth/domain/selfservice/flow/settings"
	"my.com/secrets/internal/auth/domain/selfservice/flow/verification"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/code"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/link"
	"my.com/secrets/internal/auth/domain/session"
)

func CleanSQL(t testing.TB, c *pop.Connection) {
	ctx := context.Background()
	for _, table := range []string{
		new(code.LoginCode).TableName(ctx),
		new(code.RegistrationCode).TableName(ctx),
		new(continuity.Container).TableName(ctx),
		new(courier.MessageDispatch).TableName(),
		new(courier.Message).TableName(ctx),

		new(session.Device).TableName(ctx),
		new(session.Session).TableName(ctx),
		new(login.Flow).TableName(ctx),
		new(registration.Flow).TableName(ctx),
		new(settings.Flow).TableName(ctx),

		new(link.RecoveryToken).TableName(ctx),
		new(link.VerificationToken).TableName(ctx),
		new(code.RecoveryCode).TableName(ctx),
		new(code.VerificationCode).TableName(ctx),

		new(recovery.Flow).TableName(ctx),

		new(verification.Flow).TableName(ctx),

		new(errorx.ErrorContainer).TableName(ctx),

		new(identity.CredentialIdentifier).TableName(ctx),
		new(identity.Credentials).TableName(ctx),
		new(identity.VerifiableAddress).TableName(ctx),
		new(identity.RecoveryAddress).TableName(ctx),
		new(identity.Identity).TableName(ctx),
		new(identity.CredentialsTypeTable).TableName(ctx),
		new(sessiontokenexchange.Exchanger).TableName(),
		"networks",
		"schema_migration",
	} {
		if err := c.RawQuery("DROP TABLE IF EXISTS " + table).Exec(); err != nil {
			t.Logf(`Unable to clean up table "%s": %s`, table, err)
		}
	}
	t.Logf("Successfully cleaned up database: %s", c.Dialect.Name())
}
