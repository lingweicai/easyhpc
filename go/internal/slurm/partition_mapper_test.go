// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm_test

import (
	"encoding/json"
	"testing"

	"github.com/lingweicai/easyhpc/bridge/internal/slurm"
)

// sampleSlurmPartitionsJSON is an excerpt of the JSON produced by
// "scontrol show partitions --json" for two partitions:
//   - "normal": default, INFINITE max time, no limits, NO oversubscribe
//   - "fat":    non-default, 1-day max time, restricted accounts/QOS,
//     FORCE:4 oversubscribe, specific CPU limits
const sampleSlurmPartitionsJSON = `{
  "partitions": [
    {
      "name": "normal",
      "alternate": "",
      "node_sets": "",
      "nodes": {
        "allowed_allocation": "ALL",
        "configured": "c[31-32]",
        "total": 2
      },
      "accounts": { "allowed": "ALL", "deny": "" },
      "groups":   { "allowed": "ALL" },
      "qos":      { "allowed": "ALL", "deny": "", "assigned": "" },
      "defaults": {
        "time": {"set": false, "infinite": false, "number": 0},
        "job":  ""
      },
      "maximums": {
        "nodes":          {"set": false, "infinite": true,  "number": 0},
        "time":           {"set": false, "infinite": true,  "number": 0},
        "cpus_per_node":  {"set": false, "infinite": false, "number": 0},
        "cpus_per_socket":{"set": false, "infinite": false, "number": 0},
        "over_time_limit":{"set": false, "infinite": false, "number": 0},
        "oversubscribe":  {"jobs": 1, "flags": []}
      },
      "minimums": { "nodes": 1 },
      "priority": { "job_factor": 1, "tier": 1 },
      "cpus":      { "total": 4 },
      "partition": { "state": ["UP"] },
      "tres": {
        "billing_weights": "",
        "configured": "cpu=4,mem=15872M,node=2,billing=4"
      },
      "grace_time": 0
    },
    {
      "name": "fat",
      "alternate": "normal",
      "node_sets": "",
      "nodes": {
        "allowed_allocation": "sms",
        "configured": "c[41-42]",
        "total": 2
      },
      "accounts": { "allowed": "research", "deny": "marketing" },
      "groups":   { "allowed": "ALL" },
      "qos":      { "allowed": "high", "deny": "low", "assigned": "normal" },
      "defaults": {
        "time": {"set": true, "infinite": false, "number": 60},
        "job":  ""
      },
      "maximums": {
        "nodes":          {"set": true,  "infinite": false, "number": 8},
        "time":           {"set": true,  "infinite": false, "number": 1440},
        "cpus_per_node":  {"set": true,  "infinite": false, "number": 4},
        "cpus_per_socket":{"set": false, "infinite": false, "number": 0},
        "over_time_limit":{"set": true,  "infinite": false, "number": 30},
        "oversubscribe":  {"jobs": 4, "flags": ["force"]}
      },
      "minimums": { "nodes": 2 },
      "priority": { "job_factor": 2, "tier": 2 },
      "cpus":      { "total": 8 },
      "partition": { "state": ["DOWN"] },
      "tres": {
        "billing_weights": "CPU=1.0,Mem=0.25G",
        "configured": "cpu=8,mem=32768M,node=2,billing=8"
      },
      "grace_time": 30
    }
  ]
}`

