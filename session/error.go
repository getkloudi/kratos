// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"net/http"

	"github.com/ory/herodot"
)

func RespondWithJSONErrorOnAuthenticated(h herodot.Writer, err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.WriteError(w, r, err)
	}
}

var ErrNoSessionFound = herodot.ErrUnauthorized.WithReasonf("No valid session credentials found in the request.")
