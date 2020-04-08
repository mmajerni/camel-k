// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"fmt"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
)

// New returns a new RoundTripper based on the provided RoundTripper that has been
// setup to authenticate with the remote registry "reg", in the capacity
// laid out by the specified scopes.
func New(reg name.Registry, auth authn.Authenticator, t http.RoundTripper, scopes []string) (http.RoundTripper, error) {
	// The handshake:
	//  1. Use "t" to ping() the registry for the authentication challenge.
	//
	//  2a. If we get back a 200, then simply use "t".
	//
	//  2b. If we get back a 401 with a Basic challenge, then use a transport
	//     that just attachs auth each roundtrip.
	//
	//  2c. If we get back a 401 with a Bearer challenge, then use a transport
	//     that attaches a bearer token to each request, and refreshes is on 401s.
	//     Perform an initial refresh to seed the bearer token.

	// First we ping the registry to determine the parameters of the authentication handshake
	// (if one is even necessary).
	pr, err := ping(reg, t)
	if err != nil {
		return nil, err
	}

	// Wrap the given transport in transports that use an appropriate scheme,
	// (based on the ping response) and set the user agent.
	t = &useragentTransport{
		inner: &schemeTransport{
			scheme:   pr.scheme,
			registry: reg,
			inner:    t,
		},
	}

	switch pr.challenge.Canonical() {
	case anonymous:
		return t, nil
	case basic:
		return &basicTransport{inner: t, auth: auth, target: reg.RegistryStr()}, nil
	case bearer:
		// We require the realm, which tells us where to send our Basic auth to turn it into Bearer auth.
		realm, ok := pr.parameters["realm"]
		if !ok {
			return nil, fmt.Errorf("malformed www-authenticate, missing realm: %v", pr.parameters)
		}
		service, ok := pr.parameters["service"]
		if !ok {
			// If the service parameter is not specified, then default it to the registry
			// with which we are talking.
			service = reg.String()
		}
		bt := &bearerTransport{
			inner:    t,
			basic:    auth,
			realm:    realm,
			registry: reg,
			service:  service,
			scopes:   scopes,
			scheme:   pr.scheme,
		}
		if err := bt.refresh(); err != nil {
			return nil, err
		}
		return bt, nil
	default:
		return nil, fmt.Errorf("unrecognized challenge: %s", pr.challenge)
	}
}
