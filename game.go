package vogcluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// RoomInput is the envelope Lobby publishes on
// SubjectGameRoomInput(room_id) when forwarding a player action to the
// game instance hosting the room. The Payload is opaque JSON whose
// shape depends on Action and is interpreted by game-side logic.
type RoomInput struct {
	// Action identifies the kind of input ("room.enter", "chat",
	// "table.move", "seek.create", etc.). Game instances dispatch
	// on this value.
	Action string `json:"action"`

	// UserID is the authenticated user identifier.
	UserID string `json:"user_id"`

	// LobbyInstanceID is the Lobby that received this client message.
	// The game instance uses it to address responses via
	// SubjectLobbyOutput.
	LobbyInstanceID string `json:"lobby_instance_id"`

	// ConnectionID identifies the specific WebSocket connection
	// inside that Lobby. Lobby uses this to find the right WS to
	// send a targeted reply to.
	ConnectionID string `json:"connection_id"`

	// ClientNonce is an opaque token the client may include to
	// correlate responses to requests. Game instances echo it back
	// in any RoomBroadcast or LobbyOutput they emit in reply.
	ClientNonce string `json:"client_nonce,omitempty"`

	// ReceivedAt is the wall-clock time the Lobby received the
	// client message.
	ReceivedAt time.Time `json:"received_at,omitempty"`

	// Payload is the action-specific body, opaque JSON.
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Validate reports whether the envelope is well-formed.
func (m RoomInput) Validate() error {
	if m.Action == "" {
		return errors.New("vogcluster: RoomInput.action is required")
	}
	if m.UserID == "" {
		return errors.New("vogcluster: RoomInput.user_id is required")
	}
	if m.LobbyInstanceID == "" {
		return errors.New("vogcluster: RoomInput.lobby_instance_id is required")
	}
	if m.ConnectionID == "" {
		return errors.New("vogcluster: RoomInput.connection_id is required")
	}
	return nil
}

// RoomBroadcast is the envelope a game instance publishes on
// SubjectGameRoomBroadcast(room_id) to fan out an event to every
// Lobby instance with clients in the room.
type RoomBroadcast struct {
	// Event identifies the kind of event ("player.entered",
	// "table.move", "room.frozen", etc.).
	Event string `json:"event"`

	// RoomID is the room this broadcast originates from. Duplicated
	// here so subscribers don't need to parse the subject.
	RoomID string `json:"room_id"`

	// Sequence is a monotonically increasing per-room counter that
	// lets Lobby/clients detect gaps and request a resync.
	Sequence int64 `json:"sequence"`

	// EmittedAt is the wall-clock time the game instance generated
	// the event.
	EmittedAt time.Time `json:"emitted_at,omitempty"`

	// Payload is the event-specific body, opaque JSON.
	Payload json.RawMessage `json:"payload,omitempty"`

	// ExcludeUsers may list user IDs that should NOT receive the
	// broadcast (e.g. the sender of a chat message).
	ExcludeUsers []string `json:"exclude_users,omitempty"`
}

// Validate reports whether the envelope is well-formed.
func (m RoomBroadcast) Validate() error {
	if m.Event == "" {
		return errors.New("vogcluster: RoomBroadcast.event is required")
	}
	if m.RoomID == "" {
		return errors.New("vogcluster: RoomBroadcast.room_id is required")
	}
	if m.Sequence <= 0 {
		return fmt.Errorf("vogcluster: RoomBroadcast.sequence must be positive, got %d", m.Sequence)
	}
	return nil
}
