// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/ory/kratos/examples/go/pkg"

	ory "github.com/ory/client-go"
)

// If you use Open Source this would be:
//
// var client = pkg.NewSDKForSelfHosted("http://127.0.0.1:4433")
var client = pkg.NewSDK("playground")

func toSession() *ory.Session {
	// Create a temporary user
	email, password := pkg.RandomCredentials()
	_, sessionToken := pkg.CreateIdentityWithSession(client, email, password)

	session, res, err := client.FrontendAPI.ToSessionExecute(ory.FrontendAPIToSessionRequest{}.
		XSessionToken(sessionToken))
	pkg.SDKExitOnError(err, res)
	return session
}

func main() {
	session := toSession()
	pkg.PrintJSONPretty(session)
}
