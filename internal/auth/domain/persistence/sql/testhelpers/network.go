// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"context"
	"testing"

	db "github.com/gofrs/uuid"

	courier "my.com/secrets/internal/auth/domain/courier/test"
	"my.com/secrets/internal/auth/domain/persistence"
	"my.com/secrets/internal/auth/domain/external/testhelpers"
)

func DefaultNetworkWrapper(p persistence.Persister) (courier.NetworkWrapper, courier.NetworkWrapper) {
	return func(t *testing.T, ctx context.Context) (db.UUID, courier.PersisterWrapper) {
			return testhelpers.NewNetworkUnlessExisting(t, ctx, p)
		}, func(t *testing.T, ctx context.Context) (db.UUID, courier.PersisterWrapper) {
			return testhelpers.NewNetwork(t, ctx, p)
		}
}
