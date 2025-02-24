// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package password

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/ory/x/otelx"

	"github.com/ory/kratos/text"

	"github.com/ory/kratos/ui/node"

	"github.com/ory/kratos/session"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/herodot"
	"github.com/ory/kratos/identity"
	"github.com/ory/kratos/schema"
	"github.com/ory/kratos/selfservice/flow"
	"github.com/ory/kratos/selfservice/flow/settings"
	"github.com/ory/kratos/x"
	"github.com/ory/x/decoderx"
)

func (s *Strategy) RegisterSettingsRoutes(_ *x.RouterPublic) {
}

func (s *Strategy) SettingsStrategyID() string {
	return identity.CredentialsTypePassword.String()
}

// Update Settings Flow with Password Method
//
// swagger:model updateSettingsFlowWithPasswordMethod
type updateSettingsFlowWithPasswordMethod struct {
	// Password is the updated password
	//
	// required: true
	Password string `json:"password"`

	// CSRFToken is the anti-CSRF token
	CSRFToken string `json:"csrf_token"`

	// Method
	//
	// Should be set to password when trying to update a password.
	//
	// required: true
	Method string `json:"method"`

	// Flow is flow ID.
	//
	// swagger:ignore
	Flow string `json:"flow"`

	// Transient data to pass along to any webhooks
	//
	// required: false
	TransientPayload json.RawMessage `json:"transient_payload,omitempty" form:"transient_payload"`
}

func (p *updateSettingsFlowWithPasswordMethod) GetFlowID() uuid.UUID {
	return x.ParseUUID(p.Flow)
}

func (p *updateSettingsFlowWithPasswordMethod) SetFlowID(rid uuid.UUID) {
	p.Flow = rid.String()
}

func (s *Strategy) Settings(ctx context.Context, w http.ResponseWriter, r *http.Request, f *settings.Flow, ss *session.Session) (_ *settings.UpdateContext, err error) {
	ctx, span := s.d.Tracer(ctx).Tracer().Start(ctx, "selfservice.strategy.password.Strategy.Settings")
	defer otelx.End(span, &err)

	var p updateSettingsFlowWithPasswordMethod
	ctxUpdate, err := settings.PrepareUpdate(s.d, w, r, f, ss, settings.ContinuityKey(s.SettingsStrategyID()), &p)
	if errors.Is(err, settings.ErrContinuePreviousAction) {
		return ctxUpdate, s.continueSettingsFlow(ctx, r, ctxUpdate, p)
	} else if err != nil {
		return ctxUpdate, s.handleSettingsError(ctx, w, r, ctxUpdate, p, err)
	}

	if err := flow.MethodEnabledAndAllowedFromRequest(r, f.GetFlowName(), s.SettingsStrategyID(), s.d); err != nil {
		return ctxUpdate, s.handleSettingsError(ctx, w, r, ctxUpdate, p, err)
	}

	if err := s.decodeSettingsFlow(r, &p); err != nil {
		return ctxUpdate, s.handleSettingsError(ctx, w, r, ctxUpdate, p, err)
	}

	// This does not come from the payload!
	p.Flow = ctxUpdate.Flow.ID.String()
	if err := s.continueSettingsFlow(ctx, r, ctxUpdate, p); err != nil {
		return ctxUpdate, s.handleSettingsError(ctx, w, r, ctxUpdate, p, err)
	}

	return ctxUpdate, nil
}

func (s *Strategy) decodeSettingsFlow(r *http.Request, dest interface{}) error {
	compiler, err := decoderx.HTTPRawJSONSchemaCompiler(settingsSchema)
	if err != nil {
		return errors.WithStack(err)
	}

	return decoderx.NewHTTP().Decode(r, dest, compiler,
		decoderx.HTTPKeepRequestBody(true),
		decoderx.HTTPDecoderAllowedMethods("POST", "GET"),
		decoderx.HTTPDecoderSetValidatePayloads(true),
		decoderx.HTTPDecoderJSONFollowsFormFormat(),
	)
}

