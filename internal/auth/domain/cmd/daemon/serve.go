// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package daemon

import (
	stdctx "context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/ory/x/otelx/semconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

	"github.com/ory/analytics-go/v5"
	"github.com/ory/graceful"
	"github.com/ory/x/healthx"
	"github.com/ory/x/metricsx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/otelx"
	prometheus "github.com/ory/x/prometheusx"
	"github.com/ory/x/reqlog"
	"github.com/ory/x/servicelocatorx"

	"my.com/secrets/internal/auth/domain/cmd/courier"
	"my.com/secrets/internal/auth/domain/driver"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/schema"
	"my.com/secrets/internal/auth/domain/selfservice/errorx"
	"my.com/secrets/internal/auth/domain/selfservice/flow/login"
	"my.com/secrets/internal/auth/domain/selfservice/flow/logout"
	"my.com/secrets/internal/auth/domain/selfservice/flow/recovery"
	"my.com/secrets/internal/auth/domain/selfservice/flow/registration"
	"my.com/secrets/internal/auth/domain/selfservice/flow/settings"
	"my.com/secrets/internal/auth/domain/selfservice/flow/verification"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/link"
	"my.com/secrets/internal/auth/domain/selfservice/strategy/oidc"
	"my.com/secrets/internal/auth/domain/session"
	"my.com/secrets/internal/auth/domain/x"
)

type options struct {
	ctx stdctx.Context
}

func NewOptions(ctx stdctx.Context, opts []Option) *options {
	o := new(options)
	o.ctx = ctx
	for _, f := range opts {
		f(o)
	}
	return o
}

type Option func(*options)

func WithContext(ctx stdctx.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func init() {
	graceful.DefaultShutdownTimeout = 120 * time.Second
}

func servePublic(r driver.Registry, cmd *cobra.Command, eg *errgroup.Group, slOpts *servicelocatorx.Options, opts []Option) {
	modifiers := NewOptions(cmd.Context(), opts)
	ctx := modifiers.ctx

	c := r.Config()
	l := r.Logger()
	n := negroni.New()

	for _, mw := range slOpts.HTTPMiddlewares() {
		n.UseFunc(mw)
	}

	publicLogger := reqlog.NewMiddlewareFromLogger(
		l,
		"public#"+c.SelfPublicURL(ctx).String(),
	)

	if r.Config().DisablePublicHealthRequestLog(ctx) {
		publicLogger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath)
	}

	n.UseFunc(semconv.Middleware)
	n.Use(publicLogger)
	n.Use(x.HTTPLoaderContextMiddleware(r))
	n.Use(sqa(ctx, cmd, r))

	n.Use(r.PrometheusManager())

	router := x.NewRouterPublic()
	csrf := x.NewCSRFHandler(router, r)

	// we need to always load the CORS middleware even if it is disabled, to allow hot-enabling CORS
	n.UseFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		cfg, enabled := r.Config().CORS(req.Context(), "public")
		if !enabled {
			next(w, req)
			return
		}
		cors.New(cfg).ServeHTTP(w, req, next)
	})

	n.UseFunc(x.CleanPath) // Prevent double slashes from breaking CSRF.
	r.WithCSRFHandler(csrf)
	n.UseHandler(r.CSRFHandler())

	// Disable CSRF for these endpoints
	csrf.DisablePath(healthx.AliveCheckPath)
	csrf.DisablePath(healthx.ReadyCheckPath)
	csrf.DisablePath(healthx.VersionPath)
	csrf.DisablePath(prometheus.MetricsPrometheusPath)

	r.RegisterPublicRoutes(ctx, router)
	r.PrometheusManager().RegisterRouter(router.Router)

	certs := c.GetTLSCertificatesForPublic(ctx)

	var handler http.Handler = n
	if tracer := r.Tracer(ctx); tracer.IsLoaded() {
		handler = otelx.TraceHandler(handler, otelhttp.WithTracerProvider(tracer.Provider()))
	}

	//#nosec G112 -- the correct settings are set by graceful.WithDefaults
	server := graceful.WithDefaults(&http.Server{
		Handler:           handler,
		TLSConfig:         &tls.Config{GetCertificate: certs, MinVersion: tls.VersionTLS12},
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	})
	addr := c.PublicListenOn(ctx)

	eg.Go(func() error {
		l.Printf("Starting the public httpd on: %s", addr)
		if err := graceful.GracefulContext(ctx, func() error {
			listener, err := networkx.MakeListener(addr, c.PublicSocketPermission(ctx))
			if err != nil {
				return err
			}

			if certs == nil {
				return server.Serve(listener)
			}
			return server.ServeTLS(listener, "", "")
		}, server.Shutdown); err != nil {
			if !errors.Is(err, context.Canceled) {
				l.Errorf("Failed to gracefully shutdown public httpd: %s", err)
				return err
			}
		}
		l.Println("Public httpd was shutdown gracefully")
		return nil
	})
}

