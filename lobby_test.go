package vogcluster

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestLobbyOutputRoundtrip(t *testing.T) {
	original := LobbyOutput{
		Event:         "seek.created",
		TargetUserIDs: []string{"user-123"},
		TargetConnIDs: []string{"conn-abc"},
		ClientNonce:   "nonce-xyz",
		EmittedAt:     time.Date(2026, 4, 13, 11, 0, 0, 0, time.UTC),
		Payload:       json.RawMessage(`{"seek_id":42,"rating":1500}`),
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded LobbyOutput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestLobbyOutputValidate(t *testing.T) {
	// Either TargetUserIDs or TargetConnIDs must be non-empty.
	valid1 := LobbyOutput{Event: "x", TargetUserIDs: []string{"u"}}
	if err := valid1.Validate(); err != nil {
		t.Errorf("valid1 returned error: %v", err)
	}
	valid2 := LobbyOutput{Event: "x", TargetConnIDs: []string{"c"}}
	if err := valid2.Validate(); err != nil {
		t.Errorf("valid2 returned error: %v", err)
	}

	if err := (LobbyOutput{TargetUserIDs: []string{"u"}}).Validate(); err == nil {
		t.Errorf("missing event should fail")
	}
	if err := (LobbyOutput{Event: "x"}).Validate(); err == nil {
		t.Errorf("no targets should fail")
	}
}

func TestSessionEventValidate(t *testing.T) {
	valid := &SessionEvent{UserID: "u1", RoomID: "r1", Action: "enter"}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid SessionEvent returned error: %v", err)
	}

	if err := (&SessionEvent{RoomID: "r1", Action: "enter"}).Validate(); err == nil {
		t.Error("missing UserID should fail")
	}
	if err := (&SessionEvent{UserID: "u1", Action: "enter"}).Validate(); err == nil {
		t.Error("missing RoomID should fail")
	}
	if err := (&SessionEvent{UserID: "u1", RoomID: "r1"}).Validate(); err == nil {
		t.Error("missing Action should fail")
	}
}

func TestSessionEventRoundtrip(t *testing.T) {
	original := &SessionEvent{
		UserID:     "user-42",
		RoomID:     "room-7",
		Action:     "enter",
		IP:         "192.168.1.1",
		ClientType: "web",
	}
	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	var decoded SessionEvent
	if err := Decode(data, &decoded); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if !reflect.DeepEqual(*original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", *original, decoded)
	}
}

func TestSessionEventEncodeRejectsInvalid(t *testing.T) {
	invalid := &SessionEvent{RoomID: "r1", Action: "enter"} // missing UserID
	if _, err := Encode(invalid); err == nil {
		t.Error("Encode should reject invalid SessionEvent")
	}
}
