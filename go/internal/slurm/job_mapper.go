// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

import "time"

// ---------------------------------------------------------------------------
// Raw Slurm JSON shapes (scontrol show jobs --json)
//
// These types mirror the JSON produced by Slurm's data_parser plugin (≥23.11).
// They are used only for decoding; callers should convert to the normalised
// Job model via MapSlurmJobRaw.
// ---------------------------------------------------------------------------

// SlurmOptionalInt models Slurm's "optional number" JSON shape:
//
//	{"set": true, "infinite": false, "number": N}
//
// When Set is false the value is absent/unknown; when Infinite is true the
// value is unlimited.  Number is float64 to accommodate fractional billing
// TRES values.
type SlurmOptionalInt struct {
	Set      bool    `json:"set"`
	Infinite bool    `json:"infinite"`
	Number   float64 `json:"number"`
}

// SlurmExitCodeSignalRaw models the signal sub-object inside an exit_code.
type SlurmExitCodeSignalRaw struct {
	ID   SlurmOptionalInt `json:"id"`
	Name string           `json:"name"`
}

// SlurmExitCodeRaw models the exit_code / derived_exit_code JSON object
// returned by scontrol show jobs --json.
type SlurmExitCodeRaw struct {
	Status     []string               `json:"status"`
	ReturnCode SlurmOptionalInt       `json:"return_code"`
	Signal     SlurmExitCodeSignalRaw `json:"signal"`
}

// SlurmAllocatedNodeRaw models one entry inside job_resources.allocated_nodes.
type SlurmAllocatedNodeRaw struct {
	// Sockets is an opaque map and intentionally left as a raw value;
	// callers that need socket/core detail can inspect it separately.
	NodeName        string `json:"nodename"`
	CPUsUsed        int    `json:"cpus_used"`
	MemoryUsed      int64  `json:"memory_used"`
	MemoryAllocated int64  `json:"memory_allocated"`
}

// SlurmJobResourcesRaw models the job_resources JSON object.
type SlurmJobResourcesRaw struct {
	Nodes          string                  `json:"nodes"`
	AllocatedCores int                     `json:"allocated_cores"`
	AllocatedCPUs  int                     `json:"allocated_cpus"`
	AllocatedHosts int                     `json:"allocated_hosts"`
	AllocatedNodes []SlurmAllocatedNodeRaw `json:"allocated_nodes"`
}

// SlurmJobRaw models the JSON shape of a single job entry from
// "scontrol show jobs --json".  Fields that Slurm wraps in the
// optional-number envelope are decoded as SlurmOptionalInt.
type SlurmJobRaw struct {
	// Identity
	JobID     int64  `json:"job_id"`
	Cluster   string `json:"cluster"`
	Name      string `json:"name"`
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
	Account   string `json:"account"`
	QOS       string `json:"qos"`
	Partition string `json:"partition"`

	// State
	JobState         []string `json:"job_state"`
	StateDescription string   `json:"state_description"`
	StateReason      string   `json:"state_reason"`
	Flags            []string `json:"flags"`

	// Scheduling
	Priority SlurmOptionalInt `json:"priority"`
	Hold     bool             `json:"hold"`
	Requeue  bool             `json:"requeue"`
	TimeLimit SlurmOptionalInt `json:"time_limit"`

	// Timing (Unix timestamps wrapped in optional-number envelopes)
	SubmitTime   SlurmOptionalInt `json:"submit_time"`
	EligibleTime SlurmOptionalInt `json:"eligible_time"`
	StartTime    SlurmOptionalInt `json:"start_time"`
	EndTime      SlurmOptionalInt `json:"end_time"`

	// Requested resources
	CPUs        SlurmOptionalInt `json:"cpus"`
	NodeCount   SlurmOptionalInt `json:"node_count"`
	Tasks       SlurmOptionalInt `json:"tasks"`
	CPUsPerTask SlurmOptionalInt `json:"cpus_per_task"`
	// Memory is reported per-CPU or per-node; both are captured.
	MemoryPerNode SlurmOptionalInt `json:"memory_per_node"`
	MemoryPerCPU  SlurmOptionalInt `json:"memory_per_cpu"`

	// Allocated node details
	JobResources *SlurmJobResourcesRaw `json:"job_resources"`

	// Placement
	Nodes          string `json:"nodes"`
	AllocatingNode string `json:"allocating_node"`
	BatchHost      string `json:"batch_host"`

	// I/O
	Command                 string `json:"command"`
	CurrentWorkingDirectory string `json:"current_working_directory"`
	StandardOutput          string `json:"standard_output"`
	StandardError           string `json:"standard_error"`
	StandardInput           string `json:"standard_input"`

	// Exit codes
	ExitCode        SlurmExitCodeRaw `json:"exit_code"`
	DerivedExitCode SlurmExitCodeRaw `json:"derived_exit_code"`

	// TRES strings (informational)
	TresReqStr   string `json:"tres_req_str"`
	TresAllocStr string `json:"tres_alloc_str"`
}

// SlurmJobsResponse models the top-level JSON object returned by
// "scontrol show jobs --json".
type SlurmJobsResponse struct {
	Jobs []SlurmJobRaw `json:"jobs"`
}

