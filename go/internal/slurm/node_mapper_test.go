// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm_test

import (
	"encoding/json"
	"testing"

	"github.com/lingweicai/easyhpc/bridge/internal/slurm"
)

// sampleSlurmNodesJSON is an excerpt of the JSON produced by
// "scontrol show nodes --json" containing two nodes with different partitions.
const sampleSlurmNodesJSON = `{
  "nodes": [
    {
      "architecture": "x86_64",
      "boards": 1,
      "boot_time": {"set": true, "infinite": false, "number": 1777941810},
      "cores": 1,
      "cpu_binding": 0,
      "cpu_load": 0,
      "free_mem": {"set": true, "infinite": false, "number": 6330},
      "cpus": 2,
      "effective_cpus": 2,
      "last_busy": {"set": true, "infinite": false, "number": 1777941838},
      "name": "c31",
      "state": ["IDLE"],
      "partitions": ["normal"],
      "port": 6818,
      "real_memory": 7936,
      "alloc_memory": 0,
      "alloc_cpus": 0,
      "slurmd_start_time": {"set": true, "infinite": false, "number": 1777941758},
      "sockets": 2,
      "threads": 1,
      "weight": 1,
      "tres": "cpu=2,mem=7936M,billing=2",
      "version": "23.11.10"
    },
    {
      "architecture": "x86_64",
      "boards": 1,
      "boot_time": {"set": true, "infinite": false, "number": 1777941816},
      "cores": 1,
      "cpu_binding": 0,
      "cpu_load": 0.5,
      "free_mem": {"set": false, "infinite": false, "number": 0},
      "cpus": 4,
      "effective_cpus": 4,
      "last_busy": {"set": true, "infinite": false, "number": 1777941836},
      "name": "c41",
      "state": ["ALLOCATED"],
      "partitions": ["fat"],
      "port": 6818,
      "real_memory": 16384,
      "alloc_memory": 0,
      "alloc_cpus": 4,
      "slurmd_start_time": {"set": true, "infinite": false, "number": 1777941757},
      "sockets": 4,
      "threads": 1,
      "weight": 1,
      "tres": "cpu=4,mem=16384M,billing=4",
      "version": "23.11.10"
    }
  ]
}`

// TestMapSlurmNodeRaw_BasicFields verifies that core node fields are mapped
// correctly from the scontrol JSON.
func TestMapSlurmNodeRaw_BasicFields(t *testing.T) {
	var resp slurm.SlurmNodesResponse
	if err := json.Unmarshal([]byte(sampleSlurmNodesJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(resp.Nodes))
	}

	node := slurm.MapSlurmNodeRaw(resp.Nodes[0])

	if node.NodeName != "c31" {
		t.Errorf("NodeName: got %q, want \"c31\"", node.NodeName)
	}
	if node.Arch != "x86_64" {
		t.Errorf("Arch: got %q, want \"x86_64\"", node.Arch)
	}
	if node.CPUs != 2 {
		t.Errorf("CPUs: got %d, want 2", node.CPUs)
	}
	if node.Sockets != 2 {
		t.Errorf("Sockets: got %d, want 2", node.Sockets)
	}
	if node.Mem != 7936 {
		t.Errorf("Mem: got %d, want 7936", node.Mem)
	}
}

