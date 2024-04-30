// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package verification

import (
	"context"
	"net/http"

	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/ui/node"
)

type FakeStrategy struct{}

var _ Strategy = new(FakeStrategy)

func (f FakeStrategy) VerificationStrategyID() string {
	return "fake"
}

func (f FakeStrategy) VerificationNodeGroup() node.UiNodeGroup {
	return "fake"
}

func (f FakeStrategy) PopulateVerificationMethod(*http.Request, *Flow) error {
	return nil
}

func (f FakeStrategy) Verify(_ http.ResponseWriter, _ *http.Request, _ *Flow) (err error) {
	return nil
}

func (f FakeStrategy) SendVerificationEmail(context.Context, *Flow, *identity.Identity, *identity.VerifiableAddress) error {
	return nil
}

func (f FakeStrategy) NodeGroup() node.UiNodeGroup {
	return "fake"
}