// ---------------------------------------------------------------------------
// Mapping helpers
// ---------------------------------------------------------------------------

// mapTimestamp converts a SlurmOptionalInt Unix timestamp to *time.Time.
// Returns nil when the value is unset, infinite, or zero (Slurm uses 0 to
// mean "not yet set", e.g. end_time for a running job).
func mapTimestamp(opt SlurmOptionalInt) *time.Time {
	if !opt.Set || opt.Infinite || opt.Number == 0 {
		return nil
	}
	t := time.Unix(int64(opt.Number), 0).UTC()
	return &t
}

// mapExitCode converts a SlurmExitCodeRaw into an *ExitCode.
// Returns nil when the status slice is empty (no exit information present).
func mapExitCode(raw SlurmExitCodeRaw) *ExitCode {
	if len(raw.Status) == 0 {
		return nil
	}
	ec := &ExitCode{
		Status:     raw.Status,
		ReturnCode: int(raw.ReturnCode.Number),
	}
	if raw.Signal.ID.Set && raw.Signal.ID.Number > 0 {
		ec.Signal = &ExitCodeSignal{
			ID:   int(raw.Signal.ID.Number),
			Name: raw.Signal.Name,
		}
	}
	return ec
}

// MapSlurmJobRaw converts a SlurmJobRaw (decoded from scontrol show jobs
// --json) into the normalised Job model used by the API and the frontend.
//
// Requested and Allocated resource counts are both derived from the scontrol
// JSON.  Slurm reports requested = allocated for running jobs; for pending
// jobs the allocated counts may be zero.
func MapSlurmJobRaw(raw SlurmJobRaw) Job {
	// Use the first element of job_state as the canonical state value.
	state := JobState("")
	if len(raw.JobState) > 0 {
		state = JobState(raw.JobState[0])
	}

	// Build per-node allocation details when job_resources is present.
	var jobResources *JobNodeResources
	if raw.JobResources != nil {
		nodes := make([]JobAllocatedNode, 0, len(raw.JobResources.AllocatedNodes))
		for _, n := range raw.JobResources.AllocatedNodes {
			nodes = append(nodes, JobAllocatedNode{
				NodeName:          n.NodeName,
				CPUsUsed:          n.CPUsUsed,
				MemoryUsedMB:      n.MemoryUsed,
				MemoryAllocatedMB: n.MemoryAllocated,
			})
		}
		jobResources = &JobNodeResources{
			Nodes:          raw.JobResources.Nodes,
			AllocatedCores: raw.JobResources.AllocatedCores,
			AllocatedHosts: raw.JobResources.AllocatedHosts,
			AllocatedNodes: nodes,
		}
	}

	// Compute memory in MB: prefer per-node, fall back to per-cpu.
	memMB := int64(raw.MemoryPerNode.Number)
	if memMB == 0 && raw.MemoryPerCPU.Set {
		memMB = int64(raw.MemoryPerCPU.Number) * int64(raw.CPUs.Number)
	}

	return Job{
		// Identity
		JobID:     raw.JobID,
		Cluster:   raw.Cluster,
		Name:      raw.Name,
		UserID:    raw.UserID,
		UserName:  raw.UserName,
		GroupID:   raw.GroupID,
		GroupName: raw.GroupName,
		Account:   raw.Account,
		QOS:       raw.QOS,
		Partition: raw.Partition,

		// Scheduling state
		State:            state,
		StateReason:      raw.StateReason,
		StateDescription: raw.StateDescription,
		Flags:            raw.Flags,
		Priority:         int64(raw.Priority.Number),
		Hold:             raw.Hold,
		Requeue:          raw.Requeue,
		TimeLimitMinutes: int64(raw.TimeLimit.Number),

		// Timing
		SubmitTime:   mapTimestamp(raw.SubmitTime),
		EligibleTime: mapTimestamp(raw.EligibleTime),
		StartTime:    mapTimestamp(raw.StartTime),
		EndTime:      mapTimestamp(raw.EndTime),

		// Resources – for scontrol JSON requested == allocated for running jobs.
		Requested: JobResources{
			CPUs:        int64(raw.CPUs.Number),
			Nodes:       int64(raw.NodeCount.Number),
			Tasks:       int64(raw.Tasks.Number),
			CPUsPerTask: int64(raw.CPUsPerTask.Number),
			MemMB:       memMB,
		},
		Allocated: JobResources{
			CPUs:  int64(raw.CPUs.Number),
			Nodes: int64(raw.NodeCount.Number),
			Tasks: int64(raw.Tasks.Number),
		},
		JobResources: jobResources,

		// Placement
		Nodes:          raw.Nodes,
		AllocatingNode: raw.AllocatingNode,
		BatchHost:      raw.BatchHost,

		// I/O
		Command: raw.Command,
		WorkDir: raw.CurrentWorkingDirectory,
		Output:  raw.StandardOutput,
		Error:   raw.StandardError,
		StdIn:   raw.StandardInput,

		// Exit information
		ExitCode:        mapExitCode(raw.ExitCode),
		DerivedExitCode: mapExitCode(raw.DerivedExitCode),
	}
}
