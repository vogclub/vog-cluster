package vogcluster

import (
	"errors"
	"fmt"
	"time"
)

// RoomPrepareResponse is the reply to a RoomPrepare request on
// SubjectClusterRoomsPrepare(instance_id). The target instance sends
// this back via NATS request-reply after it has either warmed state for
// the room(s) or decided it cannot accept them.
type RoomPrepareResponse struct {
	MigrationID string `json:"migration_id"`
	Accepted    bool   `json:"accepted"`
	// Reason is required when Accepted is false. Free-form machine
	// code (e.g. "out_of_capacity", "unknown_room", "tournament_mismatch").
	Reason string `json:"reason,omitempty"`
}

// Validate reports whether the response is well-formed.
func (m RoomPrepareResponse) Validate() error {
	if m.MigrationID == "" {
		return errors.New("vogcluster: RoomPrepareResponse.migration_id is required")
	}
	if !m.Accepted && m.Reason == "" {
		return errors.New("vogcluster: RoomPrepareResponse.reason is required when not accepted")
	}
	return nil
}

// RoomReady is published by a target instance on
// SubjectClusterRoomReady(server_id, room_id) after it has successfully
// committed a migrated room (i.e. it is now serving traffic for that
// room). The coordinator uses this as the terminal signal for a
// migration and updates the routing table.
type RoomReady struct {
	MigrationID string    `json:"migration_id"`
	InstanceID  string    `json:"instance_id"`
	RoomID      string    `json:"room_id"`
	ReadyAt     time.Time `json:"ready_at,omitempty"`
}

// Validate reports whether the message is well-formed.
func (m RoomReady) Validate() error {
	if m.MigrationID == "" {
		return errors.New("vogcluster: RoomReady.migration_id is required")
	}
	if m.InstanceID == "" {
		return errors.New("vogcluster: RoomReady.instance_id is required")
	}
	if m.RoomID == "" {
		return errors.New("vogcluster: RoomReady.room_id is required")
	}
	return nil
}

// InstanceRegisterReply is the reply the coordinator sends back on
// SubjectClusterGameRegister request-reply. When Accepted is false, the
// instance must exit with a non-zero code (the coordinator has rejected
// it — typically claimer_hash collision).
type InstanceRegisterReply struct {
	Accepted   bool   `json:"accepted"`
	InstanceID string `json:"instance_id"`
	// Reason is required when Accepted is false. Machine code
	// (e.g. "claimed_by_another", "invalid_payload").
	Reason string `json:"reason,omitempty"`
}

// Validate reports whether the reply is well-formed.
func (m InstanceRegisterReply) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: InstanceRegisterReply.instance_id is required")
	}
	if !m.Accepted && m.Reason == "" {
		return errors.New("vogcluster: InstanceRegisterReply.reason is required when not accepted")
	}
	return nil
}

// AdminMigrateRequest is sent by a human operator (via vog-spaces HTTP
// admin → NATS request-reply) to move a single room from its current
// instance to a specific target instance. Published on
// SubjectClusterAdminMigrate.
type AdminMigrateRequest struct {
	ServerID   int    `json:"server_id"`
	RoomID     int    `json:"room_id"`
	ToInstance string `json:"to_instance"`
	IssuedBy   string `json:"issued_by"`
	Reason     string `json:"reason,omitempty"`
}

// Validate reports whether the request is well-formed.
func (m AdminMigrateRequest) Validate() error {
	if m.ServerID <= 0 {
		return fmt.Errorf("vogcluster: AdminMigrateRequest.server_id must be positive, got %d", m.ServerID)
	}
	if m.RoomID <= 0 {
		return fmt.Errorf("vogcluster: AdminMigrateRequest.room_id must be positive, got %d", m.RoomID)
	}
	if m.ToInstance == "" {
		return errors.New("vogcluster: AdminMigrateRequest.to_instance is required")
	}
	if m.IssuedBy == "" {
		return errors.New("vogcluster: AdminMigrateRequest.issued_by is required")
	}
	return nil
}

// AdminMigrateResponse is the reply to AdminMigrateRequest.
type AdminMigrateResponse struct {
	Status      string `json:"status"` // "ok" | "rollback" | "rejected" | "rate_limited" | "grace_period" | "unknown_room" | "commit_stuck"
	MigrationID string `json:"migration_id,omitempty"`
	Detail      string `json:"detail,omitempty"`
}

// AdminDrainRequest instructs the coordinator to drain all rooms off
// the given instance. Published on SubjectClusterAdminDrain.
type AdminDrainRequest struct {
	InstanceID string `json:"instance_id"`
	IssuedBy   string `json:"issued_by"`
	Reason     string `json:"reason,omitempty"`
}

// Validate reports whether the request is well-formed.
func (m AdminDrainRequest) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: AdminDrainRequest.instance_id is required")
	}
	if m.IssuedBy == "" {
		return errors.New("vogcluster: AdminDrainRequest.issued_by is required")
	}
	return nil
}

// AdminDrainResponse is the reply to AdminDrainRequest.
type AdminDrainResponse struct {
	Status     string `json:"status"` // "ok" | "rejected" | "partial"
	Migrations int    `json:"migrations"`
	Detail     string `json:"detail,omitempty"`
}

// AdminRebalanceRequest asks the coordinator to run a rebalance pass
// immediately, bypassing the rate limiter but honoring the grace
// period. Published on SubjectClusterAdminRebalance.
type AdminRebalanceRequest struct {
	IssuedBy string `json:"issued_by"`
	DryRun   bool   `json:"dry_run,omitempty"`
}

// Validate reports whether the request is well-formed.
func (m AdminRebalanceRequest) Validate() error {
	if m.IssuedBy == "" {
		return errors.New("vogcluster: AdminRebalanceRequest.issued_by is required")
	}
	return nil
}

// AdminRebalanceResponse is the reply to AdminRebalanceRequest.
type AdminRebalanceResponse struct {
	Status     string `json:"status"` // "ok" | "noop" | "rejected"
	Migrations int    `json:"migrations"`
	Detail     string `json:"detail,omitempty"`
}

// InstanceStatusEvent is published by the coordinator on
// SubjectClusterInstanceStatus(instance_id) whenever an instance's
// status changes. Consumed by lobby, admin UI, and the audit journal.
type InstanceStatusEvent struct {
	InstanceID string         `json:"instance_id"`
	Status     InstanceStatus `json:"status"`
	ChangedAt  time.Time      `json:"changed_at,omitempty"`
	Reason     string         `json:"reason,omitempty"`
}

// Validate reports whether the event is well-formed.
func (m InstanceStatusEvent) Validate() error {
	if m.InstanceID == "" {
		return errors.New("vogcluster: InstanceStatusEvent.instance_id is required")
	}
	if !m.Status.Valid() {
		return fmt.Errorf("vogcluster: InstanceStatusEvent.status %q is invalid", m.Status)
	}
	return nil
}
