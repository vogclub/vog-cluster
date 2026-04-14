package vogcluster

import (
	"strings"
	"testing"
)

func TestStaticSubjects(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		want    string
	}{
		{"ClusterGameRegister", SubjectClusterGameRegister, "vog.cluster.game.register"},
		{"ClusterRoutingUpdate", SubjectClusterRoutingUpdate, "vog.cluster.routing.update"},
		{"StatsOnline", SubjectStatsOnline, "vog.stats.online"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.subject != tc.want {
				t.Errorf("%s = %q, want %q", tc.name, tc.subject, tc.want)
			}
		})
	}
}

func TestSubjectBuilders(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"ClusterGameHeartbeat", SubjectClusterGameHeartbeat("game-7"), "vog.cluster.game.heartbeat.game-7"},
		{"ClusterGameDrain", SubjectClusterGameDrain("game-7"), "vog.cluster.game.drain.game-7"},
		{"ClusterGameDeregister", SubjectClusterGameDeregister("game-7"), "vog.cluster.game.deregister.game-7"},
		{"ClusterRoomsAssign", SubjectClusterRoomsAssign("game-7"), "vog.cluster.rooms.assign.game-7"},
		{"ClusterRoomsRelease", SubjectClusterRoomsRelease("game-7"), "vog.cluster.rooms.release.game-7"},
		{"ClusterRoomsPrepare", SubjectClusterRoomsPrepare("game-7"), "vog.cluster.rooms.prepare.game-7"},
		{"ClusterCommand", SubjectClusterCommand("game-7"), "vog.cluster.command.game-7"},
		{"GameRoomInput", SubjectGameRoomInput("42"), "vog.game.room.42.input"},
		{"GameRoomBroadcast", SubjectGameRoomBroadcast("42"), "vog.game.room.42.broadcast"},
		{"LobbyOutput", SubjectLobbyOutput("lobby-2"), "vog.lobby.lobby-2.output"},
		{"RatingUpdated", SubjectRatingUpdated("backgammon"), "vog.rating.updated.backgammon"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
			}
		})
	}
}

func TestSubjectBuildersRejectEmpty(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("expected panic on empty instance ID")
		}
	}()
	_ = SubjectClusterGameHeartbeat("")
}

func TestSubjectBuildersRejectInvalid(t *testing.T) {
	// NATS subject tokens must not contain dots, whitespace, control
	// characters, or wildcard characters.
	cases := []struct {
		name  string
		token string
	}{
		{"dot", "foo.bar"},
		{"space", "foo bar"},
		{"tab", "foo\tbar"},
		{"newline", "foo\nbar"},
		{"carriage_return", "foo\rbar"},
		{"null", "foo\x00bar"},
		{"asterisk", "foo*"},
		{"greater_than", "foo>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("expected panic on token %q", tc.token)
				}
			}()
			_ = SubjectGameRoomInput(tc.token)
		})
	}
}

func TestRequiresJetStream(t *testing.T) {
	jsSubjects := []string{
		SubjectClusterGameRegister,
		SubjectClusterGameHeartbeat("game-7"),
		SubjectClusterGameDrain("game-7"),
		SubjectClusterGameDeregister("game-7"),
		SubjectClusterRoomsAssign("game-7"),
		SubjectClusterRoomsRelease("game-7"),
		SubjectClusterRoomsPrepare("game-7"),
		SubjectClusterCommand("game-7"),
		SubjectClusterRoutingUpdate,
		SubjectRatingUpdated("chess"),
	}
	for _, s := range jsSubjects {
		if RequiresJetStream(s) != true {
			t.Errorf("RequiresJetStream(%q) = false, want true", s)
		}
	}

	coreSubjects := []string{
		SubjectGameRoomInput("42"),
		SubjectGameRoomBroadcast("42"),
		SubjectLobbyOutput("lobby-2"),
		SubjectStatsOnline,
	}
	for _, s := range coreSubjects {
		if RequiresJetStream(s) != false {
			t.Errorf("RequiresJetStream(%q) = true, want false", s)
		}
	}
}

func TestSubjectsAreUniquePrefixes(t *testing.T) {
	// All vog-cluster subjects should start with "vog."
	subjects := []string{
		SubjectClusterGameRegister,
		SubjectClusterGameHeartbeat("x"),
		SubjectGameRoomInput("x"),
		SubjectLobbyOutput("x"),
		SubjectStatsOnline,
		SubjectRatingUpdated("x"),
	}
	for _, s := range subjects {
		if !strings.HasPrefix(s, "vog.") {
			t.Errorf("subject %q does not start with vog.", s)
		}
	}
}

func TestRequiresJetStreamCoversAllKnownSubjects(t *testing.T) {
	// Every subject this package can produce must be explicitly
	// classified by RequiresJetStream. If you add a new subject and
	// this test fails, add the corresponding case to the switch.
	known := []string{
		SubjectClusterGameRegister,
		SubjectClusterRoutingUpdate,
		SubjectStatsOnline,
		SubjectClusterGameHeartbeat("x"),
		SubjectClusterGameDrain("x"),
		SubjectClusterGameDeregister("x"),
		SubjectClusterRoomsAssign("x"),
		SubjectClusterRoomsRelease("x"),
		SubjectClusterRoomsPrepare("x"),
		SubjectClusterCommand("x"),
		SubjectGameRoomInput("x"),
		SubjectGameRoomBroadcast("x"),
		SubjectLobbyOutput("x"),
		SubjectRatingUpdated("x"),
	}
	// Build the set of expected JetStream subjects from the spec.
	jetstream := map[string]bool{
		SubjectClusterGameRegister:        true,
		SubjectClusterRoutingUpdate:       true,
		SubjectClusterGameHeartbeat("x"):  true,
		SubjectClusterGameDrain("x"):      true,
		SubjectClusterGameDeregister("x"): true,
		SubjectClusterRoomsAssign("x"):    true,
		SubjectClusterRoomsRelease("x"):   true,
		SubjectClusterRoomsPrepare("x"):   true,
		SubjectClusterCommand("x"):        true,
		SubjectRatingUpdated("x"):         true,
	}
	for _, subj := range known {
		got := RequiresJetStream(subj)
		want := jetstream[subj]
		if got != want {
			t.Errorf("RequiresJetStream(%q) = %v, want %v", subj, got, want)
		}
	}
}

func TestRequiresJetStreamUnknownReturnsFalse(t *testing.T) {
	if RequiresJetStream("vog.future.thing") != false {
		t.Errorf("unknown subject should default to Core NATS (false)")
	}
	if RequiresJetStream("not.a.vog.subject") != false {
		t.Errorf("non-vog subject should return false")
	}
}
