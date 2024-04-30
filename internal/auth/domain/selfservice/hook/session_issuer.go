// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"net/http"

	"my.com/secrets/internal/auth/domain/identity"
	"my.com/secrets/internal/auth/domain/ui/node"

	"github.com/pkg/errors"

	"github.com/ory/x/otelx"
	"my.com/secrets/internal/auth/domain/driver/config"
	"my.com/secrets/internal/auth/domain/selfservice/flow"
	"my.com/secrets/internal/auth/domain/selfservice/flow/registration"
	"my.com/secrets/internal/auth/domain/selfservice/sessiontokenexchange"
	"my.com/secrets/internal/auth/domain/session"
	"my.com/secrets/internal/auth/domain/x"
)

var (
	_ registration.PostHookPostPersistExecutor = new(SessionIssuer)
)

type (
	sessionIssuerDependencies interface {
		session.ManagementProvider
		session.PersistenceProvider
		sessiontokenexchange.PersistenceProvider
		config.Provider
		x.WriterProvider
	}
	SessionIssuerProvider interface {
		HookSessionIssuer() *SessionIssuer
	}
	SessionIssuer struct {
		r sessionIssuerDependencies
	}
)

func NewSessionIssuer(r sessionIssuerDependencies) *SessionIssuer {
	return &SessionIssuer{r: r}
}

func (e *SessionIssuer) ExecutePostRegistrationPostPersistHook(w http.ResponseWriter, r *http.Request, a *registration.Flow, s *session.Session) error {
	return otelx.WithSpan(r.Context(), "selfservice.hook.SessionIssuer.ExecutePostRegistrationPostPersistHook", func(ctx context.Context) error {
		return e.executePostRegistrationPostPersistHook(w, r.WithContext(ctx), a, s)
	})
}

func (e *SessionIssuer) executePostRegistrationPostPersistHook(w http.ResponseWriter, r *http.Request, a *registration.Flow, s *session.Session) error {
	if a.Type == flow.TypeAPI {
		// We don't want to redirect with the code, if the flow was submitted with an ID token.
		// This is the case for Sign in with native Apple SDK or Google SDK.
		if s.AuthenticatedVia(identity.CredentialsTypeOIDC) && a.IDToken == "" {
			if handled, err := e.r.SessionManager().MaybeRedirectAPICodeFlow(w, r, a, s.ID, node.OpenIDConnectGroup); err != nil {
				return errors.WithStack(err)
			} else if handled {
				return nil
			}
		}

		a.AddContinueWith(flow.NewContinueWithSetToken(s.Token))
		e.r.Writer().Write(w, r, &registration.APIFlowResponse{
			Session:      s,
			Token:        s.Token,
			Identity:     s.Identity,
			ContinueWith: a.ContinueWithItems,
		})
		return errors.WithStack(registration.ErrHookAbortFlow)
	}

	// cookie is issued both for browser and for SPA flows
	if err := e.r.SessionManager().IssueCookie(r.Context(), w, r, s); err != nil {
		return err
	}

	// SPA flows additionally send the session
	if x.IsJSONRequest(r) {
		e.r.Writer().Write(w, r, &registration.APIFlowResponse{
			Session:      s,
			Identity:     s.Identity,
			ContinueWith: a.ContinueWithItems,
		})
		return errors.WithStack(registration.ErrHookAbortFlow)
	}

	return nil
}
