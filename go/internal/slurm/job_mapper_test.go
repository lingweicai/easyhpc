// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lingweicai/easyhpc/bridge/internal/slurm"
)

// sampleSlurmJobJSON is an excerpt of the JSON produced by
// "scontrol show jobs --json" for a running 4-CPU job across 2 nodes.
const sampleSlurmJobJSON = `{
  "jobs": [
    {
      "job_id": 560,
      "cluster": "cluster",
      "name": "sleep.script",
      "user_id": 1000,
      "user_name": "dev",
      "group_id": 1000,
      "group_name": "dev",
      "account": "dev",
      "qos": "normal",
      "partition": "normal",
      "job_state": ["RUNNING"],
      "state_reason": "None",
      "state_description": "",
      "flags": ["EXACT_TASK_COUNT_REQUESTED", "JOB_WAS_RUNNING"],
      "priority": {"set": true, "infinite": false, "number": 1},
      "hold": false,
      "requeue": true,
      "time_limit": {"set": true, "infinite": false, "number": 1},
      "submit_time": {"set": true, "infinite": false, "number": 1777866658},
      "eligible_time": {"set": true, "infinite": false, "number": 1777866658},
      "start_time": {"set": true, "infinite": false, "number": 1777866658},
      "end_time": {"set": true, "infinite": false, "number": 1777866718},
      "cpus": {"set": true, "infinite": false, "number": 4},
      "node_count": {"set": true, "infinite": false, "number": 2},
      "tasks": {"set": true, "infinite": false, "number": 4},
      "cpus_per_task": {"set": true, "infinite": false, "number": 1},
      "memory_per_node": {"set": true, "infinite": false, "number": 0},
      "memory_per_cpu": {"set": false, "infinite": false, "number": 0},
      "job_resources": {
        "nodes": "c[31-32]",
        "allocated_cores": 4,
        "allocated_cpus": 0,
        "allocated_hosts": 2,
        "allocated_nodes": [
          {"nodename": "c31", "cpus_used": 2, "memory_used": 0, "memory_allocated": 1},
          {"nodename": "c32", "cpus_used": 2, "memory_used": 0, "memory_allocated": 1}
        ]
      },
      "nodes": "c[31-32]",
      "allocating_node": "sms94",
      "batch_host": "c31",
      "command": "/home/dev/sleep.script",
      "current_working_directory": "/home/dev",
      "standard_output": "/home/dev/slurm-560.out",
      "standard_error": "/home/dev/slurm-560.out",
      "standard_input": "/dev/null",
      "exit_code": {
        "status": ["SUCCESS"],
        "return_code": {"set": true, "infinite": false, "number": 0},
        "signal": {"id": {"set": false, "infinite": false, "number": 0}, "name": ""}
      },
      "derived_exit_code": {
        "status": ["SUCCESS"],
        "return_code": {"set": true, "infinite": false, "number": 0},
        "signal": {"id": {"set": false, "infinite": false, "number": 0}, "name": ""}
      },
      "tres_req_str": "cpu=4,mem=2M,node=2,billing=4",
      "tres_alloc_str": "cpu=4,mem=2M,node=2,billing=4"
    }
  ]
}`

// TestMapSlurmJobRaw_Identity checks that identity fields are mapped correctly.
func TestMapSlurmJobRaw_Identity(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(resp.Jobs))
	}

	job := slurm.MapSlurmJobRaw(resp.Jobs[0])

	if job.JobID != 560 {
		t.Errorf("JobID: got %d, want 560", job.JobID)
	}
	if job.Cluster != "cluster" {
		t.Errorf("Cluster: got %q, want \"cluster\"", job.Cluster)
	}
	if job.Name != "sleep.script" {
		t.Errorf("Name: got %q, want \"sleep.script\"", job.Name)
	}
	if job.UserName != "dev" {
		t.Errorf("UserName: got %q, want \"dev\"", job.UserName)
	}
	if job.UserID != 1000 {
		t.Errorf("UserID: got %d, want 1000", job.UserID)
	}
	if job.Account != "dev" {
		t.Errorf("Account: got %q, want \"dev\"", job.Account)
	}
	if job.QOS != "normal" {
		t.Errorf("QOS: got %q, want \"normal\"", job.QOS)
	}
	if job.Partition != "normal" {
		t.Errorf("Partition: got %q, want \"normal\"", job.Partition)
	}
}

