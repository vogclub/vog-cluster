package vogcluster

import (
	"errors"
	"fmt"
	"time"
)

// GameRegister is published by a game instance on startup to subject
// SubjectClusterGameRegister. The coordinator (vog-spaces) adds the
// instance to its registry with status `pending` and begins assigning rooms.
type GameRegister struct {
	// InstanceID uniquely identifies this game instance for the
	// lifetime of its run. Format is service-defined; recommended
	// "{hostname}-{pid}-{startup_unix}".
	InstanceID string `json:"instance_id"`

	// Capacity is the total number of slots the instance is willing
	// to host. See spec slot formula: tables*seats*5 + max_accounts.
	Capacity int `json:"capacity"`

	// Version is the build/release version of the game binary.
	Version string `json:"version"`

	// Address is the network address (host:port) the instance is
	// reachable on for direct connections, if any.
	Address string `json:"address"`

	// TournamentOnly marks the instance as accepting only tournament
	// rooms. The coordinator excludes such instances from regular pool.
	TournamentOnly bool `json:"tournament_only,omitempty"`

	// TournamentID is required when TournamentOnly is true and links
	// the instance to a specific tournament managed by vog-tourneys.
	TournamentID string `json:"tournament_id,omitempty"`

	// RegisteredAt is the wall-clock time the instance started
	// registration, in UTC.
	RegisteredAt time.Time `json:"registered_at,omitempty"`
}

// Validate reports whether the message is well-formed.
func (m GameRegister) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: GameRegister.instance_id is required")
	}
	if m.Capacity <= 0 {
		return fmt.Errorf("vogcluster: GameRegister.capacity must be positive, got %d", m.Capacity)
	}
	if m.Version == "" {
		return errors.New("vogcluster: GameRegister.version is required")
	}
	if m.Address == "" {
		return errors.New("vogcluster: GameRegister.address is required")
	}
	if m.TournamentOnly && m.TournamentID == "" {
		return errors.New("vogcluster: GameRegister.tournament_id is required when tournament_only is true")
	}
	return nil
}

// RoomLoad reports per-room load metrics inside a heartbeat.
type RoomLoad struct {
	RoomID    string `json:"room_id"`
	Tables    int    `json:"tables"`
	Players   int    `json:"players"`
	SlotsUsed int    `json:"slots_used"`
}

// Validate reports whether the per-room load report is well-formed.
func (r RoomLoad) Validate() error {
	if r.RoomID == "" {
		return errors.New("vogcluster: RoomLoad.room_id is required")
	}
	if r.Tables < 0 {
		return fmt.Errorf("vogcluster: RoomLoad.tables must be non-negative, got %d", r.Tables)
	}
	if r.Players < 0 {
		return fmt.Errorf("vogcluster: RoomLoad.players must be non-negative, got %d", r.Players)
	}
	if r.SlotsUsed < 0 {
		return fmt.Errorf("vogcluster: RoomLoad.slots_used must be non-negative, got %d", r.SlotsUsed)
	}
	return nil
}

// GameHeartbeat is published periodically (~5s) by each active game
// instance on subject SubjectClusterGameHeartbeat(instance_id).
type GameHeartbeat struct {
	InstanceID string         `json:"instance_id"`
	Status     InstanceStatus `json:"status"`
	SlotsUsed  int            `json:"slots_used"`
	SlotsTotal int            `json:"slots_total"`
	Rooms      []RoomLoad     `json:"rooms,omitempty"`
	ObservedAt time.Time      `json:"observed_at,omitempty"`
}

// Validate reports whether the heartbeat is well-formed.
func (m GameHeartbeat) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: GameHeartbeat.instance_id is required")
	}
	if !m.Status.Valid() {
		return fmt.Errorf("vogcluster: GameHeartbeat.status %q is invalid", m.Status)
	}
	if m.SlotsTotal <= 0 {
		return fmt.Errorf("vogcluster: GameHeartbeat.slots_total must be positive, got %d", m.SlotsTotal)
	}
	if m.SlotsUsed < 0 {
		return fmt.Errorf("vogcluster: GameHeartbeat.slots_used must be non-negative, got %d", m.SlotsUsed)
	}
	if m.SlotsUsed > m.SlotsTotal {
		return fmt.Errorf("vogcluster: GameHeartbeat.slots_used (%d) exceeds slots_total (%d)", m.SlotsUsed, m.SlotsTotal)
	}
	for i, r := range m.Rooms {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("vogcluster: GameHeartbeat.rooms[%d]: %w", i, err)
		}
	}
	return nil
}

// GameDrain is published by a game instance to announce it is starting
// graceful shutdown. The coordinator transitions it to draining state.
type GameDrain struct {
	InstanceID string    `json:"instance_id"`
	Reason     string    `json:"reason,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
}

// Validate reports whether the message is well-formed.
func (m GameDrain) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: GameDrain.instance_id is required")
	}
	return nil
}

// GameDeregister is published by a game instance immediately before
// shutdown so the coordinator can redistribute its rooms.
type GameDeregister struct {
	InstanceID string    `json:"instance_id"`
	StoppedAt  time.Time `json:"stopped_at,omitempty"`
}

// Validate reports whether the message is well-formed.
func (m GameDeregister) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: GameDeregister.instance_id is required")
	}
	return nil
}
