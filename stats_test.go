package vogcluster

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestStatsOnlineRoundtrip(t *testing.T) {
	original := StatsOnline{
		Total:      12345,
		ObservedAt: time.Date(2026, 4, 13, 11, 0, 0, 0, time.UTC),
		Servers: []StatsServerCount{
			{GameType: "backgammon", Count: 5000},
			{GameType: "chess", Count: 3500},
		},
		Rooms: []StatsRoomCount{
			{RoomID: "br-rapid-1", Players: 87, TablesActive: 12},
			{RoomID: "ch-blitz-3", Players: 42, TablesActive: 7},
		},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded StatsOnline
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestStatsOnlineValidate(t *testing.T) {
	if err := (StatsOnline{Total: 100}).Validate(); err != nil {
		t.Errorf("valid stats returned error: %v", err)
	}
	if err := (StatsOnline{Total: -1}).Validate(); err == nil {
		t.Errorf("negative total should fail")
	}
	if err := (StatsOnline{
		Total:   100,
		Servers: []StatsServerCount{{GameType: "", Count: 1}},
	}).Validate(); err == nil {
		t.Errorf("empty game_type should fail")
	}
	if err := (StatsOnline{
		Total: 100,
		Rooms: []StatsRoomCount{{RoomID: "", Players: 1}},
	}).Validate(); err == nil {
		t.Errorf("empty room_id should fail")
	}
}