func serveAdmin(r driver.Registry, cmd *cobra.Command, eg *errgroup.Group, slOpts *servicelocatorx.Options, opts []Option) {
	modifiers := NewOptions(cmd.Context(), opts)
	ctx := modifiers.ctx

	c := r.Config()
	l := r.Logger()
	n := negroni.New()

	for _, mw := range slOpts.HTTPMiddlewares() {
		n.UseFunc(mw)
	}

	adminLogger := reqlog.NewMiddlewareFromLogger(
		l,
		"admin#"+c.SelfPublicURL(ctx).String(),
	)

	if r.Config().DisableAdminHealthRequestLog(ctx) {
		adminLogger.ExcludePaths(x.AdminPrefix+healthx.AliveCheckPath, x.AdminPrefix+healthx.ReadyCheckPath, x.AdminPrefix+prometheus.MetricsPrometheusPath)
	}
	n.UseFunc(semconv.Middleware)
	n.Use(adminLogger)
	n.UseFunc(x.RedirectAdminMiddleware)
	n.Use(x.HTTPLoaderContextMiddleware(r))
	n.Use(sqa(ctx, cmd, r))
	n.Use(r.PrometheusManager())

	router := x.NewRouterAdmin()
	r.RegisterAdminRoutes(ctx, router)
	r.PrometheusManager().RegisterRouter(router.Router)

	n.UseHandler(router)
	certs := c.GetTLSCertificatesForAdmin(ctx)

	var handler http.Handler = n
	if tracer := r.Tracer(ctx); tracer.IsLoaded() {
		handler = otelx.TraceHandler(handler,
			otelhttp.WithTracerProvider(tracer.Provider()),
			otelhttp.WithFilter(func(req *http.Request) bool {
				return req.URL.Path != x.AdminPrefix+prometheus.MetricsPrometheusPath
			}),
		)
	}

	//#nosec G112 -- the correct settings are set by graceful.WithDefaults
	server := graceful.WithDefaults(&http.Server{
		Handler:           handler,
		TLSConfig:         &tls.Config{GetCertificate: certs, MinVersion: tls.VersionTLS12},
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       600 * time.Second,
	})

	addr := c.AdminListenOn(ctx)

	eg.Go(func() error {
		l.Printf("Starting the admin httpd on: %s", addr)
		if err := graceful.GracefulContext(ctx, func() error {
			listener, err := networkx.MakeListener(addr, c.AdminSocketPermission(ctx))
			if err != nil {
				return err
			}

			if certs == nil {
				return server.Serve(listener)
			}
			return server.ServeTLS(listener, "", "")
		}, server.Shutdown); err != nil {
			if !errors.Is(err, context.Canceled) {
				l.Errorf("Failed to gracefully shutdown admin httpd: %s", err)
				return err
			}
		}
		l.Println("Admin httpd was shutdown gracefully")
		return nil
	})
}

func sqa(ctx stdctx.Context, cmd *cobra.Command, d driver.Registry) *metricsx.Service {
	// Creates only ones
	// instance
	return metricsx.New(
		cmd,
		d.Logger(),
		d.Config().GetProvider(ctx),
		&metricsx.Options{
			Service:       "kratos",
			DeploymentId:  metricsx.Hash(d.Persister().NetworkID(ctx).String()),
			DBDialect:     d.Persister().GetConnection(ctx).Dialect.Details().Dialect,
			IsDevelopment: d.Config().IsInsecureDevMode(ctx),
			WriteKey:      "qQlI6q8Q4WvkzTjKQSor4sHYOikHIvvi",
			WhitelistedPaths: []string{
				"/",
				healthx.AliveCheckPath,
				healthx.ReadyCheckPath,
				healthx.VersionPath,

				oidc.RouteBase,

				login.RouteInitBrowserFlow,
				login.RouteInitAPIFlow,
				login.RouteGetFlow,
				login.RouteSubmitFlow,

				logout.RouteInitBrowserFlow,
				logout.RouteSubmitFlow,
				logout.RouteAPIFlow,

				registration.RouteInitBrowserFlow,
				registration.RouteInitAPIFlow,
				registration.RouteGetFlow,
				registration.RouteSubmitFlow,

				session.RouteWhoami,

				x.AdminPrefix + "/" + schema.SchemasPath,
				x.AdminPrefix + identity.RouteCollection,

				settings.RouteInitBrowserFlow,
				settings.RouteInitAPIFlow,
				settings.RouteGetFlow,
				settings.RouteSubmitFlow,

				verification.RouteInitAPIFlow,
				verification.RouteInitBrowserFlow,
				verification.RouteGetFlow,
				verification.RouteSubmitFlow,

				recovery.RouteInitAPIFlow,
				recovery.RouteInitBrowserFlow,
				recovery.RouteGetFlow,
				recovery.RouteSubmitFlow,

				link.RouteAdminCreateRecoveryLink,

				errorx.RouteGet,
				prometheus.MetricsPrometheusPath,
			},
			BuildVersion: config.Version,
			BuildHash:    config.Commit,
			BuildTime:    config.Date,
			Config: &analytics.Config{
				Endpoint:             "https://sqa.ory.sh",
				GzipCompressionLevel: 6,
				BatchMaxSize:         500 * 1000,
				BatchSize:            1000,
				Interval:             time.Hour * 6,
			},
		},
	)
}

func bgTasks(d driver.Registry, cmd *cobra.Command, opts []Option) error {
	modifiers := NewOptions(cmd.Context(), opts)
	ctx := modifiers.ctx

	if d.Config().IsBackgroundCourierEnabled(ctx) {
		return courier.Watch(ctx, d)
	}

	return nil
}

func ServeAll(d driver.Registry, slOpts *servicelocatorx.Options, opts []Option) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		mods := NewOptions(cmd.Context(), opts)
		ctx := mods.ctx

		g, ctx := errgroup.WithContext(ctx)
		cmd.SetContext(ctx)
		opts = append(opts, WithContext(ctx))

		servePublic(d, cmd, g, slOpts, opts)
		serveAdmin(d, cmd, g, slOpts, opts)
		g.Go(func() error {
			return bgTasks(d, cmd, opts)
		})
		return g.Wait()
	}
}