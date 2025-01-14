// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"fmt"

	"my.com/secrets/internal/auth/domain/selfservice/strategy/code"
)

var CodeRegex = fmt.Sprintf(`(\d{%d})`, code.CodeLength)
