package vogcluster

// InstanceStatus represents the lifecycle state of a vog-game instance
// as tracked by the coordinator (vog-spaces).
type InstanceStatus string

const (
	// InstanceStatusPending means the instance has registered but
	// has not yet received any room assignments.
	InstanceStatusPending InstanceStatus = "pending"

	// InstanceStatusActive means the instance is operating normally
	// and may accept new room assignments.
	InstanceStatusActive InstanceStatus = "active"

	// InstanceStatusDraining means the instance is being removed
	// from rotation. Existing rooms are finishing or migrating;
	// no new rooms will be assigned.
	InstanceStatusDraining InstanceStatus = "draining"

	// InstanceStatusFrozen means the coordinator has stopped receiving
	// heartbeats. Rooms on this instance are frozen pending the freeze
	// timeout.
	InstanceStatusFrozen InstanceStatus = "frozen"

	// InstanceStatusDead means the freeze timeout has expired and the
	// instance's rooms have been redistributed.
	InstanceStatusDead InstanceStatus = "dead"
)

// Valid reports whether s is one of the defined InstanceStatus values.
func (s InstanceStatus) Valid() bool {
	switch s {
	case InstanceStatusPending,
		InstanceStatusActive,
		InstanceStatusDraining,
		InstanceStatusFrozen,
		InstanceStatusDead:
		return true
	}
	return false
}

// RoomPriority controls how the coordinator handles a room under
// resource pressure. Lower-priority rooms are shed first.
type RoomPriority string

const (
	// RoomPriorityCritical rooms are always assigned, even when
	// total capacity is insufficient.
	RoomPriorityCritical RoomPriority = "critical"

	// RoomPriorityHigh rooms are assigned after critical.
	RoomPriorityHigh RoomPriority = "high"

	// RoomPriorityNormal is the default room priority.
	RoomPriorityNormal RoomPriority = "normal"

	// RoomPriorityLow rooms are skipped first when capacity
	// is insufficient.
	RoomPriorityLow RoomPriority = "low"
)

// Valid reports whether p is one of the defined RoomPriority values.
func (p RoomPriority) Valid() bool {
	switch p {
	case RoomPriorityCritical,
		RoomPriorityHigh,
		RoomPriorityNormal,
		RoomPriorityLow:
		return true
	}
	return false
}

// Rank returns an ordering value where lower means "assigned first".
// Critical = 0, High = 1, Normal = 2, Low = 3. Invalid priorities
// return math.MaxInt to sort last.
func (p RoomPriority) Rank() int {
	switch p {
	case RoomPriorityCritical:
		return 0
	case RoomPriorityHigh:
		return 1
	case RoomPriorityNormal:
		return 2
	case RoomPriorityLow:
		return 3
	}
	return 1<<31 - 1
}

// Transport identifies whether a NATS subject requires JetStream
// (durable, replayable) or plain Core NATS (best-effort, low-latency).
type Transport int

const (
	// TransportCore uses plain Core NATS publish/subscribe.
	// Suitable for low-latency, loss-tolerant traffic.
	TransportCore Transport = iota

	// TransportJetStream requires the subject to be backed by a
	// JetStream stream. Used for messages that must not be lost.
	TransportJetStream
)

// String returns the lowercase name of the transport.
func (t Transport) String() string {
	switch t {
	case TransportCore:
		return "core"
	case TransportJetStream:
		return "jetstream"
	}
	return "unknown"
}
