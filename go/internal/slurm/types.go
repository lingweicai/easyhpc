// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

import "time"

// Cluster represents the Slurm cluster and holds configuration data,
// node lists, job queues, partitions, etc.
// Source: slurm.conf, sacctmgr, SlurmDB.
type Cluster struct {
	ClusterName   string      `json:"cluster_name"`
	SlurmdVersion string      `json:"slurmd_version"`
	SlurmVersion  string      `json:"slurm_version"`
	ControlHost   string      `json:"control_host"`
	NodeCount     int         `json:"node_count"`
	Partitions    []Partition `json:"partitions"`
}

// Partition defines a job scheduling domain (e.g. batch, debug).
// Source: scontrol show partitions, slurm.conf.
type Partition struct {
	PartitionName string `json:"partition_name"`
	Nodes         []Node `json:"nodes"`
	State         string `json:"state"`      // up, down, drained, maint, etc.
	MaxTime       string `json:"max_time"`
	TotalNodes    int    `json:"total_nodes"`
	TotalCPUs     int    `json:"total_cpus"`
	Default       bool   `json:"default"`
	MaxJobs       int    `json:"max_jobs"`
}

// Node represents a Slurm compute or login node.
// Source: scontrol show nodes, slurmctld.
type Node struct {
	NodeName   string   `json:"node_name"`
	Arch       string   `json:"arch"`
	CPUs       int      `json:"cpus"`
	Mem        int      `json:"mem"`
	State      string   `json:"state"`      // idle, allocated, down, drained, maint, etc.
	Partitions []string `json:"partitions"` // names of partitions the node belongs to
	Sockets    int      `json:"sockets"`
	FreeMem    int      `json:"free_mem"`
	CPULoad    float64  `json:"cpu_load"`
}

// JobState is a typed string for Slurm job states, matching the values
// reported by scontrol/squeue (e.g. "RUNNING", "PENDING").
type JobState string

const (
	JobStatePending   JobState = "PENDING"
	JobStateRunning   JobState = "RUNNING"
	JobStateCompleted JobState = "COMPLETED"
	JobStateFailed    JobState = "FAILED"
	JobStateCancelled JobState = "CANCELLED"
	JobStateTimeout   JobState = "TIMEOUT"
	JobStateSuspended JobState = "SUSPENDED"
	JobStatePreempted JobState = "PREEMPTED"
	JobStateNodeFail  JobState = "NODE_FAIL"
	JobStateBootFail  JobState = "BOOT_FAIL"
	JobStateRequeued  JobState = "REQUEUED"
	JobStateResizing  JobState = "RESIZING"
)

// ExitCodeSignal holds the optional signal component of a Slurm exit code.
type ExitCodeSignal struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ExitCode captures how a Slurm job or job step terminated, including the
// return code and any terminating signal.
type ExitCode struct {
	// Status is a human-readable status list, e.g. ["SUCCESS"] or ["SIGNALED"].
	Status     []string        `json:"status,omitempty"`
	ReturnCode int             `json:"return_code,omitempty"`
	Signal     *ExitCodeSignal `json:"signal,omitempty"`
}

// JobResources describes the CPU/node/task/memory counts for a job's
// requested or allocated resources.
type JobResources struct {
	CPUs        int64 `json:"cpus,omitempty"`
	Nodes       int64 `json:"nodes,omitempty"`
	Tasks       int64 `json:"tasks,omitempty"`
	MemMB       int64 `json:"mem_mb,omitempty"`
	CPUsPerTask int64 `json:"cpus_per_task,omitempty"`
}

// JobAllocatedNode describes per-node resource usage within a job.
type JobAllocatedNode struct {
	NodeName          string `json:"node_name"`
	CPUsUsed          int    `json:"cpus_used,omitempty"`
	MemoryUsedMB      int64  `json:"memory_used_mb,omitempty"`
	MemoryAllocatedMB int64  `json:"memory_allocated_mb,omitempty"`
}

// JobNodeResources holds the detailed per-node allocation reported by
// scontrol show jobs in the job_resources field.
type JobNodeResources struct {
	Nodes          string             `json:"nodes,omitempty"`
	AllocatedCores int                `json:"allocated_cores,omitempty"`
	AllocatedHosts int                `json:"allocated_hosts,omitempty"`
	AllocatedNodes []JobAllocatedNode `json:"allocated_nodes,omitempty"`
}