// TestMapSlurmJobRaw_State checks state and flag fields.
func TestMapSlurmJobRaw_State(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	job := slurm.MapSlurmJobRaw(resp.Jobs[0])

	if job.State != slurm.JobStateRunning {
		t.Errorf("State: got %q, want %q", job.State, slurm.JobStateRunning)
	}
	if job.StateReason != "None" {
		t.Errorf("StateReason: got %q, want \"None\"", job.StateReason)
	}
	if len(job.Flags) != 2 {
		t.Errorf("Flags: got %d elements, want 2", len(job.Flags))
	}
	if job.Requeue != true {
		t.Errorf("Requeue: got false, want true")
	}
	if job.TimeLimitMinutes != 1 {
		t.Errorf("TimeLimitMinutes: got %d, want 1", job.TimeLimitMinutes)
	}
}

// TestMapSlurmJobRaw_Timestamps checks that Unix timestamps are converted to
// correct UTC time.Time values.
func TestMapSlurmJobRaw_Timestamps(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	job := slurm.MapSlurmJobRaw(resp.Jobs[0])

	wantSubmit := time.Unix(1777866658, 0).UTC()
	wantEnd := time.Unix(1777866718, 0).UTC()

	if job.SubmitTime == nil {
		t.Fatal("SubmitTime is nil, want non-nil")
	}
	if !job.SubmitTime.Equal(wantSubmit) {
		t.Errorf("SubmitTime: got %v, want %v", *job.SubmitTime, wantSubmit)
	}
	if job.EndTime == nil {
		t.Fatal("EndTime is nil, want non-nil")
	}
	if !job.EndTime.Equal(wantEnd) {
		t.Errorf("EndTime: got %v, want %v", *job.EndTime, wantEnd)
	}
}

// TestMapSlurmJobRaw_Resources checks that CPU/node/task resource fields are
// populated from the scontrol JSON.
func TestMapSlurmJobRaw_Resources(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	job := slurm.MapSlurmJobRaw(resp.Jobs[0])

	if job.Requested.CPUs != 4 {
		t.Errorf("Requested.CPUs: got %d, want 4", job.Requested.CPUs)
	}
	if job.Requested.Nodes != 2 {
		t.Errorf("Requested.Nodes: got %d, want 2", job.Requested.Nodes)
	}
	if job.Requested.Tasks != 4 {
		t.Errorf("Requested.Tasks: got %d, want 4", job.Requested.Tasks)
	}
	if job.Requested.CPUsPerTask != 1 {
		t.Errorf("Requested.CPUsPerTask: got %d, want 1", job.Requested.CPUsPerTask)
	}
	if job.Allocated.CPUs != 4 {
		t.Errorf("Allocated.CPUs: got %d, want 4", job.Allocated.CPUs)
	}
	if job.Allocated.Nodes != 2 {
		t.Errorf("Allocated.Nodes: got %d, want 2", job.Allocated.Nodes)
	}
}

// TestMapSlurmJobRaw_NodeResources checks the per-node allocation details.
func TestMapSlurmJobRaw_NodeResources(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	job := slurm.MapSlurmJobRaw(resp.Jobs[0])

	if job.JobResources == nil {
		t.Fatal("JobResources is nil, want non-nil")
	}
	if job.JobResources.Nodes != "c[31-32]" {
		t.Errorf("JobResources.Nodes: got %q, want \"c[31-32]\"", job.JobResources.Nodes)
	}
	if job.JobResources.AllocatedCores != 4 {
		t.Errorf("JobResources.AllocatedCores: got %d, want 4", job.JobResources.AllocatedCores)
	}
	if len(job.JobResources.AllocatedNodes) != 2 {
		t.Fatalf("JobResources.AllocatedNodes: got %d, want 2", len(job.JobResources.AllocatedNodes))
	}

	c31 := job.JobResources.AllocatedNodes[0]
	if c31.NodeName != "c31" {
		t.Errorf("AllocatedNodes[0].NodeName: got %q, want \"c31\"", c31.NodeName)
	}
	if c31.CPUsUsed != 2 {
		t.Errorf("AllocatedNodes[0].CPUsUsed: got %d, want 2", c31.CPUsUsed)
	}
}

