package vogcluster_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsd "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	vogcluster "vogclub.com/vog-cluster"
)

func startEmbeddedNATS(t *testing.T) *server.Server {
	t.Helper()
	opts := &server.Options{
		Host:       "127.0.0.1",
		Port:       -1, // pick a free port
		NoLog:      true,
		NoSigs:     true,
		MaxPayload: 1024 * 1024,
	}
	s := natsd.RunServer(opts)
	if !s.ReadyForConnections(2 * time.Second) {
		t.Fatal("embedded NATS server did not become ready")
	}
	t.Cleanup(func() { s.Shutdown() })
	return s
}

func TestRoundtripGameRegister(t *testing.T) {
	srv := startEmbeddedNATS(t)
	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer nc.Close()

	received := make(chan vogcluster.GameRegister, 1)
	sub, err := nc.Subscribe(vogcluster.SubjectClusterGameRegister, func(m *nats.Msg) {
		var msg vogcluster.GameRegister
		if err := vogcluster.Decode(m.Data, &msg); err != nil {
			t.Errorf("decode: %v", err)
			return
		}
		received <- msg
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	original := vogcluster.GameRegister{
		InstanceID: "game-7",
		Capacity:   5000,
		Version:    "1.2.3",
		Address:    "10.0.0.42:9001",
	}
	data, err := vogcluster.Encode(original)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err := nc.Publish(vogcluster.SubjectClusterGameRegister, data); err != nil {
		t.Fatalf("publish: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	select {
	case got := <-received:
		if got.InstanceID != original.InstanceID {
			t.Errorf("instance_id mismatch: got %q want %q", got.InstanceID, original.InstanceID)
		}
		if got.Capacity != original.Capacity {
			t.Errorf("capacity mismatch: got %d want %d", got.Capacity, original.Capacity)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for message")
	}
}

func TestRoundtripRoomBroadcastWithDynamicSubject(t *testing.T) {
	srv := startEmbeddedNATS(t)
	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer nc.Close()

	roomID := "br-rapid-1"
	subject := vogcluster.SubjectGameRoomBroadcast(roomID)

	received := make(chan vogcluster.RoomBroadcast, 1)
	sub, err := nc.Subscribe(subject, func(m *nats.Msg) {
		var msg vogcluster.RoomBroadcast
		if err := vogcluster.Decode(m.Data, &msg); err != nil {
			t.Errorf("decode: %v", err)
			return
		}
		received <- msg
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	original := vogcluster.RoomBroadcast{
		Event:    "table.move",
		RoomID:   roomID,
		Sequence: 1,
		Payload:  json.RawMessage(`{"table_id":7}`),
	}
	data, err := vogcluster.Encode(original)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err := nc.Publish(subject, data); err != nil {
		t.Fatalf("publish: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	select {
	case got := <-received:
		if got.RoomID != roomID {
			t.Errorf("room_id mismatch: got %q want %q", got.RoomID, roomID)
		}
		if got.Sequence != original.Sequence {
			t.Errorf("sequence mismatch: got %d want %d", got.Sequence, original.Sequence)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for room broadcast")
	}
}

func TestSubscribeWithWildcard(t *testing.T) {
	srv := startEmbeddedNATS(t)
	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer nc.Close()

	// Coordinator-side subscribe to all heartbeats from any instance.
	received := make(chan vogcluster.GameHeartbeat, 4)
	sub, err := nc.Subscribe("vog.cluster.game.heartbeat.*", func(m *nats.Msg) {
		var msg vogcluster.GameHeartbeat
		if err := vogcluster.Decode(m.Data, &msg); err != nil {
			t.Errorf("decode: %v", err)
			return
		}
		received <- msg
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	for _, id := range []string{"game-1", "game-2", "game-3"} {
		hb := vogcluster.GameHeartbeat{
			InstanceID: id,
			Status:     vogcluster.InstanceStatusActive,
			SlotsTotal: 1000,
			SlotsUsed:  100,
		}
		data, err := vogcluster.Encode(hb)
		if err != nil {
			t.Fatalf("encode: %v", err)
		}
		if err := nc.Publish(vogcluster.SubjectClusterGameHeartbeat(id), data); err != nil {
			t.Fatalf("publish: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	seen := map[string]bool{}
	for len(seen) < 3 {
		select {
		case got := <-received:
			seen[got.InstanceID] = true
		case <-ctx.Done():
			t.Fatalf("timeout, only saw %v", seen)
		}
	}
}