// Job represents a normalized user-submitted HPC job.
// It is the canonical model used both by the event bridge and the React frontend.
// Source: scontrol show jobs --json, squeue.
type Job struct {
	// Identity
	JobID     int64  `json:"job_id"`
	Cluster   string `json:"cluster,omitempty"`
	Name      string `json:"name,omitempty"`
	UserID    int64  `json:"user_id,omitempty"`
	UserName  string `json:"user_name"`
	GroupID   int64  `json:"group_id,omitempty"`
	GroupName string `json:"group_name,omitempty"`
	Account   string `json:"account,omitempty"`
	QOS       string `json:"qos,omitempty"`
	Partition string `json:"partition"`

	// Scheduling state
	State            JobState `json:"state"`
	StateReason      string   `json:"state_reason,omitempty"`
	StateDescription string   `json:"state_description,omitempty"`
	Flags            []string `json:"flags,omitempty"`
	Priority         int64    `json:"priority,omitempty"`
	Hold             bool     `json:"hold,omitempty"`
	Requeue          bool     `json:"requeue,omitempty"`
	TimeLimitMinutes int64    `json:"time_limit_minutes,omitempty"`

	// Timing – pointers because absent/unknown must be distinguishable from zero.
	SubmitTime   *time.Time `json:"submit_time,omitempty"`
	EligibleTime *time.Time `json:"eligible_time,omitempty"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`

	// Resources
	Requested    JobResources      `json:"requested"`
	Allocated    JobResources      `json:"allocated"`
	JobResources *JobNodeResources `json:"job_resources,omitempty"`

	// Placement
	Nodes          string `json:"nodes,omitempty"`
	AllocatingNode string `json:"allocating_node,omitempty"`
	BatchHost      string `json:"batch_host,omitempty"`

	// I/O paths and working directory
	Command string `json:"command,omitempty"`
	WorkDir string `json:"work_dir,omitempty"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
	StdIn   string `json:"stdin,omitempty"`

	// Exit information
	ExitCode        *ExitCode `json:"exit_code,omitempty"`
	DerivedExitCode *ExitCode `json:"derived_exit_code,omitempty"`

	// Runtime statistics (populated from sacct/sstat when available)
	CPUUsageSeconds float64 `json:"cpu_usage_seconds,omitempty"`
	MemUsageMB      float64 `json:"mem_usage_mb,omitempty"`
}

// JobStep represents a subtask of a job (sacct step) with specific resource/state.
// Source: sacctmgr, scontrol show jobsteps.
type JobStep struct {
	JobStepID string  `json:"job_step_id"`
	JobID     string  `json:"job_id"`
	StepID    string  `json:"step_id"`
	State     string  `json:"state"`      // running, completed, failed, etc.
	CPUUsage  float64 `json:"cpu_usage"`
	MemUsage  float64 `json:"mem_usage"`
	TimeSpent string  `json:"time_spent"` // duration string, e.g. "1h23m"
}

// User represents a user submitting jobs and utilizing resources.
// Source: sacctmgr list users, slurmdbd.
type User struct {
	UserID        string  `json:"user_id"`
	UserName      string  `json:"user_name"`
	TotalJobs     int     `json:"total_jobs"`
	ActiveJobs    int     `json:"active_jobs"`
	TotalCPUUsage float64 `json:"total_cpu_usage"`
	TotalMemUsage float64 `json:"total_mem_usage"`
}

// Account represents billing and resource allocation for users and groups.
// Source: sacctmgr list accounts, slurmdbd.
type Account struct {
	AccountID     string   `json:"account_id"`
	AccountName   string   `json:"account_name"`
	UserNames     []string `json:"user_names"`
	TotalJobs     int      `json:"total_jobs"`
	TotalCPUUsage float64  `json:"total_cpu_usage"`
	TotalMemUsage float64  `json:"total_mem_usage"`
}

// Reservation represents reserved resources for specific jobs, users, or accounts.
// Source: scontrol show reservations, slurmdbd.
type Reservation struct {
	ReservationID string    `json:"reservation_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Nodes         []Node    `json:"nodes"`
	CPUs          int       `json:"cpus"`
	State         string    `json:"state"` // active, expired, cancelled
}

// SlurmLog represents a log entry from Slurm's slurm.log file.
// Source: slurm.log.
type SlurmLog struct {
	LogID     string    `json:"log_id"`
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"event_type"` // e.g. job_start, job_end, node_failure
	Message   string    `json:"message"`
	Severity  string    `json:"severity"` // info, warn, error, critical
}

// Event represents a system event that occurs within Slurm, such as job state
// changes, node state changes, or system warnings.
// Source: strigger, slurmctld, slurm.log.
type Event struct {
	EventID            string    `json:"event_id"`
	EventType          string    `json:"event_type"`           // node_up, node_down, job_start, job_end, etc.
	Timestamp          time.Time `json:"timestamp"`
	AssociatedObjectID string    `json:"associated_object_id"` // job_id, node_name, etc.
	EventMessage       string    `json:"event_message"`
}

// LogEvent represents a parsed event from the slurmctld log file.
// Used internally by the log watcher for real-time log streaming.
type LogEvent struct {
	Level     string
	Message   string
	Timestamp string
}

// Cache holds the last-known state of all Slurm resources.
type Cache struct {
	Clusters     []Cluster
	Partitions   []Partition
	Nodes        []Node
	Jobs         []Job
	JobSteps     []JobStep
	Reservations []Reservation
	Users        []User
	Accounts     []Account
	SlurmLogs    []SlurmLog
	Events       []Event
	LastUpdated  time.Time
}
