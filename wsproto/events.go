// Package wsproto defines the WebSocket wire vocabulary shared across
// vog-game, vog-lobby, and vog-web.
//
// All identifiers follow <domain>.<action> lowercase dot-separated form.
// Consumers MUST import these constants instead of inlining string
// literals. See vog-arch/docs/ws-protocol.md for the full mapping
// including frame types and payload field names.
package wsproto

// Event codes — values of the `e` field inside r.event payloads.
const (
	// EvtWhoJoin ("w.join") — a user appeared in the room's presence set.
	EvtWhoJoin = "w.join"
	// EvtWhoLeave ("w.leave") — a user disappeared from the room.
	EvtWhoLeave = "w.leave"
	// EvtWhoSnap ("w.snap") — full presence snapshot (on r.enter reply).
	EvtWhoSnap = "w.snap"

	// EvtPlayerEnter ("p.enter") — a player joined an active game session.
	EvtPlayerEnter = "p.enter"
	// EvtPlayerLeave ("p.leave") — a player left an active game session.
	EvtPlayerLeave = "p.leave"

	// EvtGameFreeze ("g.freeze") — game paused (disconnect/timeout/left_room).
	EvtGameFreeze = "g.freeze"
	// EvtGameResume ("g.resume") — game resumed.
	EvtGameResume = "g.resume"

	// EvtRoomState ("r.state") — full room state snapshot (players/seeks/active_games).
	EvtRoomState = "r.state"
)

// Scope values — values of the `scope` field on r.enter frames.
const (
	// ScopePresence ("p") — receive w.* only, no game events.
	ScopePresence = "p"
	// ScopeFull ("f") — receive every room event.
	ScopeFull = "f"
)
