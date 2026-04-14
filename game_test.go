package vogcluster

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestRoomInputRoundtrip(t *testing.T) {
	original := RoomInput{
		Action:          "table.move",
		UserID:          "user-123",
		LobbyInstanceID: "lobby-2",
		ConnectionID:    "conn-abc",
		ClientNonce:     "nonce-xyz",
		ReceivedAt:      time.Date(2026, 4, 13, 11, 0, 0, 0, time.UTC),
		Payload:         json.RawMessage(`{"table_id":7,"move":"24/18 13/11"}`),
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded RoomInput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestRoomInputValidate(t *testing.T) {
	valid := RoomInput{
		Action:          "chat",
		UserID:          "user-123",
		LobbyInstanceID: "lobby-2",
		ConnectionID:    "conn-abc",
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid input returned error: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*RoomInput)
	}{
		{"empty action", func(m *RoomInput) { m.Action = "" }},
		{"empty user_id", func(m *RoomInput) { m.UserID = "" }},
		{"empty lobby_instance_id", func(m *RoomInput) { m.LobbyInstanceID = "" }},
		{"empty connection_id", func(m *RoomInput) { m.ConnectionID = "" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := valid
			tc.mutate(&m)
			if err := m.Validate(); err == nil {
				t.Errorf("expected validation error, got nil")
			}
		})
	}
}

func TestRoomBroadcastRoundtrip(t *testing.T) {
	original := RoomBroadcast{
		Event:        "table.move",
		RoomID:       "br-rapid-1",
		Sequence:     1042,
		EmittedAt:    time.Date(2026, 4, 13, 11, 0, 0, 0, time.UTC),
		Payload:      json.RawMessage(`{"table_id":7,"by":"user-123","move":"24/18 13/11"}`),
		ExcludeUsers: []string{"user-123"},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded RoomBroadcast
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestRoomBroadcastValidate(t *testing.T) {
	valid := RoomBroadcast{Event: "table.move", RoomID: "br-rapid-1", Sequence: 1}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid broadcast returned error: %v", err)
	}
	if err := (RoomBroadcast{RoomID: "br-rapid-1", Sequence: 1}).Validate(); err == nil {
		t.Errorf("missing event should fail")
	}
	if err := (RoomBroadcast{Event: "x", Sequence: 1}).Validate(); err == nil {
		t.Errorf("missing room_id should fail")
	}
	if err := (RoomBroadcast{Event: "x", RoomID: "r"}).Validate(); err == nil {
		t.Errorf("zero sequence should fail")
	}
}
