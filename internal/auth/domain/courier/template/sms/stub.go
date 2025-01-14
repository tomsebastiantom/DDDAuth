// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sms

import (
	"context"
	"encoding/json"
	"os"

	"my.com/secrets/internal/auth/domain/courier/template"
)

type (
	TestStub struct {
		d template.Dependencies
		m *TestStubModel
	}

	TestStubModel struct {
		To       string                 `json:"to"`
		Body     string                 `json:"body"`
		Identity map[string]interface{} `json:"identity"`
	}
)

func NewTestStub(d template.Dependencies, m *TestStubModel) *TestStub {
	return &TestStub{d: d, m: m}
}

func (t *TestStub) PhoneNumber() (string, error) {
	return t.m.To, nil
}

func (t *TestStub) SMSBody(ctx context.Context) (string, error) {
	return template.LoadText(ctx, t.d, os.DirFS(t.d.CourierConfig().CourierTemplatesRoot(ctx)), "otp/test_stub/sms.body.gotmpl", "otp/test_stub/sms.body*", t.m, "")
}

func (t *TestStub) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.m)
}

func (t *TestStub) TemplateType() template.TemplateType {
	return template.TypeTestStub
}
