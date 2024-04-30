// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package login

import (
	"net/http"

	"my.com/secrets/internal/auth/domain/session"
)

func RequiresAAL2ForTest(e HookExecutor, r *http.Request, s *session.Session) (bool, error) {
	return e.requiresAAL2(r, s, nil) // *login.Flow is nil to avoid an import cycle
}
