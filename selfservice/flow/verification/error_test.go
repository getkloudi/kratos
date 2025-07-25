// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package verification_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ory/kratos/x/nosurfx"

	"github.com/gofrs/uuid"

	"github.com/ory/x/jsonx"

	"github.com/ory/kratos/ui/node"

	"github.com/gobuffalo/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/assertx"
	"github.com/ory/x/urlx"

	"github.com/ory/herodot"

	"github.com/ory/kratos/driver/config"
	"github.com/ory/kratos/internal"
	"github.com/ory/kratos/internal/testhelpers"
	"github.com/ory/kratos/schema"
	"github.com/ory/kratos/selfservice/flow"
	"github.com/ory/kratos/selfservice/flow/verification"
	"github.com/ory/kratos/text"
	"github.com/ory/kratos/x"
)

func TestHandleError(t *testing.T) {
	ctx := context.Background()
	conf, reg := internal.NewFastRegistryWithMocks(t)
	conf.MustSet(ctx, config.ViperKeySelfServiceVerificationEnabled, true)

	public, _ := testhelpers.NewKratosServer(t, reg)

	router := http.NewServeMux()
	ts := httptest.NewServer(router)
	t.Cleanup(ts.Close)

	testhelpers.NewVerificationUIFlowEchoServer(t, reg)
	testhelpers.NewErrorTestServer(t, reg)

	h := reg.VerificationFlowErrorHandler()
	sdk := testhelpers.NewSDKClient(public)

	var verificationFlow *verification.Flow
	var flowError error
	var methodName node.UiNodeGroup
	router.HandleFunc("GET /error", func(w http.ResponseWriter, r *http.Request) {
		h.WriteFlowError(w, r, verificationFlow, methodName, flowError)
	})

	reset := func() {
		verificationFlow = nil
		flowError = nil
		methodName = ""
	}

	newFlow := func(t *testing.T, ttl time.Duration, ft flow.Type) *verification.Flow {
		t.Helper()
		req := &http.Request{URL: urlx.ParseOrPanic("/")}
		strategy, err := reg.GetActiveVerificationStrategy(context.Background())
		require.NoError(t, err)
		f, err := verification.NewFlow(conf, ttl, nosurfx.FakeCSRFToken, req, strategy, ft)
		require.NoError(t, err)
		require.NoError(t, reg.VerificationFlowPersister().CreateVerificationFlow(context.Background(), f))
		f, err = reg.VerificationFlowPersister().GetVerificationFlow(context.Background(), f.ID)
		require.NoError(t, err)
		return f
	}

	expectErrorUI := func(t *testing.T) (map[string]interface{}, *http.Response) {
		t.Helper()
		res, err := ts.Client().Get(ts.URL + "/error")
		require.NoError(t, err)
		defer res.Body.Close()
		require.Contains(t, res.Request.URL.String(), conf.SelfServiceFlowErrorURL(ctx).String()+"?id=")

		sse, _, err := sdk.FrontendAPI.GetFlowError(context.Background()).Id(res.Request.URL.Query().Get("id")).Execute()
		require.NoError(t, err)

		return sse.Error, nil
	}

	anHourAgo := time.Now().Add(-time.Hour)

	t.Run("case=error with nil flow defaults to error ui redirect", func(t *testing.T) {
		t.Cleanup(reset)

		flowError = herodot.ErrInternalServerError.WithReason("system error")
		methodName = node.UiNodeGroup(verification.VerificationStrategyLink)

		sse, _ := expectErrorUI(t)
		assertx.EqualAsJSON(t, flowError, sse)
	})

	t.Run("case=error with nil flow detects application/json", func(t *testing.T) {
		t.Cleanup(reset)

		flowError = herodot.ErrInternalServerError.WithReason("system error")
		methodName = node.UiNodeGroup(verification.VerificationStrategyLink)

		res, err := ts.Client().Do(testhelpers.NewHTTPGetJSONRequest(t, ts.URL+"/error"))
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Contains(t, res.Header.Get("Content-Type"), "application/json")
		assert.NotContains(t, res.Request.URL.String(), conf.SelfServiceFlowErrorURL(ctx).String()+"?id=")

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "system error")
	})

	for _, tc := range []struct {
		n string
		t flow.Type
	}{
		{"api", flow.TypeAPI},
		{"spa", flow.TypeBrowser},
	} {
		t.Run("flow="+tc.n, func(t *testing.T) {
			t.Run("case=expired error", func(t *testing.T) {
				t.Cleanup(reset)

				verificationFlow = newFlow(t, time.Minute, flow.TypeAPI)
				flowError = flow.NewFlowExpiredError(anHourAgo)
				methodName = node.UiNodeGroup(verification.VerificationStrategyLink)

				res, err := ts.Client().Do(testhelpers.NewHTTPGetJSONRequest(t, ts.URL+"/error"))
				require.NoError(t, err)
				defer res.Body.Close()
				require.Contains(t, res.Request.URL.String(), public.URL+verification.RouteGetFlow)
				require.Equal(t, http.StatusOK, res.StatusCode, "%+v", res.Request)

				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, int(text.ErrorValidationVerificationFlowExpired), int(gjson.GetBytes(body, "ui.messages.0.id").Int()), string(body))
				assert.NotEqual(t, verificationFlow.ID.String(), gjson.GetBytes(body, "id").String())
			})

			t.Run("case=validation error", func(t *testing.T) {
				t.Cleanup(reset)

				verificationFlow = newFlow(t, time.Minute, tc.t)
				flowError = schema.NewInvalidCredentialsError()
				methodName = node.UiNodeGroup(verification.VerificationStrategyLink)

				res, err := ts.Client().Do(testhelpers.NewHTTPGetJSONRequest(t, ts.URL+"/error"))
				require.NoError(t, err)
				defer res.Body.Close()
				require.Equal(t, http.StatusBadRequest, res.StatusCode)

				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, int(text.ErrorValidationInvalidCredentials), int(gjson.GetBytes(body, "ui.messages.0.id").Int()), "%s", body)
				assert.Equal(t, verificationFlow.ID.String(), gjson.GetBytes(body, "id").String())
			})

			t.Run("case=generic error", func(t *testing.T) {
				t.Cleanup(reset)

				verificationFlow = newFlow(t, time.Minute, tc.t)
				flowError = herodot.ErrInternalServerError.WithReason("system error")
				methodName = node.UiNodeGroup(verification.VerificationStrategyLink)

				res, err := ts.Client().Do(testhelpers.NewHTTPGetJSONRequest(t, ts.URL+"/error"))
				require.NoError(t, err)
				defer res.Body.Close()
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)

				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.JSONEq(t, x.MustEncodeJSON(t, flowError), gjson.GetBytes(body, "error").Raw)
			})
		})
	}

	t.Run("flow=browser", func(t *testing.T) {
		expectVerificationUI := func(t *testing.T) (*verification.Flow, *http.Response) {
			res, err := ts.Client().Get(ts.URL + "/error")
			require.NoError(t, err)
			defer res.Body.Close()
			assert.Contains(t, res.Request.URL.String(), conf.SelfServiceFlowVerificationUI(ctx).String()+"?flow=")

			vf, err := reg.VerificationFlowPersister().GetVerificationFlow(context.Background(), uuid.FromStringOrNil(res.Request.URL.Query().Get("flow")))
			require.NoError(t, err)
			return vf, res
		}

		t.Run("case=expired error", func(t *testing.T) {
			t.Cleanup(reset)

			verificationFlow = &verification.Flow{Type: flow.TypeBrowser}
			flowError = flow.NewFlowExpiredError(anHourAgo)
			methodName = node.LinkGroup

			lf, _ := expectVerificationUI(t)
			require.Len(t, lf.UI.Messages, 1, "%s", jsonx.TestMarshalJSONString(t, lf))
			assert.Equal(t, int(text.ErrorValidationVerificationFlowExpired), int(lf.UI.Messages[0].ID))
		})

		t.Run("case=validation error", func(t *testing.T) {
			t.Cleanup(reset)

			verificationFlow = newFlow(t, time.Minute, flow.TypeBrowser)
			flowError = schema.NewInvalidCredentialsError()
			methodName = node.LinkGroup

			lf, _ := expectVerificationUI(t)
			require.NotEmpty(t, lf.UI, x.MustEncodeJSON(t, lf))
			require.Len(t, lf.UI.Messages, 1, x.MustEncodeJSON(t, lf))
			assert.Equal(t, int(text.ErrorValidationInvalidCredentials), int(lf.UI.Messages[0].ID), x.MustEncodeJSON(t, lf))
		})

		t.Run("case=generic error", func(t *testing.T) {
			t.Cleanup(reset)

			verificationFlow = newFlow(t, time.Minute, flow.TypeBrowser)
			flowError = herodot.ErrInternalServerError.WithReason("system error")
			methodName = node.LinkGroup

			sse, _ := expectErrorUI(t)
			assertx.EqualAsJSON(t, flowError, sse)
		})

		t.Run("case=fails to retry flow if recovery strategy id is not valid", func(t *testing.T) {
			t.Cleanup(func() {
				reset()
				conf.MustSet(ctx, config.ViperKeySelfServiceVerificationUse, "code")
			})

			verificationFlow = newFlow(t, 0, flow.TypeBrowser)
			verificationFlow.Active = "not-valid"
			flowError = flow.NewFlowExpiredError(anHourAgo)

			conf.MustSet(ctx, config.ViperKeySelfServiceVerificationUse, "not-valid")
			sse, _ := expectErrorUI(t)
			testhelpers.SnapshotTExcept(t, sse, nil)
		})
	})
}
