# vog-cluster

NATS protocol contracts for the vog horizontally-scaled game cluster.

This package provides typed message structs, NATS subject constants and
builders, transport classification (Core NATS vs JetStream), and a small
codec helper that validates messages on encode and decode. It is the
shared vocabulary for vog-spaces (coordinator), vog-game (game
instances), vog-lobby (gateway), vog-games (rating publisher), and
vog-tourneys (tournament events).

See the [horizontal scaling spec](https://github.com/khorost/vog-arch/blob/main/docs/superpowers/specs/2026-04-13-horizontal-scaling-design.md)
in vog-arch for the full architecture.

## Install

```bash
go get vogclub.com/vog-cluster
```

## Subjects

| Subject | Direction | Transport | Message Type |
|---------|-----------|-----------|--------------|
| `vog.cluster.game.register` | game -> coordinator | JetStream | `GameRegister` |
| `vog.cluster.game.heartbeat.{instance_id}` | game -> coordinator | JetStream | `GameHeartbeat` |
| `vog.cluster.game.drain.{instance_id}` | game -> coordinator | JetStream | `GameDrain` |
| `vog.cluster.game.deregister.{instance_id}` | game -> coordinator | JetStream | `GameDeregister` |
| `vog.cluster.rooms.assign.{instance_id}` | coordinator -> game | JetStream | `RoomAssign` |
| `vog.cluster.rooms.release.{instance_id}` | coordinator -> game | JetStream | `RoomRelease` |
| `vog.cluster.rooms.prepare.{instance_id}` | coordinator -> game | JetStream | `RoomPrepare` |
| `vog.cluster.command.{instance_id}` | coordinator -> game | JetStream | `ClusterCommand` |
| `vog.cluster.routing.update` | coordinator -> all lobby | JetStream | `RoutingUpdate` |
| `vog.game.room.{room_id}.input` | lobby -> game | Core NATS | `RoomInput` |
| `vog.game.room.{room_id}.broadcast` | game -> lobby (fan-out) | Core NATS | `RoomBroadcast` |
| `vog.lobby.{lobby_instance_id}.output` | game -> specific lobby | Core NATS | `LobbyOutput` |
| `vog.stats.online` | coordinator -> lobby | Core NATS | `StatsOnline` |
| `vog.rating.updated.{game_type}` | vog-games -> game | JetStream | `RatingUpdated` |

Subjects with placeholders (`{...}`) are built with helper functions.
Use `RequiresJetStream(subject)` to classify any subject at runtime.

## Usage

### Publishing

```go
import (
    "github.com/nats-io/nats.go"
    vogcluster "vogclub.com/vog-cluster"
)

nc, _ := nats.Connect(nats.DefaultURL)

msg := vogcluster.GameRegister{
    InstanceID: "game-7",
    Capacity:   5000,
    Version:    "1.2.3",
    Address:    "10.0.0.42:9001",
}

data, err := vogcluster.Encode(msg)
if err != nil {
    // Validation failed before anything went over the wire.
    return err
}
nc.Publish(vogcluster.SubjectClusterGameRegister, data)
```

### Subscribing

```go
nc.Subscribe(vogcluster.SubjectClusterGameRegister, func(m *nats.Msg) {
    var msg vogcluster.GameRegister
    if err := vogcluster.Decode(m.Data, &msg); err != nil {
        // Malformed or invalid message.
        log.Printf("bad register: %v", err)
        return
    }
    // msg is guaranteed to have all required fields populated.
    handleRegister(msg)
})
```

### Wildcards

To receive heartbeats from any instance, subscribe with a NATS wildcard:

```go
nc.Subscribe("vog.cluster.game.heartbeat.*", func(m *nats.Msg) {
    var hb vogcluster.GameHeartbeat
    if err := vogcluster.Decode(m.Data, &hb); err != nil {
        return
    }
    // hb.InstanceID identifies the source.
})
```

## Validation

Every message type implements `Validate() error`. `Encode` calls it
before marshalling and `Decode` calls it after unmarshalling, so
malformed payloads never reach handler code. Call `Validate()` directly
when building messages programmatically and you want an early failure.

## Transport classification

JetStream-backed subjects are listed in the table above. Use
`RequiresJetStream(subject string) bool` to decide which publish path
to use at runtime, or hardcode the choice per call site.

JetStream is required for messages that must not be lost: cluster
coordination commands and rating updates. Core NATS is used for
high-volume real-time game traffic where loss is acceptable because
the next event will overwrite state.

## Testing

```bash
go test ./...
```

The integration tests start an embedded NATS server (no external
dependencies) and exercise full publish/subscribe roundtrips for
representative message types.
</content>
</invoke>