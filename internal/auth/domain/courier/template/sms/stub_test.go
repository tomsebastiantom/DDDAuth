// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sms_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/courier/template/sms"
	"my.com/secrets/internal/auth/domain/external"
)

func TestNewTestStub(t *testing.T) {
	_, reg := external.NewFastRegistryWithMocks(t)

	const (
		expectedPhone = "+12345678901"
		expectedBody  = "test sms"
	)

	tpl := sms.NewTestStub(reg, &sms.TestStubModel{To: expectedPhone, Body: expectedBody})

	actualBody, err := tpl.SMSBody(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "stub sms body test sms\n", actualBody)

	actualPhone, err := tpl.PhoneNumber()
	require.NoError(t, err)
	assert.Equal(t, expectedPhone, actualPhone)
}
