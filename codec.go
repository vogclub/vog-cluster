package vogcluster

import (
	"encoding/json"
	"fmt"
)

// Validator is implemented by every message type in this package and
// reports whether the value's required fields are populated and any
// invariants hold.
type Validator interface {
	Validate() error
}

// Encode marshals msg to compact JSON after validating it. It is the
// recommended way to serialize messages for NATS publication: a sender
// that misses a required field gets an error before anything goes
// out on the wire.
func Encode(msg Validator) ([]byte, error) {
	if err := msg.Validate(); err != nil {
		return nil, fmt.Errorf("vogcluster: cannot encode invalid message: %w", err)
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("vogcluster: marshal: %w", err)
	}
	return data, nil
}

// Decode unmarshals data into out and then calls out.Validate(). out
// must be a pointer to a message type that implements Validator. A
// receiver that gets a malformed message gets an error instead of
// silently processing garbage.
func Decode(data []byte, out Validator) error {
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("vogcluster: unmarshal: %w", err)
	}
	if err := out.Validate(); err != nil {
		return fmt.Errorf("vogcluster: decoded message failed validation: %w", err)
	}
	return nil
}
