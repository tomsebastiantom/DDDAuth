// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my.com/secrets/internal/auth/domain/courier/template/sms"
	"my.com/secrets/internal/auth/domain/external"
)

func TestNewOTPMessage(t *testing.T) {
	_, reg := external.NewFastRegistryWithMocks(t)

	const (
		expectedPhone = "+12345678901"
		otp           = "012345"
	)

	tpl := sms.NewVerificationCodeValid(reg, &sms.VerificationCodeValidModel{To: expectedPhone, VerificationCode: otp})

	expectedBody := fmt.Sprintf("Your verification code is: %s\n", otp)

	actualBody, err := tpl.SMSBody(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)

	actualPhone, err := tpl.PhoneNumber()
	require.NoError(t, err)
	assert.Equal(t, expectedPhone, actualPhone)
}
