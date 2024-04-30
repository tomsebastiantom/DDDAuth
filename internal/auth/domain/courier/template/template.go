// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/x"
)

type Dependencies interface {
	CourierConfig() config.CourierConfigs
	x.HTTPClientProvider
}
