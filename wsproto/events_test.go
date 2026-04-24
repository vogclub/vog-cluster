package wsproto

import "testing"

func TestEventCodes(t *testing.T) {
	cases := map[string]string{
		"EvtWhoJoin":     EvtWhoJoin,
		"EvtWhoLeave":    EvtWhoLeave,
		"EvtWhoSnap":     EvtWhoSnap,
		"EvtPlayerEnter": EvtPlayerEnter,
		"EvtPlayerLeave": EvtPlayerLeave,
		"EvtGameFreeze":  EvtGameFreeze,
		"EvtGameResume":  EvtGameResume,
		"EvtRoomState":   EvtRoomState,
	}
	want := map[string]string{
		"EvtWhoJoin":     "w.join",
		"EvtWhoLeave":    "w.leave",
		"EvtWhoSnap":     "w.snap",
		"EvtPlayerEnter": "p.enter",
		"EvtPlayerLeave": "p.leave",
		"EvtGameFreeze":  "g.freeze",
		"EvtGameResume":  "g.resume",
		"EvtRoomState":   "r.state",
	}
	for k, got := range cases {
		if got != want[k] {
			t.Errorf("%s = %q, want %q", k, got, want[k])
		}
	}
}

func TestScopeValues(t *testing.T) {
	if ScopePresence != "p" {
		t.Errorf("ScopePresence = %q, want %q", ScopePresence, "p")
	}
	if ScopeFull != "f" {
		t.Errorf("ScopeFull = %q, want %q", ScopeFull, "f")
	}
}
