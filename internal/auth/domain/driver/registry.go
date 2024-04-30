// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"io/fs"

	"github.com/ory/x/contextx"
	"github.com/ory/x/jsonnetsecure"
	"github.com/ory/x/otelx"
	prometheus "github.com/ory/x/prometheusx"
	"my.com/secrets/internal/auth/domain/selfservice/sessiontokenexchange"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"

	"github.com/ory/nosurf"

	"github.com/ory/x/logrusx"

	"my.com/secrets/internal/auth/domain/continuity"
	"my.com/secrets/internal/auth/domain/courier"
	"my.com/secrets/internal/auth/domain/hash"
	"my.com/secrets/internal/auth/domain/schema"
	"my.com/secrets/internal/auth/domain/selfservice/flow/recovery"
	"my.com/secrets/internal/auth/domain/selfservice/flow/settings"
	"my.com/secrets/internal/auth/domain/selfservice/flow/verification"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/code"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/link"

	"github.com/ory/x/healthx"

	"my.com/secrets/internal/auth/domain/persistence"
	"my.com/secrets/internal/auth/domain/selfservice/flow/login"
	"my.com/secrets/internal/auth/domain/selfservice/flow/logout"
	"my.com/secrets/internal/auth/domain/selfservice/flow/registration"

	"my.com/secrets/internal/auth/domain/x"

	"github.com/ory/x/dbal"

	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/selfservice/errorx"
	password2 "my.com/secrets/internal/auth/domain/selfservice/strategy/password"
	"my.com/secrets/internal/auth/domain/session"
)

type Registry interface {
	dbal.Driver

	Init(ctx context.Context, ctxer contextx.Contextualizer, opts ...RegistryOption) error

	WithLogger(l *logrusx.Logger) Registry
	WithJsonnetVMProvider(jsonnetsecure.VMProvider) Registry

	WithCSRFHandler(c nosurf.Handler)
	WithCSRFTokenGenerator(cg x.CSRFToken)

	MetricsHandler() *prometheus.Handler
	HealthHandler(ctx context.Context) *healthx.Handler
	CookieManager(ctx context.Context) sessions.StoreExact
	ContinuityCookieManager(ctx context.Context) sessions.StoreExact

	RegisterRoutes(ctx context.Context, public *x.RouterPublic, admin *x.RouterAdmin)
	RegisterPublicRoutes(ctx context.Context, public *x.RouterPublic)
	RegisterAdminRoutes(ctx context.Context, admin *x.RouterAdmin)
	PrometheusManager() *prometheus.MetricsManager
	Tracer(context.Context) *otelx.Tracer
	SetTracer(*otelx.Tracer)

	config.Provider
	CourierConfig() config.CourierConfigs
	WithConfig(c *config.Config) Registry
	WithContextualizer(ctxer contextx.Contextualizer) Registry

	x.CSRFProvider
	x.WriterProvider
	x.LoggingProvider
	x.HTTPClientProvider
	jsonnetsecure.VMProvider

	continuity.ManagementProvider
	continuity.PersistenceProvider

	courier.Provider

	persistence.Provider

	errorx.ManagementProvider
	errorx.HandlerProvider
	errorx.PersistenceProvider

	hash.HashProvider

	identity.HandlerProvider
	identity.ValidationProvider
	identity.PoolProvider
	identity.PrivilegedPoolProvider
	identity.ManagementProvider
	identity.ActiveCredentialsCounterStrategyProvider

	courier.HandlerProvider
	courier.PersistenceProvider

	schema.HandlerProvider
	schema.IdentityTraitsProvider

	password2.ValidationProvider

	session.HandlerProvider
	session.ManagementProvider
	session.PersistenceProvider
	session.TokenizerProvider

	settings.HandlerProvider
	settings.ErrorHandlerProvider
	settings.FlowPersistenceProvider
	settings.StrategyProvider

	login.FlowPersistenceProvider
	login.ErrorHandlerProvider
	login.HooksProvider
	login.HookExecutorProvider
	login.HandlerProvider
	login.StrategyProvider

	logout.HandlerProvider

	registration.FlowPersistenceProvider
	registration.ErrorHandlerProvider
	registration.HooksProvider
	registration.HookExecutorProvider
	registration.HandlerProvider
	registration.StrategyProvider

	verification.FlowPersistenceProvider
	verification.ErrorHandlerProvider
	verification.HandlerProvider
	verification.StrategyProvider

	sessiontokenexchange.PersistenceProvider

	link.SenderProvider
	link.VerificationTokenPersistenceProvider
	link.RecoveryTokenPersistenceProvider

	code.SenderProvider
	code.RecoveryCodePersistenceProvider

	recovery.FlowPersistenceProvider
	recovery.ErrorHandlerProvider
	recovery.HandlerProvider
	recovery.StrategyProvider

	x.CSRFTokenGeneratorProvider
}

func NewRegistryFromDSN(ctx context.Context, c *config.Config, l *logrusx.Logger) (Registry, error) {
	driver, err := dbal.GetDriverFor(c.DSN(ctx))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	tracer, err := otelx.New("Ory Kratos", l, c.Tracing(ctx))
	if err != nil {
		l.WithError(err).Fatalf("failed to initialize tracer")
		tracer = otelx.NewNoop(l, c.Tracing(ctx))
	}
	registry.SetTracer(tracer)

	return registry.WithLogger(l).WithConfig(c), nil
}

type options struct {
	skipNetworkInit         bool
	config                  *config.Config
	replaceTracer           func(*otelx.Tracer) *otelx.Tracer
	inspect                 func(Registry) error
	extraMigrations         []fs.FS
	replacementStrategies   []NewStrategy
	extraHooks              map[string]func(config.SelfServiceHook) any
	disableMigrationLogging bool
	jsonnetPool             jsonnetsecure.Pool
}

type RegistryOption func(*options)

func SkipNetworkInit(o *options) {
	o.skipNetworkInit = true
}

func WithJsonnetPool(pool jsonnetsecure.Pool) RegistryOption {
	return func(o *options) {
		o.jsonnetPool = pool
	}
}

func WithConfig(config *config.Config) RegistryOption {
	return func(o *options) {
		o.config = config
	}
}

func ReplaceTracer(f func(*otelx.Tracer) *otelx.Tracer) RegistryOption {
	return func(o *options) {
		o.replaceTracer = f
	}
}

type NewStrategy func(deps any) any

// WithReplaceStrategies adds a strategy to the registry. This is useful if you want to
// add a custom strategy to the registry. Default strategies with the same
// name/ID will be overwritten.
func WithReplaceStrategies(s ...NewStrategy) RegistryOption {
	return func(o *options) {
		o.replacementStrategies = append(o.replacementStrategies, s...)
	}
}

func WithExtraHooks(hooks map[string]func(config.SelfServiceHook) any) RegistryOption {
	return func(o *options) {
		o.extraHooks = hooks
	}
}

func Inspect(f func(reg Registry) error) RegistryOption {
	return func(o *options) {
		o.inspect = f
	}
}

func WithExtraMigrations(m ...fs.FS) RegistryOption {
	return func(o *options) {
		o.extraMigrations = append(o.extraMigrations, m...)
	}
}

func WithDisabledMigrationLogging() RegistryOption {
	return func(o *options) {
		o.disableMigrationLogging = true
	}
}

func newOptions(os []RegistryOption) *options {
	o := new(options)
	for _, f := range os {
		f(o)
	}
	return o
}
