package vogcluster

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestRatingUpdatedRoundtrip(t *testing.T) {
	original := RatingUpdated{
		UserID:     "user-123",
		GameType:   "backgammon",
		RatingType: "br-rapid",
		OldRating:  1500,
		NewRating:  1512,
		GameID:     "g-987654",
		UpdatedAt:  time.Date(2026, 4, 13, 11, 0, 0, 0, time.UTC),
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded RatingUpdated
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("roundtrip mismatch:\noriginal: %+v\ndecoded:  %+v", original, decoded)
	}
}

func TestRatingUpdatedValidate(t *testing.T) {
	valid := RatingUpdated{
		UserID: "u", GameType: "g", RatingType: "r",
		OldRating: 1500, NewRating: 1510, GameID: "g1",
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid rating returned error: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*RatingUpdated)
	}{
		{"empty user_id", func(m *RatingUpdated) { m.UserID = "" }},
		{"empty game_type", func(m *RatingUpdated) { m.GameType = "" }},
		{"empty rating_type", func(m *RatingUpdated) { m.RatingType = "" }},
		{"empty game_id", func(m *RatingUpdated) { m.GameID = "" }},
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