// TestMapSlurmJobRaw_Placement checks placement fields (nodes, batch host, etc.).
func TestMapSlurmJobRaw_Placement(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	job := slurm.MapSlurmJobRaw(resp.Jobs[0])

	if job.Nodes != "c[31-32]" {
		t.Errorf("Nodes: got %q, want \"c[31-32]\"", job.Nodes)
	}
	if job.AllocatingNode != "sms94" {
		t.Errorf("AllocatingNode: got %q, want \"sms94\"", job.AllocatingNode)
	}
	if job.BatchHost != "c31" {
		t.Errorf("BatchHost: got %q, want \"c31\"", job.BatchHost)
	}
}

// TestMapSlurmJobRaw_IO checks I/O path and working directory fields.
func TestMapSlurmJobRaw_IO(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	job := slurm.MapSlurmJobRaw(resp.Jobs[0])

	if job.Command != "/home/dev/sleep.script" {
		t.Errorf("Command: got %q", job.Command)
	}
	if job.WorkDir != "/home/dev" {
		t.Errorf("WorkDir: got %q", job.WorkDir)
	}
	if job.Output != "/home/dev/slurm-560.out" {
		t.Errorf("Output: got %q", job.Output)
	}
	if job.Error != "/home/dev/slurm-560.out" {
		t.Errorf("Error: got %q", job.Error)
	}
	if job.StdIn != "/dev/null" {
		t.Errorf("StdIn: got %q", job.StdIn)
	}
}

// TestMapSlurmJobRaw_ExitCode checks exit code mapping.
func TestMapSlurmJobRaw_ExitCode(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	job := slurm.MapSlurmJobRaw(resp.Jobs[0])

	if job.ExitCode == nil {
		t.Fatal("ExitCode is nil, want non-nil")
	}
	if len(job.ExitCode.Status) == 0 || job.ExitCode.Status[0] != "SUCCESS" {
		t.Errorf("ExitCode.Status: got %v, want [\"SUCCESS\"]", job.ExitCode.Status)
	}
	if job.ExitCode.ReturnCode != 0 {
		t.Errorf("ExitCode.ReturnCode: got %d, want 0", job.ExitCode.ReturnCode)
	}
	// Signal ID is not set, so Signal should be nil.
	if job.ExitCode.Signal != nil {
		t.Errorf("ExitCode.Signal: got %v, want nil", job.ExitCode.Signal)
	}
}

// TestMapSlurmJobRaw_ZeroTimestampIsNil verifies that a zero Unix timestamp
// (number: 0, which Slurm uses to mean "not set") maps to nil.
func TestMapSlurmJobRaw_ZeroTimestampIsNil(t *testing.T) {
	raw := slurm.SlurmJobRaw{
		JobID:     1,
		UserName:  "test",
		Partition: "normal",
		JobState:  []string{"PENDING"},
		// SubmitTime left at zero value → should be nil
	}
	job := slurm.MapSlurmJobRaw(raw)
	if job.SubmitTime != nil {
		t.Errorf("SubmitTime: expected nil for zero timestamp, got %v", *job.SubmitTime)
	}
	if job.StartTime != nil {
		t.Errorf("StartTime: expected nil for zero timestamp, got %v", *job.StartTime)
	}
}

// TestMapSlurmJobRaw_EmptyJobState verifies graceful handling of an empty
// job_state array.
func TestMapSlurmJobRaw_EmptyJobState(t *testing.T) {
	raw := slurm.SlurmJobRaw{
		JobID:     2,
		UserName:  "test",
		Partition: "normal",
		JobState:  []string{},
	}
	job := slurm.MapSlurmJobRaw(raw)
	if job.State != "" {
		t.Errorf("State: expected empty string for empty job_state, got %q", job.State)
	}
}

// TestSlurmJobsResponseUnmarshal validates round-trip JSON unmarshal of the
// SlurmJobsResponse type against the sample JSON.
func TestSlurmJobsResponseUnmarshal(t *testing.T) {
	var resp slurm.SlurmJobsResponse
	if err := json.Unmarshal([]byte(sampleSlurmJobJSON), &resp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if len(resp.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(resp.Jobs))
	}
	raw := resp.Jobs[0]
	if raw.JobID != 560 {
		t.Errorf("raw.JobID: got %d, want 560", raw.JobID)
	}
	if !raw.CPUs.Set {
		t.Error("raw.CPUs.Set: expected true")
	}
	if raw.CPUs.Number != 4 {
		t.Errorf("raw.CPUs.Number: got %v, want 4", raw.CPUs.Number)
	}
}
