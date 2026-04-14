package vogcluster

import (
	"encoding/json"
	"reflect"
	"testing"
)

func sampleRoomConfig() RoomConfig {
	return RoomConfig{
		RoomID:       "br-rapid-1",
		Name:         "Backgammon Rapid #1",
		GameType:     "backgammon",
		RatingType:   "br-rapid",
		Priority:     RoomPriorityNormal,
		Tables:       20,
		Seats:        2,
		MaxAccounts:  100,
		MaxObservers: 200,
		Settings:     json.RawMessage(`{"clock":"5+3","variant":"standard"}`),
	}
}

func TestRoomConfigSlotCost(t *testing.T) {
	c := sampleRoomConfig()
	got := c.SlotCost()
	want := 20*2*5 + 100 // 300
	if got != want {
		t.Errorf("SlotCost() = %d, want %d", got, want)
	}
}

func TestRoomConfigValidate(t *testing.T) {
	valid := sampleRoomConfig()
	if err := valid.Validate(); err != nil {
		t.Errorf("valid room returned error: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*RoomConfig)
	}{
		{"empty room_id", func(c *RoomConfig) { c.RoomID = "" }},
		{"empty game_type", func(c *RoomConfig) { c.GameType = "" }},
		{"empty rating_type", func(c *RoomConfig) { c.RatingType = "" }},
		{"invalid priority", func(c *RoomConfig) { c.Priority = "urgent" }},
		{"zero tables", func(c *RoomConfig) { c.Tables = 0 }},
		{"zero seats", func(c *RoomConfig) { c.Seats = 0 }},
		{"negative max_accounts", func(c *RoomConfig) { c.MaxAccounts = -1 }},
		{"negative max_observers", func(c *RoomConfig) { c.MaxObservers = -1 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := sampleRoomConfig()
			tc.mutate(&c)
			if err := c.Validate(); err == nil {
				t.Errorf("expected validation error, got nil")
			}
		})
	}
}

func TestRoomAssignRoundtrip(t *testing.T) {
	original := RoomAssign{
		InstanceID: "game-7",
		Rooms:      []RoomConfig{sampleRoomConfig()},
		IssuedBy:   "coordinator-primary",
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded RoomAssign
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestRoomAssignValidate(t *testing.T) {
	valid := RoomAssign{InstanceID: "game-7", Rooms: []RoomConfig{sampleRoomConfig()}}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid assign returned error: %v", err)
	}

	if err := (RoomAssign{Rooms: []RoomConfig{sampleRoomConfig()}}).Validate(); err == nil {
		t.Errorf("missing instance_id should fail")
	}
	if err := (RoomAssign{InstanceID: "game-7"}).Validate(); err == nil {
		t.Errorf("empty rooms should fail")
	}
	bad := sampleRoomConfig()
	bad.Tables = 0
	if err := (RoomAssign{InstanceID: "game-7", Rooms: []RoomConfig{bad}}).Validate(); err == nil {
		t.Errorf("invalid room should fail")
	}
}

func TestRoomReleaseValidate(t *testing.T) {
	if err := (RoomRelease{InstanceID: "game-7", RoomIDs: []string{"a"}}).Validate(); err != nil {
		t.Errorf("valid release returned error: %v", err)
	}
	if err := (RoomRelease{RoomIDs: []string{"a"}}).Validate(); err == nil {
		t.Errorf("missing instance_id should fail")
	}
	if err := (RoomRelease{InstanceID: "game-7"}).Validate(); err == nil {
		t.Errorf("empty room_ids should fail")
	}
}

func TestRoomPrepareValidate(t *testing.T) {
	valid := RoomPrepare{InstanceID: "game-7", Rooms: []RoomConfig{sampleRoomConfig()}}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid prepare returned error: %v", err)
	}
	if err := (RoomPrepare{Rooms: []RoomConfig{sampleRoomConfig()}}).Validate(); err == nil {
		t.Errorf("missing instance_id should fail")
	}
}

func TestRoutingUpdateRoundtrip(t *testing.T) {
	original := RoutingUpdate{
		Version: 42,
		Entries: []RoutingEntry{
			{RoomID: "br-rapid-1", InstanceID: "game-7", Frozen: false},
			{RoomID: "ch-blitz-3", InstanceID: "game-2", Frozen: true},
		},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded RoutingUpdate
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestRoutingUpdateValidate(t *testing.T) {
	valid := RoutingUpdate{Version: 1, Entries: []RoutingEntry{{RoomID: "r", InstanceID: "i"}}}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid update returned error: %v", err)
	}
	if err := (RoutingUpdate{Version: 0, Entries: []RoutingEntry{{RoomID: "r", InstanceID: "i"}}}).Validate(); err == nil {
		t.Errorf("zero version should fail")
	}
	if err := (RoutingUpdate{Version: 1, Entries: []RoutingEntry{{InstanceID: "i"}}}).Validate(); err == nil {
		t.Errorf("missing room_id should fail")
	}
	// Frozen entries are allowed to have empty instance_id.
	if err := (RoutingUpdate{Version: 1, Entries: []RoutingEntry{{RoomID: "r", Frozen: true}}}).Validate(); err != nil {
		t.Errorf("frozen entry without instance should be valid: %v", err)
	}
	// Non-frozen entries require instance_id.
	if err := (RoutingUpdate{Version: 1, Entries: []RoutingEntry{{RoomID: "r"}}}).Validate(); err == nil {
		t.Errorf("non-frozen entry without instance should fail")
	}
}

func TestClusterCommandValidate(t *testing.T) {
	valid := ClusterCommand{Command: ClusterCommandDrain, IssuedBy: "admin"}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid command returned error: %v", err)
	}
	if err := (ClusterCommand{IssuedBy: "admin"}).Validate(); err == nil {
		t.Errorf("missing command should fail")
	}
	if err := (ClusterCommand{Command: "destroy", IssuedBy: "admin"}).Validate(); err == nil {
		t.Errorf("unknown command should fail")
	}
}
