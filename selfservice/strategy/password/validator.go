// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package password

import (
	"bufio"
	"context"
	"crypto/sha1" //#nosec G505 -- sha1 is used for k-anonymity
	stderrs "errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace/noop"

	"github.com/ory/kratos/text"

	"github.com/arbovm/levenshtein"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"

	"github.com/ory/herodot"
	"github.com/ory/kratos/driver/config"
	"github.com/ory/x/httpx"
	"github.com/ory/x/otelx"
)

const hashCacheItemTTL = time.Hour

// Validator implements a validation strategy for passwords. One example is that the password
// has to have at least 6 characters and at least one lower and one uppercase password.
type Validator interface {
	// Validate returns nil if the password is passing the validation strategy and an error otherwise. If a validation error
	// occurs, a regular error will be returned. If some other type of error occurs (e.g. HTTP request failed), an error
	// of type *herodot.DefaultError will be returned.
	Validate(ctx context.Context, identifier, password string) error
}

type ValidationProvider interface {
	PasswordValidator() Validator
}

var (
	_                       Validator = new(DefaultPasswordValidator)
	ErrNetworkFailure                 = stderrs.New("unable to check if password has been leaked because an unexpected network error occurred")
	ErrUnexpectedStatusCode           = stderrs.New("unexpected status code")
)

// DefaultPasswordValidator implements Validator. It is based on best
// practices as defined in the following blog posts:
//
// - https://www.troyhunt.com/passwords-evolved-authentication-guidance-for-the-modern-era/
// - https://www.microsoft.com/en-us/research/wp-content/uploads/2016/06/Microsoft_Password_Guidance-1.pdf
//
// Additionally passwords are being checked against Troy Hunt's
// [haveibeenpwnd](https://haveibeenpwned.com/API/v2#SearchingPwnedPasswordsByRange) service to check if the
// password has been breached in a previous data leak using k-anonymity.
type DefaultPasswordValidator struct {
	reg    validatorDependencies
	Client *retryablehttp.Client
	hashes *ristretto.Cache[string, int64]

	minIdentifierPasswordDist            int
	maxIdentifierPasswordSubstrThreshold float32
}

type validatorDependencies interface {
	config.Provider
}

func NewDefaultPasswordValidatorStrategy(reg validatorDependencies) (*DefaultPasswordValidator, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, int64]{
		NumCounters:        10 * 10000,
		MaxCost:            60 * 10000, // BCrypt hash size is 60 bytes
		BufferItems:        64,
		IgnoreInternalCost: true,
	})
	// sanity check - this should never happen unless above configuration variables are invalid
	if err != nil {
		return nil, errors.Wrap(err, "error while setting up validator cache")
	}
	return &DefaultPasswordValidator{
		Client: httpx.NewResilientClient(
			httpx.ResilientClientWithConnectionTimeout(time.Second),
			// Tracing still works correctly even though we pass a no-op tracer
			// here, because the otelhttp package will preferentially use the
			// tracer from the incoming request context over this one.
			httpx.ResilientClientWithTracer(noop.NewTracerProvider().Tracer("github.com/ory/kratos/selfservice/strategy/password"))),
		reg:                       reg,
		hashes:                    cache,
		minIdentifierPasswordDist: 5, maxIdentifierPasswordSubstrThreshold: 0.5}, nil
}

func b20(src []byte) string {
	return fmt.Sprintf("%X", src)
}

// code inspired by https://rosettacode.org/wiki/Longest_Common_Substring#Go
func lcsLength(a, b string) int {
	lengths := make([]int, len(a)*len(b))
	greatestLength := 0
	for i, x := range a {
		for j, y := range b {
			if x == y {
				curr := 1
				if i != 0 && j != 0 {
					curr = lengths[(i-1)*len(b)+j-1] + 1
				}

				if curr > greatestLength {
					greatestLength = curr
				}
				lengths[i*len(b)+j] = curr
			}
		}
	}
	return greatestLength
}