func (s *Strategy) continueSettingsFlow(ctx context.Context, r *http.Request, ctxUpdate *settings.UpdateContext, p updateSettingsFlowWithPasswordMethod) error {
	if err := flow.MethodEnabledAndAllowed(ctx, flow.SettingsFlow, s.SettingsStrategyID(), p.Method, s.d); err != nil {
		return err
	}

	if err := flow.EnsureCSRF(s.d, r, ctxUpdate.Flow.Type, s.d.Config().DisableAPIFlowEnforcement(ctx), s.d.GenerateCSRFToken, p.CSRFToken); err != nil {
		return err
	}

	if ctxUpdate.Session.AuthenticatedAt.Add(s.d.Config().SelfServiceFlowSettingsPrivilegedSessionMaxAge(ctx)).Before(time.Now()) {
		return errors.WithStack(settings.NewFlowNeedsReAuth())
	}

	if len(p.Password) == 0 {
		return schema.NewRequiredError("#/password", "password")
	}

	hpw, errC := make(chan []byte), make(chan error)
	go func() {
		defer close(hpw)
		defer close(errC)
		h, err := s.d.Hasher(ctx).Generate(ctx, []byte(p.Password))
		if err != nil {
			errC <- err
			return
		}
		hpw <- h
	}()

	i, err := s.d.PrivilegedIdentityPool().GetIdentityConfidential(ctx, ctxUpdate.Session.Identity.ID)
	if err != nil {
		return err
	}

	i.UpsertCredentialsConfig(s.ID(), []byte("{}"), 0)
	if err := s.validateCredentials(ctx, i, p.Password); err != nil {
		return err
	}

	select {
	case err := <-errC:
		return err
	case h := <-hpw:
		co, err := json.Marshal(&identity.CredentialsPassword{HashedPassword: string(h)})
		if err != nil {
			return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to encode password options to JSON: %s", err))
		}
		i.UpsertCredentialsConfig(s.ID(), co, 0)
	}
	ctxUpdate.UpdateIdentity(i)

	return nil
}

func (s *Strategy) PopulateSettingsMethod(ctx context.Context, r *http.Request, _ *identity.Identity, f *settings.Flow) (err error) {
	_, span := s.d.Tracer(ctx).Tracer().Start(ctx, "selfservice.strategy.password.Strategy.PopulateSettingsMethod")
	defer otelx.End(span, &err)

	f.UI.SetCSRF(s.d.GenerateCSRFToken(r))
	f.UI.Nodes.Upsert(NewPasswordNode("password", node.InputAttributeAutocompleteNewPassword).WithMetaLabel(text.NewInfoNodeInputPassword()))
	f.UI.Nodes.Append(node.NewInputField("method", "password", node.PasswordGroup, node.InputAttributeTypeSubmit).WithMetaLabel(text.NewInfoNodeLabelSave()))

	return nil
}

func (s *Strategy) handleSettingsError(ctx context.Context, w http.ResponseWriter, r *http.Request, ctxUpdate *settings.UpdateContext, p updateSettingsFlowWithPasswordMethod, err error) error {
	// Do not pause flow if the flow type is an API flow as we can't save cookies in those flows.
	if e := new(settings.FlowNeedsReAuth); errors.As(err, &e) && ctxUpdate.Flow != nil && ctxUpdate.Flow.Type == flow.TypeBrowser {
		if err := s.d.ContinuityManager().Pause(ctx, w, r, settings.ContinuityKey(s.SettingsStrategyID()), settings.ContinuityOptions(p, ctxUpdate.GetSessionIdentity())...); err != nil {
			return err
		}
	}

	if ctxUpdate.Flow != nil {
		ctxUpdate.Flow.UI.ResetMessages()
		ctxUpdate.Flow.UI.SetCSRF(s.d.GenerateCSRFToken(r))
	}

	return err
}
