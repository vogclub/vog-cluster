package vogcluster

import (
	"fmt"
	"strings"
)

// Static subject names. These have no dynamic components and can be
// used directly in nats.Conn.Publish / Subscribe.
const (
	// SubjectClusterGameRegister: game instance -> coordinator,
	// registers a new instance with capacity and metadata.
	SubjectClusterGameRegister = "vog.cluster.game.register"

	// SubjectClusterRoutingUpdate: coordinator -> all Lobby instances,
	// publishes updates to the room->instance routing table.
	SubjectClusterRoutingUpdate = "vog.cluster.routing.update"

	// SubjectStatsOnline: coordinator -> Lobby, periodic aggregate
	// online counters for guests and authorized clients.
	SubjectStatsOnline = "vog.stats.online"
)

// SubjectClusterGameHeartbeat returns the subject a specific game
// instance publishes its periodic heartbeat to.
func SubjectClusterGameHeartbeat(instanceID string) string {
	return buildSubject("vog.cluster.game.heartbeat", instanceID)
}

// SubjectClusterGameDrain returns the subject a game instance uses
// to announce it is draining.
func SubjectClusterGameDrain(instanceID string) string {
	return buildSubject("vog.cluster.game.drain", instanceID)
}

// SubjectClusterGameDeregister returns the subject a game instance
// uses to announce it is shutting down.
func SubjectClusterGameDeregister(instanceID string) string {
	return buildSubject("vog.cluster.game.deregister", instanceID)
}

// SubjectClusterRoomsAssign returns the subject the coordinator uses
// to send room assignments to a specific game instance.
func SubjectClusterRoomsAssign(instanceID string) string {
	return buildSubject("vog.cluster.rooms.assign", instanceID)
}

// SubjectClusterRoomsRelease returns the subject the coordinator uses
// to instruct a game instance to release rooms.
func SubjectClusterRoomsRelease(instanceID string) string {
	return buildSubject("vog.cluster.rooms.release", instanceID)
}

// SubjectClusterRoomsPrepare returns the subject the coordinator uses
// to ask a game instance to prepare rooms (during migration).
func SubjectClusterRoomsPrepare(instanceID string) string {
	return buildSubject("vog.cluster.rooms.prepare", instanceID)
}

// SubjectClusterCommand returns the subject the coordinator uses to
// send admin-initiated commands to a specific game instance.
func SubjectClusterCommand(instanceID string) string {
	return buildSubject("vog.cluster.command", instanceID)
}

// SubjectGameRoomInput returns the subject Lobby uses to deliver
// player actions to the game instance hosting the given room.
func SubjectGameRoomInput(roomID string) string {
	return buildTripleSubject("vog.game.room", roomID, "input")
}

// SubjectGameRoomBroadcast returns the subject a game instance
// publishes room events to. All Lobby instances with clients in
// the room subscribe to it.
func SubjectGameRoomBroadcast(roomID string) string {
	return buildTripleSubject("vog.game.room", roomID, "broadcast")
}

// SubjectLobbyOutput returns the subject a game instance uses to
// deliver targeted messages to a specific Lobby instance.
func SubjectLobbyOutput(lobbyInstanceID string) string {
	return buildTripleSubject("vog.lobby", lobbyInstanceID, "output")
}

// SubjectRatingUpdated returns the subject vog-games publishes
// rating updates to, partitioned by game type.
func SubjectRatingUpdated(gameType string) string {
	return buildSubject("vog.rating.updated", gameType)
}

// RequiresJetStream reports whether the given subject must be backed
// by a JetStream stream rather than plain Core NATS. The classification
// matches the spec's transport table.
//
// Only subjects defined in this package are classified. Any subject
// outside the known prefixes returns false (Core NATS) — this is the
// safe default for unknown traffic, but means that any new subject
// added to this package MUST also be added to the switch below.
// Forgetting to do so causes silent loss for messages that should have
// been durable. Tests in subjects_test.go enforce that every defined
// subject is covered.
func RequiresJetStream(subject string) bool {
	switch {
	case subject == SubjectClusterRoutingUpdate:
		return true
	case subject == SubjectStatsOnline:
		return false
	case strings.HasPrefix(subject, "vog.cluster."):
		return true
	case strings.HasPrefix(subject, "vog.rating.updated."):
		return true
	case strings.HasPrefix(subject, "vog.game.room."):
		return false
	case strings.HasPrefix(subject, "vog.lobby."):
		return false
	}
	return false
}

// buildSubject joins a prefix with a single token, validating the token.
func buildSubject(prefix, token string) string {
	mustValidToken(token)
	return prefix + "." + token
}

// buildTripleSubject builds a "{prefix}.{token}.{suffix}" subject,
// validating the token.
func buildTripleSubject(prefix, token, suffix string) string {
	mustValidToken(token)
	return prefix + "." + token + "." + suffix
}

// mustValidToken panics if token is empty or contains characters that
// are illegal in a NATS subject token (dot, space, wildcards).
func mustValidToken(token string) {
	if token == "" {
		panic("vogcluster: empty subject token")
	}
	if strings.ContainsAny(token, ". \t\n\r\x00*>") {
		panic(fmt.Sprintf("vogcluster: invalid subject token %q", token))
	}
}
