package vogcluster

import (
	"testing"
)

func TestInstanceStatusValid(t *testing.T) {
	tests := []struct {
		name   string
		status InstanceStatus
		want   bool
	}{
		{"pending", InstanceStatusPending, true},
		{"active", InstanceStatusActive, true},
		{"draining", InstanceStatusDraining, true},
		{"frozen", InstanceStatusFrozen, true},
		{"dead", InstanceStatusDead, true},
		{"empty", InstanceStatus(""), false},
		{"unknown", InstanceStatus("running"), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.status.Valid(); got != tc.want {
				t.Errorf("InstanceStatus(%q).Valid() = %v, want %v", tc.status, got, tc.want)
			}
		})
	}
}

func TestRoomPriorityValid(t *testing.T) {
	tests := []struct {
		name     string
		priority RoomPriority
		want     bool
	}{
		{"critical", RoomPriorityCritical, true},
		{"high", RoomPriorityHigh, true},
		{"normal", RoomPriorityNormal, true},
		{"low", RoomPriorityLow, true},
		{"empty", RoomPriority(""), false},
		{"unknown", RoomPriority("urgent"), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.priority.Valid(); got != tc.want {
				t.Errorf("RoomPriority(%q).Valid() = %v, want %v", tc.priority, got, tc.want)
			}
		})
	}
}

func TestRoomPriorityRank(t *testing.T) {
	// Lower rank = assigned first
	if RoomPriorityCritical.Rank() >= RoomPriorityHigh.Rank() {
		t.Errorf("critical should have lower rank than high")
	}
	if RoomPriorityHigh.Rank() >= RoomPriorityNormal.Rank() {
		t.Errorf("high should have lower rank than normal")
	}
	if RoomPriorityNormal.Rank() >= RoomPriorityLow.Rank() {
		t.Errorf("normal should have lower rank than low")
	}
}

func TestTransportString(t *testing.T) {
	if TransportCore.String() != "core" {
		t.Errorf("TransportCore.String() = %q, want %q", TransportCore.String(), "core")
	}
	if TransportJetStream.String() != "jetstream" {
		t.Errorf("TransportJetStream.String() = %q, want %q", TransportJetStream.String(), "jetstream")
	}
}
