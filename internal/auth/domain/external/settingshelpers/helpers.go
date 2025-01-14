// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package settingshelpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/external/testhelpers"
	"my.com/secrets/internal/auth/domain/selfservice/flow/settings"
)

func ExpectValidationError(t *testing.T, ts *httptest.Server, hc *http.Client, conf *config.Config, flow string, values func(url.Values)) string {
	isSPA := flow == "spa"
	isAPI := flow == "api"
	ctx := context.Background()
	return testhelpers.SubmitSettingsForm(t, isAPI, isSPA, hc, ts, values,
		testhelpers.ExpectStatusCode(isAPI || isSPA, http.StatusBadRequest, http.StatusOK),
		testhelpers.ExpectURL(isAPI || isSPA, ts.URL+settings.RouteSubmitFlow, conf.SelfServiceFlowSettingsUI(ctx).String()))
}
