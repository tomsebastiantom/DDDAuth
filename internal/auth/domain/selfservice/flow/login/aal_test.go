// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package login_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/flow"
	"my.com/secrets/internal/auth/domain/selfservice/flow/login"
)

func TestCheckAAL(t *testing.T) {
	f := &login.Flow{RequestedAAL: identity.AuthenticatorAssuranceLevel1}
	assert.NoError(t, login.CheckAAL(f, identity.AuthenticatorAssuranceLevel1))
	assert.ErrorIs(t, login.CheckAAL(f, identity.AuthenticatorAssuranceLevel2), flow.ErrStrategyNotResponsible)
}