// TestMapSlurmNodeRaw_State verifies that the state array is reduced to a
// single string (first element).
func TestMapSlurmNodeRaw_State(t *testing.T) {
	var resp slurm.SlurmNodesResponse
	if err := json.Unmarshal([]byte(sampleSlurmNodesJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	node := slurm.MapSlurmNodeRaw(resp.Nodes[0])
	if node.State != "IDLE" {
		t.Errorf("State: got %q, want \"IDLE\"", node.State)
	}

	node2 := slurm.MapSlurmNodeRaw(resp.Nodes[1])
	if node2.State != "ALLOCATED" {
		t.Errorf("State: got %q, want \"ALLOCATED\"", node2.State)
	}
}

// TestMapSlurmNodeRaw_Partitions verifies partition list mapping.
func TestMapSlurmNodeRaw_Partitions(t *testing.T) {
	var resp slurm.SlurmNodesResponse
	if err := json.Unmarshal([]byte(sampleSlurmNodesJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	node := slurm.MapSlurmNodeRaw(resp.Nodes[0])
	if len(node.Partitions) != 1 || node.Partitions[0] != "normal" {
		t.Errorf("Partitions: got %v, want [\"normal\"]", node.Partitions)
	}

	node2 := slurm.MapSlurmNodeRaw(resp.Nodes[1])
	if len(node2.Partitions) != 1 || node2.Partitions[0] != "fat" {
		t.Errorf("Partitions: got %v, want [\"fat\"]", node2.Partitions)
	}
}

// TestMapSlurmNodeRaw_FreeMem_Set verifies that FreeMem is populated when the
// optional envelope has Set=true and Infinite=false.
func TestMapSlurmNodeRaw_FreeMem_Set(t *testing.T) {
	var resp slurm.SlurmNodesResponse
	if err := json.Unmarshal([]byte(sampleSlurmNodesJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	node := slurm.MapSlurmNodeRaw(resp.Nodes[0])
	if node.FreeMem != 6330 {
		t.Errorf("FreeMem: got %d, want 6330", node.FreeMem)
	}
}

// TestMapSlurmNodeRaw_FreeMem_NotSet verifies that FreeMem is 0 when the
// optional envelope has Set=false.
func TestMapSlurmNodeRaw_FreeMem_NotSet(t *testing.T) {
	var resp slurm.SlurmNodesResponse
	if err := json.Unmarshal([]byte(sampleSlurmNodesJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	node := slurm.MapSlurmNodeRaw(resp.Nodes[1])
	if node.FreeMem != 0 {
		t.Errorf("FreeMem: got %d, want 0 (not set)", node.FreeMem)
	}
}

// TestMapSlurmNodeRaw_FreeMem_Infinite verifies that FreeMem is 0 when the
// optional envelope has Infinite=true.
func TestMapSlurmNodeRaw_FreeMem_Infinite(t *testing.T) {
	raw := slurm.SlurmNodeRaw{
		Name:   "test",
		State:  []string{"IDLE"},
		FreeMem: slurm.SlurmOptionalInt{Set: true, Infinite: true, Number: 999},
	}
	node := slurm.MapSlurmNodeRaw(raw)
	if node.FreeMem != 0 {
		t.Errorf("FreeMem: got %d, want 0 (infinite)", node.FreeMem)
	}
}

// TestMapSlurmNodeRaw_EmptyState verifies graceful handling of an empty state
// array.
func TestMapSlurmNodeRaw_EmptyState(t *testing.T) {
	raw := slurm.SlurmNodeRaw{
		Name:  "test",
		State: []string{},
	}
	node := slurm.MapSlurmNodeRaw(raw)
	if node.State != "" {
		t.Errorf("State: expected empty string for empty state array, got %q", node.State)
	}
}

// TestSlurmNodesResponseUnmarshal validates round-trip JSON unmarshal of the
// SlurmNodesResponse type against the sample JSON.
func TestSlurmNodesResponseUnmarshal(t *testing.T) {
	var resp slurm.SlurmNodesResponse
	if err := json.Unmarshal([]byte(sampleSlurmNodesJSON), &resp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if len(resp.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(resp.Nodes))
	}
	raw := resp.Nodes[0]
	if raw.Name != "c31" {
		t.Errorf("raw.Name: got %q, want \"c31\"", raw.Name)
	}
	if !raw.FreeMem.Set {
		t.Error("raw.FreeMem.Set: expected true")
	}
	if raw.FreeMem.Number != 6330 {
		t.Errorf("raw.FreeMem.Number: got %v, want 6330", raw.FreeMem.Number)
	}
}
