package vogcluster

import (
	"fmt"
	"time"
)

// StatsServerCount is the per-game-type online count.
type StatsServerCount struct {
	GameType string `json:"game_type"`
	Count    int    `json:"count"`
}

// StatsRoomCount is the per-room online count.
type StatsRoomCount struct {
	RoomID       string `json:"room_id"`
	Players      int    `json:"players"`
	TablesActive int    `json:"tables_active"`
}

// StatsOnline is the periodic aggregate-online snapshot the coordinator
// publishes on SubjectStatsOnline. Lobby relays it to SSE streams
// (guests) and to authorized WebSocket clients.
type StatsOnline struct {
	// Total is the sum of all online connections across the cluster.
	Total int `json:"total"`

	// Servers reports per-game-type online counts.
	Servers []StatsServerCount `json:"servers,omitempty"`

	// Rooms reports per-room counts. May be a partial slice
	// (e.g. only rooms with non-zero traffic) to keep payload small.
	Rooms []StatsRoomCount `json:"rooms,omitempty"`

	// ObservedAt is the wall-clock time the snapshot was assembled.
	ObservedAt time.Time `json:"observed_at,omitempty"`
}

// Validate reports whether the snapshot is well-formed.
func (m StatsOnline) Validate() error {
	if m.Total < 0 {
		return fmt.Errorf("vogcluster: StatsOnline.total must be non-negative, got %d", m.Total)
	}
	for i, s := range m.Servers {
		if s.GameType == "" {
			return fmt.Errorf("vogcluster: StatsOnline.servers[%d].game_type is required", i)
		}
		if s.Count < 0 {
			return fmt.Errorf("vogcluster: StatsOnline.servers[%d].count must be non-negative, got %d", i, s.Count)
		}
	}
	for i, r := range m.Rooms {
		if r.RoomID == "" {
			return fmt.Errorf("vogcluster: StatsOnline.rooms[%d].room_id is required", i)
		}
		if r.Players < 0 {
			return fmt.Errorf("vogcluster: StatsOnline.rooms[%d].players must be non-negative", i)
		}
		if r.TablesActive < 0 {
			return fmt.Errorf("vogcluster: StatsOnline.rooms[%d].tables_active must be non-negative", i)
		}
	}
	return nil
}
