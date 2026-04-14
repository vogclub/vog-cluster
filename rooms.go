package vogcluster

import (
	"encoding/json"
	"errors"
	"fmt"
)

// RoomConfig fully describes a room as the coordinator hands it off to
// a game instance. The instance uses this to load room state on assign
// or migration.
type RoomConfig struct {
	RoomID       string          `json:"room_id"`
	Name         string          `json:"name"`
	GameType     string          `json:"game_type"`
	RatingType   string          `json:"rating_type"`
	Priority     RoomPriority    `json:"priority"`
	Tables       int             `json:"tables"`
	Seats        int             `json:"seats"`
	MaxAccounts  int             `json:"max_accounts"`
	MaxObservers int             `json:"max_observers"`
	Settings     json.RawMessage `json:"settings,omitempty"`
}

// SlotCost returns the room's load weight using the spec formula:
//
//	tables * seats * 5 + max_accounts
//
// The coordinator uses this for bin-packing during balancing.
func (r RoomConfig) SlotCost() int {
	return r.Tables*r.Seats*5 + r.MaxAccounts
}

// Validate reports whether the room configuration is well-formed.
func (r RoomConfig) Validate() error {
	if r.RoomID == "" {
		return errors.New("vogcluster: RoomConfig.room_id is required")
	}
	if r.GameType == "" {
		return errors.New("vogcluster: RoomConfig.game_type is required")
	}
	if r.RatingType == "" {
		return errors.New("vogcluster: RoomConfig.rating_type is required")
	}
	if !r.Priority.Valid() {
		return fmt.Errorf("vogcluster: RoomConfig.priority %q is invalid", r.Priority)
	}
	if r.Tables <= 0 {
		return fmt.Errorf("vogcluster: RoomConfig.tables must be positive, got %d", r.Tables)
	}
	if r.Seats <= 0 {
		return fmt.Errorf("vogcluster: RoomConfig.seats must be positive, got %d", r.Seats)
	}
	if r.MaxAccounts < 0 {
		return fmt.Errorf("vogcluster: RoomConfig.max_accounts must be non-negative, got %d", r.MaxAccounts)
	}
	if r.MaxObservers < 0 {
		return fmt.Errorf("vogcluster: RoomConfig.max_observers must be non-negative, got %d", r.MaxObservers)
	}
	return nil
}

// RoomAssign is sent by the coordinator to a game instance with the
// rooms it should host. The instance loads each room's state and starts
// processing room input traffic.
type RoomAssign struct {
	InstanceID string       `json:"instance_id"`
	Rooms      []RoomConfig `json:"rooms"`
	IssuedBy   string       `json:"issued_by,omitempty"`
}

// Validate reports whether the assignment is well-formed.
func (m RoomAssign) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: RoomAssign.instance_id is required")
	}
	if len(m.Rooms) == 0 {
		return errors.New("vogcluster: RoomAssign.rooms must not be empty")
	}
	for i, r := range m.Rooms {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("vogcluster: RoomAssign.rooms[%d]: %w", i, err)
		}
	}
	return nil
}

// RoomRelease is sent by the coordinator to a game instance to remove
// rooms from its responsibility (drain, migration, or admin action).
type RoomRelease struct {
	InstanceID string   `json:"instance_id"`
	RoomIDs    []string `json:"room_ids"`
	Reason     string   `json:"reason,omitempty"`
}

// Validate reports whether the release is well-formed.
func (m RoomRelease) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: RoomRelease.instance_id is required")
	}
	if len(m.RoomIDs) == 0 {
		return errors.New("vogcluster: RoomRelease.room_ids must not be empty")
	}
	return nil
}

// RoomPrepare is sent during migration: the coordinator asks the
// destination instance to load room state in advance, before the source
// instance is told to release.
type RoomPrepare struct {
	InstanceID  string       `json:"instance_id"`
	Rooms       []RoomConfig `json:"rooms"`
	MigrationID string       `json:"migration_id,omitempty"`
}

// Validate reports whether the prepare is well-formed.
func (m RoomPrepare) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: RoomPrepare.instance_id is required")
	}
	if len(m.Rooms) == 0 {
		return errors.New("vogcluster: RoomPrepare.rooms must not be empty")
	}
	for i, r := range m.Rooms {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("vogcluster: RoomPrepare.rooms[%d]: %w", i, err)
		}
	}
	return nil
}

// RoutingEntry maps a single room to its current owning instance.
// A frozen entry has Frozen=true and may have an empty InstanceID
// when the previous owner is dead and no new owner has been picked yet.
type RoutingEntry struct {
	RoomID     string `json:"room_id"`
	InstanceID string `json:"instance_id,omitempty"`
	Frozen     bool   `json:"frozen,omitempty"`
}

// RoutingUpdate is broadcast by the coordinator to all Lobby instances
// after every change to the routing table. Lobby uses this to refresh
// its in-memory cache.
type RoutingUpdate struct {
	// Version is a monotonically increasing counter. Lobby ignores
	// updates with a version <= the one it has already applied.
	Version int64 `json:"version"`

	// Entries are the rooms whose routing changed. A full snapshot
	// uses Entries containing every known room.
	Entries []RoutingEntry `json:"entries"`

	// Snapshot is true when Entries represents the complete current
	// routing table (used after coordinator failover). When false,
	// Entries are deltas applied on top of the previous version.
	Snapshot bool `json:"snapshot,omitempty"`
}

// Validate reports whether the update is well-formed.
func (m RoutingUpdate) Validate() error {
	if m.Version <= 0 {
		return fmt.Errorf("vogcluster: RoutingUpdate.version must be positive, got %d", m.Version)
	}
	for i, e := range m.Entries {
		if e.RoomID == "" {
			return fmt.Errorf("vogcluster: RoutingUpdate.entries[%d].room_id is required", i)
		}
		if !e.Frozen && e.InstanceID == "" {
			return fmt.Errorf("vogcluster: RoutingUpdate.entries[%d].instance_id is required for non-frozen entries", i)
		}
	}
	return nil
}

// ClusterCommandKind enumerates the admin commands the coordinator
// may issue to a specific game instance.
type ClusterCommandKind string

const (
	// ClusterCommandDrain instructs the instance to begin draining.
	ClusterCommandDrain ClusterCommandKind = "drain"

	// ClusterCommandShutdown instructs the instance to shut down
	// after completing in-flight requests.
	ClusterCommandShutdown ClusterCommandKind = "shutdown"

	// ClusterCommandRefresh asks the instance to reload its room
	// configuration from the coordinator.
	ClusterCommandRefresh ClusterCommandKind = "refresh"
)

// Valid reports whether k is one of the defined commands.
func (k ClusterCommandKind) Valid() bool {
	switch k {
	case ClusterCommandDrain, ClusterCommandShutdown, ClusterCommandRefresh:
		return true
	}
	return false
}

// ClusterCommand is sent on SubjectClusterCommand(instance_id) for
// admin-initiated actions. The target instance ID is encoded in the
// subject; the message body carries metadata.
type ClusterCommand struct {
	Command  ClusterCommandKind `json:"command"`
	IssuedBy string             `json:"issued_by,omitempty"`
	Reason   string             `json:"reason,omitempty"`
}

// Validate reports whether the command is well-formed.
func (m ClusterCommand) Validate() error {
	if !m.Command.Valid() {
		return fmt.Errorf("vogcluster: ClusterCommand.command %q is invalid", m.Command)
	}
	return nil
}
