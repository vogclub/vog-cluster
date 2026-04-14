package vogcluster_test

import (
	"testing"

	vogcluster "vogclub.com/vog-cluster"
)

func TestHeaderClaimerHashConstant(t *testing.T) {
	if got := vogcluster.HeaderClaimerHash; got != "X-Claimer-Hash" {
		t.Fatalf("HeaderClaimerHash = %q, want %q", got, "X-Claimer-Hash")
	}
}

func TestRoomPrepareResponseValidate(t *testing.T) {
	cases := []struct {
		name    string
		msg     vogcluster.RoomPrepareResponse
		wantErr bool
	}{
		{"ok", vogcluster.RoomPrepareResponse{MigrationID: "m-1", Accepted: true}, false},
		{"reject with reason", vogcluster.RoomPrepareResponse{MigrationID: "m-1", Accepted: false, Reason: "out_of_capacity"}, false},
		{"missing migration id", vogcluster.RoomPrepareResponse{Accepted: true}, true},
		{"reject without reason", vogcluster.RoomPrepareResponse{MigrationID: "m-1", Accepted: false}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() err = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestRoomReadyValidate(t *testing.T) {
	if err := (vogcluster.RoomReady{MigrationID: "m-1", InstanceID: "game-1", RoomID: "1:101"}).Validate(); err != nil {
		t.Fatalf("valid RoomReady: %v", err)
	}
	if err := (vogcluster.RoomReady{}).Validate(); err == nil {
		t.Fatalf("empty RoomReady should fail")
	}
}

func TestInstanceRegisterReplyValidate(t *testing.T) {
	cases := []struct {
		name    string
		msg     vogcluster.InstanceRegisterReply
		wantErr bool
	}{
		{"accept", vogcluster.InstanceRegisterReply{Accepted: true, InstanceID: "game-1"}, false},
		{"reject with reason", vogcluster.InstanceRegisterReply{Accepted: false, InstanceID: "game-1", Reason: "claimed_by_another"}, false},
		{"missing id", vogcluster.InstanceRegisterReply{Accepted: true}, true},
		{"reject without reason", vogcluster.InstanceRegisterReply{Accepted: false, InstanceID: "game-1"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() err = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestAdminMigrateRequestValidate(t *testing.T) {
	good := vogcluster.AdminMigrateRequest{
		ServerID:   1,
		RoomID:     101,
		ToInstance: "game-2",
		IssuedBy:   "admin@example.com",
	}
	if err := good.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	bad := vogcluster.AdminMigrateRequest{ServerID: 0, RoomID: 101, ToInstance: "game-2"}
	if err := bad.Validate(); err == nil {
		t.Fatalf("server_id=0 should fail")
	}
}

func TestAdminDrainRequestValidate(t *testing.T) {
	if err := (vogcluster.AdminDrainRequest{InstanceID: "game-1", IssuedBy: "admin"}).Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if err := (vogcluster.AdminDrainRequest{}).Validate(); err == nil {
		t.Fatalf("empty should fail")
	}
}

func TestAdminRebalanceRequestValidate(t *testing.T) {
	if err := (vogcluster.AdminRebalanceRequest{IssuedBy: "admin"}).Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if err := (vogcluster.AdminRebalanceRequest{}).Validate(); err == nil {
		t.Fatalf("missing issued_by should fail")
	}
}

func TestInstanceStatusEventValidate(t *testing.T) {
	if err := (vogcluster.InstanceStatusEvent{InstanceID: "game-1", Status: vogcluster.InstanceStatusActive}).Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if err := (vogcluster.InstanceStatusEvent{InstanceID: "", Status: vogcluster.InstanceStatusActive}).Validate(); err == nil {
		t.Fatalf("missing id should fail")
	}
	if err := (vogcluster.InstanceStatusEvent{InstanceID: "game-1", Status: "bogus"}).Validate(); err == nil {
		t.Fatalf("bad status should fail")
	}
}
