// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// Flow State
//
// The state represents the state of the verification flow.
//
// - choose_method: ask the user to choose a method (e.g. recover account via email)
// - sent_email: the email has been sent to the user
// - passed_challenge: the request was successful and the recovery challenge was passed.
// - show_form: a form is shown to the user to perform the flow
// - success: the flow has been completed successfully
//
// swagger:enum selfServiceFlowState
type State string

// #nosec G101 -- only a key constant
// Define the various states for all flows.
// Recovery flows for V2 have different states (see below).
const (
	StateChooseMethod State = "choose_method"
	// Note: this state should actually be called `StateMessageSent`,
	// where a 'Message' is a code or link sent to an address (e.g. `email`, `sms`, etc).
	StateEmailSent       State = "sent_email"
	StatePassedChallenge State = "passed_challenge"
	StateShowForm        State = "show_form"
	StateSuccess         State = "success"

	// Recovery V2.
	// The initial state is different from `StateChooseMethod` to distinguish recovery v1 vs v2.
	// This avoids the issue of the feature flag being toggled while some recovery flows are on-going.
	// This would lead to an inconsistent state machine/logic for these flows.
	StateRecoveryAwaitingAddress        State = "recovery_awaiting_address"
	StateRecoveryAwaitingAddressChoice  State = "recovery_awaiting_address_choice"
	StateRecoveryAwaitingAddressConfirm State = "recovery_confirming_address"
	StateRecoveryAwaitingCode           State = "recovery_awaiting_code"
	// The final success state is the same as in Recovery V1 (`passed_challenge`).
	// Since this is the terminal state, it is not affected by toggling the feature flag.

	// State machine diagrams:
	// - ./state_recovery_v1.mermaid
	// - ./state_recovery_v2.mermaid
)

var states = []State{
	StateChooseMethod,
	StateEmailSent,
	StatePassedChallenge,
}

func indexOf(current State) int {
	for k, s := range states {
		if s == current {
			return k
		}
	}
	return 0
}

func HasReachedState(expected, actual State) bool {
	return indexOf(actual) >= indexOf(expected)
}

func IsStateRecoveryV2(state State) bool {
	return strings.HasPrefix(state.String(), "recovery_")
}

func NextState(current State) State {
	if current == StatePassedChallenge {
		return StatePassedChallenge
	}

	return states[indexOf(current)+1]
}

// For some reason using sqlxx.NullString as the State type does not work here.
// Reimplementing the Scanner interface on type State does work and allows
// the state to be NULL in the database.

// MarshalJSON returns m as the JSON encoding of m.
func (ns State) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ns))
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *State) UnmarshalJSON(data []byte) error {
	if ns == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	if len(data) == 0 {
		return nil
	}
	return errors.WithStack(json.Unmarshal(data, (*string)(ns)))
}

// Scan implements the Scanner interface.
func (ns *State) Scan(value interface{}) error {
	var v sql.NullString
	if err := (&v).Scan(value); err != nil {
		return err
	}
	*ns = State(v.String)
	return nil
}

// Value implements the driver Valuer interface.
func (ns State) Value() (driver.Value, error) {
	if len(ns) == 0 {
		return sql.NullString{}.Value()
	}
	return sql.NullString{Valid: true, String: string(ns)}.Value()
}

// String implements the Stringer interface.
func (ns State) String() string {
	return string(ns)
}
