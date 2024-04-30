// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package persistence

import (
	"context"
	"time"

	"github.com/ory/x/networkx"
	"my.com/secrets/internal/auth/domain/selfservice/sessiontokenexchange"

	"github.com/gofrs/uuid"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/popx"

	"my.com/secrets/internal/auth/domain/continuity"
	"my.com/secrets/internal/auth/domain/courier"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/errorx"
	"my.com/secrets/internal/auth/domain/selfservice/flow/login"
	"my.com/secrets/internal/auth/domain/selfservice/flow/recovery"
	"my.com/secrets/internal/auth/domain/selfservice/flow/registration"
	"my.com/secrets/internal/auth/domain/selfservice/flow/settings"
	"my.com/secrets/internal/auth/domain/selfservice/flow/verification"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/code"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/link"
	"my.com/secrets/internal/auth/domain/session"
)

type Provider interface {
	Persister() Persister
	SetPersister(Persister)
}

type Persister interface {
	continuity.Persister
	identity.PrivilegedPool
	registration.FlowPersister
	login.FlowPersister
	settings.FlowPersister
	courier.Persister
	session.Persister
	sessiontokenexchange.Persister
	errorx.Persister
	verification.FlowPersister
	recovery.FlowPersister
	link.RecoveryTokenPersister
	link.VerificationTokenPersister
	code.RecoveryCodePersister
	code.VerificationCodePersister
	code.RegistrationCodePersister
	code.LoginCodePersister

	CleanupDatabase(context.Context, time.Duration, time.Duration, int) error
	Close(context.Context) error
	Ping() error
	MigrationStatus(c context.Context) (popx.MigrationStatuses, error)
	MigrateDown(c context.Context, steps int) error
	MigrateUp(c context.Context) error
	Migrator() *popx.Migrator
	MigrationBox() *popx.MigrationBox
	GetConnection(ctx context.Context) *pop.Connection
	Transaction(ctx context.Context, callback func(ctx context.Context, connection *pop.Connection) error) error
	Networker
}

type Networker interface {
	WithNetworkID(sid uuid.UUID) Persister
	NetworkID(ctx context.Context) uuid.UUID
	DetermineNetwork(ctx context.Context) (*networkx.Network, error)
}
