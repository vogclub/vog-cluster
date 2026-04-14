package vogcluster

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestEncodeProducesCompactJSON(t *testing.T) {
	msg := GameRegister{
		InstanceID: "game-7", Capacity: 1000, Version: "1.0", Address: "h:9",
	}
	data, err := Encode(msg)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if bytes.Contains(data, []byte("  ")) || bytes.Contains(data, []byte("\n")) {
		t.Errorf("Encode should produce compact JSON, got %q", data)
	}
	if !json.Valid(data) {
		t.Errorf("Encode produced invalid JSON: %q", data)
	}
}

func TestEncodeRejectsInvalid(t *testing.T) {
	bad := GameRegister{InstanceID: ""} // missing required fields
	if _, err := Encode(bad); err == nil {
		t.Errorf("Encode should fail validation, got nil error")
	}
}

func TestDecodeIntoStructAndValidates(t *testing.T) {
	good, err := Encode(GameRegister{
		InstanceID: "game-7", Capacity: 1000, Version: "1.0", Address: "h:9",
	})
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	var msg GameRegister
	if err := Decode(good, &msg); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if msg.InstanceID != "game-7" {
		t.Errorf("decoded wrong instance: %+v", msg)
	}
}

func TestDecodeRejectsMalformed(t *testing.T) {
	var msg GameRegister
	if err := Decode([]byte("{not json"), &msg); err == nil {
		t.Errorf("Decode should fail on malformed JSON")
	}
}

func TestDecodeRejectsInvalidPayload(t *testing.T) {
	// Valid JSON but fails Validate (missing required fields).
	var msg GameRegister
	err := Decode([]byte(`{"instance_id":""}`), &msg)
	if err == nil {
		t.Errorf("Decode should fail validation")
	}
	if !strings.Contains(err.Error(), "instance_id") {
		t.Errorf("error should mention instance_id, got %v", err)
	}
}

func TestEncodeAcceptsValidator(t *testing.T) {
	// Compile-time check: every message type with Validate satisfies
	// the Validator interface used by Encode/Decode.
	var _ Validator = GameRegister{}
	var _ Validator = GameHeartbeat{}
	var _ Validator = GameDrain{}
	var _ Validator = GameDeregister{}
	var _ Validator = RoomConfig{}
	var _ Validator = RoomAssign{}
	var _ Validator = RoomRelease{}
	var _ Validator = RoomPrepare{}
	var _ Validator = RoutingUpdate{}
	var _ Validator = ClusterCommand{}
	var _ Validator = RoomInput{}
	var _ Validator = RoomBroadcast{}
	var _ Validator = LobbyOutput{}
	var _ Validator = StatsOnline{}
	var _ Validator = RatingUpdated{}
}
