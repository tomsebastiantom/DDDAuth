// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"time"

	"my.com/secrets/internal/auth/domain/text"
	"my.com/secrets/internal/auth/domain/ui/node"

	"github.com/ory/x/sqlxx"
)

// CredentialsConfig is the struct that is being used as part of the identity credentials.
type CredentialsLookupConfig struct {
	// List of recovery codes
	RecoveryCodes []RecoveryCode `json:"recovery_codes"`
}

func (c *CredentialsLookupConfig) ToNode() *node.Node {
	messages := make([]text.Message, len(c.RecoveryCodes))
	formatted := make([]string, len(c.RecoveryCodes))
	for k, code := range c.RecoveryCodes {
		if time.Time(code.UsedAt).IsZero() {
			messages[k] = *text.NewInfoSelfServiceSettingsLookupSecret(code.Code)
			formatted[k] = code.Code
		} else {
			messages[k] = *text.NewInfoSelfServiceSettingsLookupSecretUsed(time.Time(code.UsedAt).UTC())
			formatted[k] = "used"
		}
	}

	return node.NewTextField(node.LookupCodes, text.NewInfoSelfServiceSettingsLookupSecretList(formatted, messages), node.LookupGroup).
		WithMetaLabel(text.NewInfoSelfServiceSettingsLookupSecretsLabel())
}

type RecoveryCode struct {
	// A recovery code
	Code string `json:"code"`

	// UsedAt indicates whether and when a recovery code was used.
	UsedAt sqlxx.NullTime `json:"used_at,omitempty"`
}
