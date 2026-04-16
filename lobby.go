package vogcluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// LobbyOutput is the envelope a game instance publishes on
// SubjectLobbyOutput(lobby_instance_id) for messages targeted at
// specific clients (rather than broadcast to a whole room).
//
// Lobby looks up matching connections in its in-memory map and
// delivers Payload over WS or SSE to them.
type LobbyOutput struct {
	// Event identifies the kind of message.
	Event string `json:"event"`

	// TargetUserIDs lists users that should receive the payload.
	// Lobby delivers to every active connection of each user.
	// Either TargetUserIDs or TargetConnIDs must be non-empty.
	TargetUserIDs []string `json:"target_user_ids,omitempty"`

	// TargetConnIDs lists specific connection IDs to deliver to.
	// Useful for replying to a single client request.
	TargetConnIDs []string `json:"target_conn_ids,omitempty"`

	// ClientNonce echoes the nonce from the originating RoomInput
	// when this message is a reply to a specific client request.
	ClientNonce string `json:"client_nonce,omitempty"`

	// EmittedAt is the wall-clock time the game instance generated
	// the message.
	EmittedAt time.Time `json:"emitted_at,omitempty"`

	// Payload is the event-specific body, opaque JSON.
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Validate reports whether the envelope is well-formed.
func (m LobbyOutput) Validate() error {
	if m.Event == "" {
		return errors.New("vogcluster: LobbyOutput.event is required")
	}
	if len(m.TargetUserIDs) == 0 && len(m.TargetConnIDs) == 0 {
		return errors.New("vogcluster: LobbyOutput requires at least one target_user_id or target_conn_id")
	}
	return nil
}

// SessionEvent is published by lobby instances to track user room presence.
// Published on SubjectLobbySessionEnter / SubjectLobbySessionExit.
type SessionEvent struct {
	UserID     string `json:"user_id"`
	RoomID     string `json:"room_id"`
	Action     string `json:"action"`
	IP         string `json:"ip,omitempty"`
	ClientType string `json:"client_type,omitempty"`
}

// Validate implements Validator.
func (e *SessionEvent) Validate() error {
	if e.UserID == "" {
		return fmt.Errorf("session event: user_id required")
	}
	if e.RoomID == "" {
		return fmt.Errorf("session event: room_id required")
	}
	if e.Action == "" {
		return fmt.Errorf("session event: action required")
	}
	return nil
}
