package vogcluster

import (
	"testing"
)

func TestHeaderClaimerHashConstant(t *testing.T) {
	if got := HeaderClaimerHash; got != "X-Claimer-Hash" {
		t.Fatalf("HeaderClaimerHash = %q, want %q", got, "X-Claimer-Hash")
	}
}

func TestRoomPrepareResponseValidate(t *testing.T) {
	cases := []struct {
		name    string
		msg     RoomPrepareResponse
		wantErr bool
	}{
		{"ok", RoomPrepareResponse{MigrationID: "m-1", Accepted: true}, false},
		{"reject with reason", RoomPrepareResponse{MigrationID: "m-1", Accepted: false, Reason: "out_of_capacity"}, false},
		{"missing migration id", RoomPrepareResponse{Accepted: true}, true},
		{"reject without reason", RoomPrepareResponse{MigrationID: "m-1", Accepted: false}, true},
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
	cases := []struct {
		name    string
		msg     RoomReady
		wantErr bool
	}{
		{"ok", RoomReady{MigrationID: "m-1", InstanceID: "game-1", RoomID: "1:101"}, false},
		{"missing migration id", RoomReady{InstanceID: "game-1", RoomID: "1:101"}, true},
		{"missing instance id", RoomReady{MigrationID: "m-1", RoomID: "1:101"}, true},
		{"missing room id", RoomReady{MigrationID: "m-1", InstanceID: "game-1"}, true},
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

func TestInstanceRegisterReplyValidate(t *testing.T) {
	cases := []struct {
		name    string
		msg     InstanceRegisterReply
		wantErr bool
	}{
		{"accept", InstanceRegisterReply{Accepted: true, InstanceID: "game-1"}, false},
		{"reject with reason", InstanceRegisterReply{Accepted: false, InstanceID: "game-1", Reason: "claimed_by_another"}, false},
		{"missing id", InstanceRegisterReply{Accepted: true}, true},
		{"reject without reason", InstanceRegisterReply{Accepted: false, InstanceID: "game-1"}, true},
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
	good := AdminMigrateRequest{
		ServerID:   1,
		RoomID:     101,
		ToInstance: "game-2",
		IssuedBy:   "admin@example.com",
	}
	if err := good.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	bad := AdminMigrateRequest{ServerID: 0, RoomID: 101, ToInstance: "game-2"}
	if err := bad.Validate(); err == nil {
		t.Fatalf("server_id=0 should fail")
	}
}

func TestAdminDrainRequestValidate(t *testing.T) {
	if err := (AdminDrainRequest{InstanceID: "game-1", IssuedBy: "admin"}).Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if err := (AdminDrainRequest{}).Validate(); err == nil {
		t.Fatalf("empty should fail")
	}
}

func TestAdminRebalanceRequestValidate(t *testing.T) {
	if err := (AdminRebalanceRequest{IssuedBy: "admin"}).Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if err := (AdminRebalanceRequest{}).Validate(); err == nil {
		t.Fatalf("missing issued_by should fail")
	}
}

func TestInstanceStatusEventValidate(t *testing.T) {
	if err := (InstanceStatusEvent{InstanceID: "game-1", Status: InstanceStatusActive}).Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if err := (InstanceStatusEvent{InstanceID: "", Status: InstanceStatusActive}).Validate(); err == nil {
		t.Fatalf("missing id should fail")
	}
	if err := (InstanceStatusEvent{InstanceID: "game-1", Status: "bogus"}).Validate(); err == nil {
		t.Fatalf("bad status should fail")
	}
}