func unmarshalPartitions(t *testing.T) slurm.SlurmPartitionsResponse {
	t.Helper()
	var resp slurm.SlurmPartitionsResponse
	if err := json.Unmarshal([]byte(sampleSlurmPartitionsJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Partitions) != 2 {
		t.Fatalf("expected 2 partitions, got %d", len(resp.Partitions))
	}
	return resp
}

// ---------------------------------------------------------------------------
// SlurmPartitionsResponse unmarshal
// ---------------------------------------------------------------------------

// TestSlurmPartitionsResponseUnmarshal validates that the raw JSON round-trips
// correctly into SlurmPartitionsResponse.
func TestSlurmPartitionsResponseUnmarshal(t *testing.T) {
	resp := unmarshalPartitions(t)
	normal := resp.Partitions[0]

	if normal.Name != "normal" {
		t.Errorf("Name: got %q, want \"normal\"", normal.Name)
	}
	if normal.Nodes.Total != 2 {
		t.Errorf("Nodes.Total: got %d, want 2", normal.Nodes.Total)
	}
	if normal.Nodes.Configured != "c[31-32]" {
		t.Errorf("Nodes.Configured: got %q, want \"c[31-32]\"", normal.Nodes.Configured)
	}
	if !normal.Maximums.Time.Infinite {
		t.Error("Maximums.Time.Infinite: expected true for INFINITE max time")
	}
	if len(normal.Partition.State) == 0 || normal.Partition.State[0] != "UP" {
		t.Errorf("Partition.State: got %v, want [\"UP\"]", normal.Partition.State)
	}
}

// ---------------------------------------------------------------------------
// MapSlurmPartitionRaw – identity and basic fields
// ---------------------------------------------------------------------------

func TestMapSlurmPartitionRaw_Identity(t *testing.T) {
	resp := unmarshalPartitions(t)
	p := slurm.MapSlurmPartitionRaw(resp.Partitions[0])

	if p.PartitionName != "normal" {
		t.Errorf("PartitionName: got %q, want \"normal\"", p.PartitionName)
	}
	if p.State != "UP" {
		t.Errorf("State: got %q, want \"UP\"", p.State)
	}
	if p.TotalNodes != 2 {
		t.Errorf("TotalNodes: got %d, want 2", p.TotalNodes)
	}
	if p.TotalCPUs != 4 {
		t.Errorf("TotalCPUs: got %d, want 4", p.TotalCPUs)
	}
	if p.NodeList != "c[31-32]" {
		t.Errorf("NodeList: got %q, want \"c[31-32]\"", p.NodeList)
	}
	if p.Nodes == nil {
		t.Error("Nodes: expected non-nil empty slice")
	}
	if len(p.Nodes) != 0 {
		t.Errorf("Nodes: got %d elements, want 0", len(p.Nodes))
	}
}

// TestMapSlurmPartitionRaw_TimeLimits checks MaxTime and DefaultTime formatting.
func TestMapSlurmPartitionRaw_TimeLimits(t *testing.T) {
	resp := unmarshalPartitions(t)

	// "normal" has INFINITE max time and unset default time.
	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.MaxTime != "INFINITE" {
		t.Errorf("normal MaxTime: got %q, want \"INFINITE\"", normal.MaxTime)
	}
	if normal.DefaultTime != "NONE" {
		t.Errorf("normal DefaultTime: got %q, want \"NONE\"", normal.DefaultTime)
	}

	// "fat" has 1-day (1440 min) max time and 1-hour (60 min) default time.
	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.MaxTime != "1-00:00:00" {
		t.Errorf("fat MaxTime: got %q, want \"1-00:00:00\"", fat.MaxTime)
	}
	if fat.DefaultTime != "01:00:00" {
		t.Errorf("fat DefaultTime: got %q, want \"01:00:00\"", fat.DefaultTime)
	}
}

// TestMapSlurmPartitionRaw_NodeCounts checks MinNodes and MaxNodes.
func TestMapSlurmPartitionRaw_NodeCounts(t *testing.T) {
	resp := unmarshalPartitions(t)

	// "normal": no max node limit (infinite), min 1.
	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.MinNodes != 1 {
		t.Errorf("normal MinNodes: got %d, want 1", normal.MinNodes)
	}
	if normal.MaxNodes != 0 {
		t.Errorf("normal MaxNodes: got %d, want 0 (unlimited)", normal.MaxNodes)
	}

	// "fat": max 8 nodes, min 2.
	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.MinNodes != 2 {
		t.Errorf("fat MinNodes: got %d, want 2", fat.MinNodes)
	}
	if fat.MaxNodes != 8 {
		t.Errorf("fat MaxNodes: got %d, want 8", fat.MaxNodes)
	}
}

// TestMapSlurmPartitionRaw_AccessControl checks AllowAccounts/DenyAccounts,
// AllowGroups, AllowQOS/DenyQOS, and AllocNodes.
func TestMapSlurmPartitionRaw_AccessControl(t *testing.T) {
	resp := unmarshalPartitions(t)

	// "normal" allows everyone.
	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.AllowAccounts != "ALL" {
		t.Errorf("normal AllowAccounts: got %q, want \"ALL\"", normal.AllowAccounts)
	}
	if normal.AllowGroups != "ALL" {
		t.Errorf("normal AllowGroups: got %q, want \"ALL\"", normal.AllowGroups)
	}
	if normal.AllowQOS != "ALL" {
		t.Errorf("normal AllowQOS: got %q, want \"ALL\"", normal.AllowQOS)
	}
	if normal.AllocNodes != "ALL" {
		t.Errorf("normal AllocNodes: got %q, want \"ALL\"", normal.AllocNodes)
	}

	// "fat" has restricted access.
	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.AllowAccounts != "research" {
		t.Errorf("fat AllowAccounts: got %q, want \"research\"", fat.AllowAccounts)
	}
	if fat.DenyAccounts != "marketing" {
		t.Errorf("fat DenyAccounts: got %q, want \"marketing\"", fat.DenyAccounts)
	}
	if fat.AllowQOS != "high" {
		t.Errorf("fat AllowQOS: got %q, want \"high\"", fat.AllowQOS)
	}
	if fat.DenyQOS != "low" {
		t.Errorf("fat DenyQOS: got %q, want \"low\"", fat.DenyQOS)
	}
	if fat.AllocNodes != "sms" {
		t.Errorf("fat AllocNodes: got %q, want \"sms\"", fat.AllocNodes)
	}
}

// TestMapSlurmPartitionRaw_Priority checks PriorityJobFactor and PriorityTier.
func TestMapSlurmPartitionRaw_Priority(t *testing.T) {
	resp := unmarshalPartitions(t)

	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.PriorityJobFactor != 1 {
		t.Errorf("normal PriorityJobFactor: got %d, want 1", normal.PriorityJobFactor)
	}
	if normal.PriorityTier != 1 {
		t.Errorf("normal PriorityTier: got %d, want 1", normal.PriorityTier)
	}

	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.PriorityJobFactor != 2 {
		t.Errorf("fat PriorityJobFactor: got %d, want 2", fat.PriorityJobFactor)
	}
	if fat.PriorityTier != 2 {
		t.Errorf("fat PriorityTier: got %d, want 2", fat.PriorityTier)
	}
}

// TestMapSlurmPartitionRaw_OverSubscribe checks OverSubscribe string formatting.
func TestMapSlurmPartitionRaw_OverSubscribe(t *testing.T) {
	resp := unmarshalPartitions(t)

	// "normal" has jobs=1, no force → "NO"
	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.OverSubscribe != "NO" {
		t.Errorf("normal OverSubscribe: got %q, want \"NO\"", normal.OverSubscribe)
	}

	// "fat" has jobs=4, force → "FORCE:4"
	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.OverSubscribe != "FORCE:4" {
		t.Errorf("fat OverSubscribe: got %q, want \"FORCE:4\"", fat.OverSubscribe)
	}
}

// TestMapSlurmPartitionRaw_CPULimits checks MaxCPUsPerNode and MaxCPUsPerSocket.
func TestMapSlurmPartitionRaw_CPULimits(t *testing.T) {
	resp := unmarshalPartitions(t)

	// "normal": no per-node CPU limit (not set) → 0
	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.MaxCPUsPerNode != 0 {
		t.Errorf("normal MaxCPUsPerNode: got %d, want 0", normal.MaxCPUsPerNode)
	}

	// "fat": 4 CPUs per node limit
	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.MaxCPUsPerNode != 4 {
		t.Errorf("fat MaxCPUsPerNode: got %d, want 4", fat.MaxCPUsPerNode)
	}
}

// TestMapSlurmPartitionRaw_OverTimeLimit checks OverTimeLimit formatting.
func TestMapSlurmPartitionRaw_OverTimeLimit(t *testing.T) {
	resp := unmarshalPartitions(t)

	// "normal": over_time_limit not set → empty string
	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.OverTimeLimit != "" {
		t.Errorf("normal OverTimeLimit: got %q, want \"\"", normal.OverTimeLimit)
	}

	// "fat": 30 minutes over-time limit → "00:30:00"
	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.OverTimeLimit != "00:30:00" {
		t.Errorf("fat OverTimeLimit: got %q, want \"00:30:00\"", fat.OverTimeLimit)
	}
}

// TestMapSlurmPartitionRaw_TRES checks TRESBillingWeights and TRES fields.
func TestMapSlurmPartitionRaw_TRES(t *testing.T) {
	resp := unmarshalPartitions(t)

	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.TRESBillingWeights != "" {
		t.Errorf("normal TRESBillingWeights: got %q, want \"\"", normal.TRESBillingWeights)
	}
	if normal.TRES != "cpu=4,mem=15872M,node=2,billing=4" {
		t.Errorf("normal TRES: got %q", normal.TRES)
	}

	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.TRESBillingWeights != "CPU=1.0,Mem=0.25G" {
		t.Errorf("fat TRESBillingWeights: got %q, want \"CPU=1.0,Mem=0.25G\"", fat.TRESBillingWeights)
	}
}

// TestMapSlurmPartitionRaw_Misc checks GraceTime, Alternate, QOS, and State.
func TestMapSlurmPartitionRaw_Misc(t *testing.T) {
	resp := unmarshalPartitions(t)

	normal := slurm.MapSlurmPartitionRaw(resp.Partitions[0])
	if normal.GraceTime != 0 {
		t.Errorf("normal GraceTime: got %d, want 0", normal.GraceTime)
	}
	if normal.Alternate != "" {
		t.Errorf("normal Alternate: got %q, want \"\"", normal.Alternate)
	}
	if normal.QOS != "" {
		t.Errorf("normal QOS: got %q, want \"\"", normal.QOS)
	}

	fat := slurm.MapSlurmPartitionRaw(resp.Partitions[1])
	if fat.GraceTime != 30 {
		t.Errorf("fat GraceTime: got %d, want 30", fat.GraceTime)
	}
	if fat.Alternate != "normal" {
		t.Errorf("fat Alternate: got %q, want \"normal\"", fat.Alternate)
	}
	if fat.State != "DOWN" {
		t.Errorf("fat State: got %q, want \"DOWN\"", fat.State)
	}
	if fat.QOS != "normal" {
		t.Errorf("fat QOS: got %q, want \"normal\"", fat.QOS)
	}
}

// TestMapSlurmPartitionRaw_DefaultFlagNotSet verifies that MapSlurmPartitionRaw
// always sets Default to false (callers must set it from sinfo output).
func TestMapSlurmPartitionRaw_DefaultFlagNotSet(t *testing.T) {
	resp := unmarshalPartitions(t)
	for _, raw := range resp.Partitions {
		p := slurm.MapSlurmPartitionRaw(raw)
		if p.Default {
			t.Errorf("partition %q: Default should be false (set by caller from sinfo)", p.PartitionName)
		}
	}
}

// TestMapSlurmPartitionRaw_EmptyStateArray verifies graceful handling when the
// partition state array is empty.
func TestMapSlurmPartitionRaw_EmptyStateArray(t *testing.T) {
	raw := slurm.SlurmPartitionRaw{
		Name:      "empty",
		Partition: slurm.SlurmPartitionState{State: []string{}},
	}
	p := slurm.MapSlurmPartitionRaw(raw)
	if p.State != "" {
		t.Errorf("State: expected empty string for empty state array, got %q", p.State)
	}
}

// TestMapSlurmPartitionRaw_YESOverSubscribe verifies "YES:N" formatting when
// jobs > 1 and no force flag.
func TestMapSlurmPartitionRaw_YESOverSubscribe(t *testing.T) {
	raw := slurm.SlurmPartitionRaw{
		Name: "shared",
		Maximums: slurm.SlurmPartitionMaximums{
			Oversubscribe: slurm.SlurmPartitionOversubscribe{Jobs: 2, Flags: []string{}},
		},
	}
	p := slurm.MapSlurmPartitionRaw(raw)
	if p.OverSubscribe != "YES:2" {
		t.Errorf("OverSubscribe: got %q, want \"YES:2\"", p.OverSubscribe)
	}
}
