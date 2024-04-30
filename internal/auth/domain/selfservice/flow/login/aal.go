// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package login

import (
	"github.com/pkg/errors"

	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/flow"
)

func CheckAAL(f *Flow, expected identity.AuthenticatorAssuranceLevel) error {
	if f.RequestedAAL != expected {
		return errors.WithStack(flow.ErrStrategyNotResponsible)
	}
	return nil
}
