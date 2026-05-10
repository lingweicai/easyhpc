// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lingweicai/easyhpc/bridge/internal/slurm"
)

const sampleSlurmReservationsJSON = `{
  "reservations": [
    {
      "name": "workshop",
      "start_time": {"set": true, "infinite": false, "number": 1778000000},
      "end_time": {"set": true, "infinite": false, "number": 1778003600},
      "duration": "01:00:00",
      "users": "alice,bob",
      "accounts": "research",
      "groups": "hpc",
      "node_list": "c[31-32]",
      "node_cnt": {"set": true, "infinite": false, "number": 2},
      "core_cnt": {"set": true, "infinite": false, "number": 64},
      "partition_name": "normal",
      "features": "gpu",
      "licenses": "ansys:2",
      "tres": "cpu=64,mem=128000M,node=2",
      "burst_buffer": "bb1",
      "max_start_delay": {"set": true, "infinite": false, "number": 120},
      "flags": ["MAINT", "IGNORE_JOBS"],
      "state": ["ACTIVE"]
    }
  ]
}`

func TestMapSlurmReservationRaw_KeyFields(t *testing.T) {
	var resp slurm.SlurmReservationsResponse
	if err := json.Unmarshal([]byte(sampleSlurmReservationsJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	raws := resp.Items()
	if len(raws) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(raws))
	}

	r := slurm.MapSlurmReservationRaw(raws[0])

	if r.ReservationID != "workshop" {
		t.Errorf("ReservationID: got %q, want workshop", r.ReservationID)
	}
	if r.ReservationName != "workshop" {
		t.Errorf("ReservationName: got %q, want workshop", r.ReservationName)
	}
	if r.StartTime == nil {
		t.Fatal("StartTime: got nil, want value")
	}
	if r.EndTime == nil {
		t.Fatal("EndTime: got nil, want value")
	}
	if !r.StartTime.Equal(time.Unix(1778000000, 0).UTC()) {
		t.Errorf("StartTime: got %v, want %v", *r.StartTime, time.Unix(1778000000, 0).UTC())
	}
	if !r.EndTime.Equal(time.Unix(1778003600, 0).UTC()) {
		t.Errorf("EndTime: got %v, want %v", *r.EndTime, time.Unix(1778003600, 0).UTC())
	}
	if r.NodeList != "c[31-32]" {
		t.Errorf("NodeList: got %q, want c[31-32]", r.NodeList)
	}
	if r.PartitionName != "normal" {
		t.Errorf("PartitionName: got %q, want normal", r.PartitionName)
	}
	if r.Users != "alice,bob" {
		t.Errorf("Users: got %q, want alice,bob", r.Users)
	}
	if r.Accounts != "research" {
		t.Errorf("Accounts: got %q, want research", r.Accounts)
	}
	if r.State != "ACTIVE" {
		t.Errorf("State: got %q, want ACTIVE", r.State)
	}
	if len(r.Flags) != 2 || r.Flags[0] != "MAINT" || r.Flags[1] != "IGNORE_JOBS" {
		t.Errorf("Flags: got %v, want [MAINT IGNORE_JOBS]", r.Flags)
	}
	if r.NodeCount != 2 {
		t.Errorf("NodeCount: got %d, want 2", r.NodeCount)
	}
	if r.CoreCount != 64 {
		t.Errorf("CoreCount: got %d, want 64", r.CoreCount)
	}
	if r.CPUs != 64 {
		t.Errorf("CPUs: got %d, want 64", r.CPUs)
	}
	if r.MaxStartDelay != "120" {
		t.Errorf("MaxStartDelay: got %q, want 120", r.MaxStartDelay)
	}
	if r.Nodes == nil || len(r.Nodes) != 0 {
		t.Errorf("Nodes: got %#v, want empty non-nil slice", r.Nodes)
	}
}