func (s *DefaultPasswordValidator) fetch(ctx context.Context, hpw []byte, apiDNSName string) (int64, error) {
	prefix := fmt.Sprintf("%X", hpw)[0:5]
	loc := fmt.Sprintf("https://%s/range/%s", apiDNSName, prefix)
	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", loc, nil)
	if err != nil {
		return 0, err
	}
	res, err := s.Client.Do(req)
	if err != nil {
		return 0, errors.Wrapf(ErrNetworkFailure, "%s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, errors.Wrapf(ErrUnexpectedStatusCode, "%d", res.StatusCode)
	}

	var thisCount int64

	sc := bufio.NewScanner(res.Body)
	for sc.Scan() {
		row := sc.Text()
		result := strings.Split(strings.TrimSpace(row), ":")

		// We assume a count of 1. HIBP API sometimes responds without the
		// colon, so we just assume that the leak count is one.
		//
		// See https://github.com/ory/kratos/issues/2145
		count := int64(1)
		if len(result) == 2 {
			count, err = strconv.ParseInt(strings.ReplaceAll(result[1], ",", ""), 10, 64)
			if err != nil {
				return 0, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Expected password hash to contain a count formatted as int but got: %s", result[1]))
			}
		}

		s.hashes.SetWithTTL(prefix+result[0], count, 1, hashCacheItemTTL)
		if prefix+result[0] == b20(hpw) {
			thisCount = count
		}
	}

	if err := sc.Err(); err != nil {
		return 0, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to initialize string scanner: %s", err))
	}

	s.hashes.SetWithTTL(b20(hpw), thisCount, 1, hashCacheItemTTL)
	return thisCount, nil
}

func (s *DefaultPasswordValidator) Validate(ctx context.Context, identifier, password string) error {
	return otelx.WithSpan(ctx, "password.DefaultPasswordValidator.Validate", func(ctx context.Context) error {
		return s.validate(ctx, identifier, password)
	})
}

func (s *DefaultPasswordValidator) validate(ctx context.Context, identifier, password string) error {
	passwordPolicyConfig := s.reg.Config().PasswordPolicyConfig(ctx)

	//nolint:gosec // disable G115
	if len(password) < int(passwordPolicyConfig.MinPasswordLength) {
		//nolint:gosec // disable G115
		return text.NewErrorValidationPasswordMinLength(int(passwordPolicyConfig.MinPasswordLength), len(password))
	}

	if passwordPolicyConfig.IdentifierSimilarityCheckEnabled && len(identifier) > 0 {
		compIdentifier, compPassword := strings.ToLower(identifier), strings.ToLower(password)
		dist := levenshtein.Distance(compIdentifier, compPassword)
		lcs := float32(lcsLength(compIdentifier, compPassword)) / float32(len(compPassword))
		if dist < s.minIdentifierPasswordDist || lcs > s.maxIdentifierPasswordSubstrThreshold {
			return text.NewErrorValidationPasswordIdentifierTooSimilar()
		}
	}

	if !passwordPolicyConfig.HaveIBeenPwnedEnabled {
		return nil
	}

	//#nosec G401 -- sha1 is used for k-anonymity
	h := sha1.New()
	if _, err := h.Write([]byte(password)); err != nil {
		return err
	}
	hpw := h.Sum(nil)

	c, ok := s.hashes.Get(b20(hpw))
	if !ok {
		var err error
		c, err = s.fetch(ctx, hpw, passwordPolicyConfig.HaveIBeenPwnedHost)
		if (errors.Is(err, ErrNetworkFailure) || errors.Is(err, ErrUnexpectedStatusCode)) && passwordPolicyConfig.IgnoreNetworkErrors {
			return nil
		} else if err != nil {
			return err
		}
	}

	//nolint:gosec // disable G115
	if c > int64(s.reg.Config().PasswordPolicyConfig(ctx).MaxBreaches) {
		return text.NewErrorValidationPasswordTooManyBreaches(c)
	}

	return nil
}
