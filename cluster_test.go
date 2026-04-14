package vogcluster

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestGameRegisterRoundtrip(t *testing.T) {
	original := GameRegister{
		InstanceID:     "game-7",
		Capacity:       5000,
		Version:        "1.2.3",
		Address:        "10.0.0.42:9001",
		TournamentOnly: false,
		TournamentID:   "",
		RegisteredAt:   time.Date(2026, 4, 13, 10, 30, 0, 0, time.UTC),
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded GameRegister
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestGameRegisterValidate(t *testing.T) {
	tests := []struct {
		name    string
		msg     GameRegister
		wantErr bool
	}{
		{
			name: "valid",
			msg: GameRegister{
				InstanceID: "game-7", Capacity: 1000, Version: "1.0.0", Address: "host:9000",
			},
			wantErr: false,
		},
		{
			name:    "missing instance_id",
			msg:     GameRegister{Capacity: 1000, Version: "1.0.0", Address: "host:9000"},
			wantErr: true,
		},
		{
			name:    "zero capacity",
			msg:     GameRegister{InstanceID: "game-7", Version: "1.0.0", Address: "host:9000"},
			wantErr: true,
		},
		{
			name:    "negative capacity",
			msg:     GameRegister{InstanceID: "game-7", Capacity: -1, Version: "1.0.0", Address: "host:9000"},
			wantErr: true,
		},
		{
			name:    "missing version",
			msg:     GameRegister{InstanceID: "game-7", Capacity: 1000, Address: "host:9000"},
			wantErr: true,
		},
		{
			name:    "missing address",
			msg:     GameRegister{InstanceID: "game-7", Capacity: 1000, Version: "1.0.0"},
			wantErr: true,
		},
		{
			name: "tournament_only without tournament_id",
			msg: GameRegister{
				InstanceID: "game-7", Capacity: 1000, Version: "1.0.0", Address: "host:9000",
				TournamentOnly: true,
			},
			wantErr: true,
		},
		{
			name: "tournament_only with tournament_id",
			msg: GameRegister{
				InstanceID: "game-7", Capacity: 1000, Version: "1.0.0", Address: "host:9000",
				TournamentOnly: true, TournamentID: "tourney-42",
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestGameHeartbeatRoundtrip(t *testing.T) {
	original := GameHeartbeat{
		InstanceID: "game-7",
		Status:     InstanceStatusActive,
		SlotsUsed:  1234,
		SlotsTotal: 5000,
		Rooms: []RoomLoad{
			{RoomID: "br-rapid-1", Tables: 20, Players: 87, SlotsUsed: 287},
			{RoomID: "ch-blitz-3", Tables: 15, Players: 42, SlotsUsed: 192},
		},
		ObservedAt: time.Date(2026, 4, 13, 10, 31, 5, 0, time.UTC),
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded GameHeartbeat
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestGameHeartbeatValidate(t *testing.T) {
	valid := GameHeartbeat{
		InstanceID: "game-7", Status: InstanceStatusActive,
		SlotsUsed: 100, SlotsTotal: 1000,
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid heartbeat returned error: %v", err)
	}

	invalid := []struct {
		name string
		msg  GameHeartbeat
	}{
		{"empty id", GameHeartbeat{Status: InstanceStatusActive, SlotsTotal: 1000}},
		{"invalid status", GameHeartbeat{InstanceID: "g", Status: "running", SlotsTotal: 1000}},
		{"slots used > total", GameHeartbeat{InstanceID: "g", Status: InstanceStatusActive, SlotsUsed: 2000, SlotsTotal: 1000}},
		{"negative slots used", GameHeartbeat{InstanceID: "g", Status: InstanceStatusActive, SlotsUsed: -1, SlotsTotal: 1000}},
		{"zero total slots", GameHeartbeat{InstanceID: "g", Status: InstanceStatusActive}},
		{"bad nested room", GameHeartbeat{
			InstanceID: "g", Status: InstanceStatusActive, SlotsTotal: 1000,
			Rooms: []RoomLoad{{RoomID: "", Tables: 1, Players: 1, SlotsUsed: 1}},
		}},
		{"negative room tables", GameHeartbeat{
			InstanceID: "g", Status: InstanceStatusActive, SlotsTotal: 1000,
			Rooms: []RoomLoad{{RoomID: "r", Tables: -1}},
		}},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.msg.Validate(); err == nil {
				t.Errorf("expected validation error, got nil")
			}
		})
	}
}

func TestRoomLoadValidate(t *testing.T) {
	valid := RoomLoad{RoomID: "r", Tables: 10, Players: 50, SlotsUsed: 150}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid room load returned error: %v", err)
	}

	tests := []struct {
		name string
		load RoomLoad
	}{
		{"empty room_id", RoomLoad{Tables: 1, Players: 1, SlotsUsed: 1}},
		{"negative tables", RoomLoad{RoomID: "r", Tables: -1}},
		{"negative players", RoomLoad{RoomID: "r", Players: -1}},
		{"negative slots_used", RoomLoad{RoomID: "r", SlotsUsed: -1}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.load.Validate(); err == nil {
				t.Errorf("expected validation error, got nil")
			}
		})
	}
}

func TestGameDrainAndDeregisterValidate(t *testing.T) {
	if err := (GameDrain{InstanceID: "game-7"}).Validate(); err != nil {
		t.Errorf("valid drain returned error: %v", err)
	}
	if err := (GameDrain{}).Validate(); err == nil {
		t.Errorf("empty drain should fail validation")
	}
	if err := (GameDeregister{InstanceID: "game-7"}).Validate(); err != nil {
		t.Errorf("valid deregister returned error: %v", err)
	}
	if err := (GameDeregister{}).Validate(); err == nil {
		t.Errorf("empty deregister should fail validation")
	}
}
