// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/x"
)

func TestGenerateState(t *testing.T) {
	flowID := x.NewUUID().String()
	state := generateState(flowID).String()
	assert.NotEmpty(t, state)

	s, err := parseState(state)
	require.NoError(t, err)
	assert.Equal(t, flowID, s.FlowID)
	assert.NotEmpty(t, s.Data)
}
