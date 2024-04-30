// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"context"
	"fmt"
	"testing"

	"my.com/secrets/internal/auth/domain/driver/config"
)

func StrategyEnable(t *testing.T, c *config.Config, strategy string, enable bool) {
	ctx := context.Background()
	c.MustSet(ctx, fmt.Sprintf("%s.%s.enabled", config.ViperKeySelfServiceStrategyConfig, strategy), enable)
}
