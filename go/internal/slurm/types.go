// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

import "time"

// Node represents a Slurm compute node.
type Node struct {
	Name      string `json:"name"`
	State     string `json:"state"`
	Partition string `json:"partition"`
	CPUs      string `json:"cpus"`
	Memory    string `json:"memory"`
	Reason    string `json:"reason,omitempty"`
}

// Partition represents a Slurm partition.
type Partition struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	TotalCPUs  string `json:"total_cpus"`
	TotalNodes string `json:"total_nodes"`
}

// Cluster represents a Slurm cluster.
type Cluster struct {
	Name        string `json:"name"`
	ControlHost string `json:"control_host"`
	ControlPort string `json:"control_port"`
}

// Job represents a Slurm batch job.
type Job struct {
	JobID     string `json:"job_id"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Account   string `json:"account"`
	State     string `json:"state"`
	Partition string `json:"partition"`
	Nodes     string `json:"nodes"`
	Time      string `json:"time"`
	TimeLimit string `json:"time_limit"`
	Reason    string `json:"reason,omitempty"`
}

// Reservation represents a Slurm advance reservation.
type Reservation struct {
	Name      string `json:"name"`
	State     string `json:"state"`
	Nodes     string `json:"nodes"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Duration  string `json:"duration"`
	Users     string `json:"users"`
	Accounts  string `json:"accounts"`
}

// User represents a Slurm user (from sacctmgr).
type User struct {
	Name           string `json:"name"`
	DefaultAccount string `json:"default_account"`
	Admin          string `json:"admin"`
}

// Account represents a Slurm accounting account (from sacctmgr).
type Account struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Organization string `json:"organization"`
}

// LogEvent represents a parsed event from the slurmctld log file.
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
// Source: scontrol show partitions --json (Slurm ≥23.11 data_parser), slurm.conf.
//
// Field descriptions follow the Slurm scontrol partition specification:
// https://slurm.schedmd.com/scontrol.html#SECTION_PARTITIONS---SPECIFICATIONS-FOR-CREATE-AND-UPDATE-COMMANDS
type Partition struct {
	// Identity
	PartitionName string `json:"partition_name"`

	// Nodes cross-referenced from node data by partition membership.
	Nodes []Node `json:"nodes"`

	// State is the current partition state: UP, DOWN, DRAIN, or INACTIVE.
	State string `json:"state"`

	// Time limits. "INFINITE" means no limit; "NONE" means unset (use MaxTime).
	MaxTime     string `json:"max_time"`
	DefaultTime string `json:"default_time,omitempty"`

	// Node counts
	TotalNodes int `json:"total_nodes"`
	TotalCPUs  int `json:"total_cpus"`
	MinNodes   int `json:"min_nodes,omitempty"`
	MaxNodes   int `json:"max_nodes,omitempty"` // 0 means unlimited

	// Access control (empty string or "ALL" means no restriction).
	AllowGroups   string `json:"allow_groups,omitempty"`
	AllowAccounts string `json:"allow_accounts,omitempty"`
	AllowQOS      string `json:"allow_qos,omitempty"`
	DenyAccounts  string `json:"deny_accounts,omitempty"`
	DenyQOS       string `json:"deny_qos,omitempty"`
	// AllocNodes restricts which nodes may be used to launch jobs in this
	// partition (corresponds to AllocNodes in slurm.conf / scontrol).
	AllocNodes string `json:"alloc_nodes,omitempty"`

	// CPU limits per node/socket (0 means unlimited).
	MaxCPUsPerNode   int `json:"max_cpus_per_node,omitempty"`
	MaxCPUsPerSocket int `json:"max_cpus_per_socket,omitempty"`

	// Scheduling behaviour
	Default bool `json:"default"` // true when this is the default partition

	// PriorityJobFactor and PriorityTier influence relative priority of jobs.
	PriorityJobFactor int `json:"priority_job_factor,omitempty"`
	PriorityTier      int `json:"priority_tier,omitempty"`

	// OverSubscribe describes whether jobs may share nodes, e.g. "NO",
	// "YES:N", or "FORCE:N".
	OverSubscribe string `json:"over_subscribe,omitempty"`

	// PreemptMode is the preemption mode for this partition (e.g. "OFF",
	// "REQUEUE"). Note: not yet exported by Slurm's JSON data_parser.
	PreemptMode string `json:"preempt_mode,omitempty"`

	// OverTimeLimit is the maximum over-run in minutes beyond MaxTime allowed
	// before the job is killed ("INFINITE" or a minute count).
	OverTimeLimit string `json:"over_time_limit,omitempty"`

	// TRES and billing
	TRESBillingWeights string `json:"tres_billing_weights,omitempty"`
	TRES               string `json:"tres,omitempty"`

	// QOS is the Quality of Service name associated with the partition.
	QOS string `json:"qos,omitempty"`

	// GraceTime is the number of seconds a job running in this partition is
	// given to clean up after its time limit is reached.
	GraceTime int `json:"grace_time,omitempty"`

	// Alternate is the name of an alternate partition to use when the job
	// cannot run in this one.
	Alternate string `json:"alternate,omitempty"`

	// NodeList is the Slurm hostlist expression for the configured nodes
	// (e.g. "c[31-40]").
	NodeList string `json:"node_list,omitempty"`

	// MaxJobs is retained for backward compatibility; Slurm's JSON output
	// does not currently include a per-partition job count limit.
	MaxJobs int `json:"max_jobs,omitempty"`
}

// Node represents a Slurm compute or login node.
// Source: scontrol show nodes --json (via MapSlurmNodeRaw), slurmctld.
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
	State     string  `json:"state"` // running, completed, failed, etc.
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
	// Identity
	ReservationID   string `json:"reservation_id"`
	ReservationName string `json:"reservation_name,omitempty"`

	// Timing
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Duration  string     `json:"duration,omitempty"`

	// Access control
	// Users/Accounts/Groups follow Slurm's reservation format and are typically
	// comma-separated lists as emitted by scontrol.
	Users    string `json:"users,omitempty"`
	Accounts string `json:"accounts,omitempty"`
	Groups   string `json:"groups,omitempty"`

	// Resources / nodes
	Nodes         []Node `json:"nodes"`
	NodeList      string `json:"node_list,omitempty"`
	NodeCount     int    `json:"node_count,omitempty"`
	CoreCount     int    `json:"core_count,omitempty"`
	CPUs          int    `json:"cpus,omitempty"`
	PartitionName string `json:"partition_name,omitempty"`
	Features      string `json:"features,omitempty"`
	Licenses      string `json:"licenses,omitempty"`
	TRES          string `json:"tres,omitempty"`
	BurstBuffer   string `json:"burst_buffer,omitempty"`

	// Flags / behavior
	Flags         []string `json:"flags,omitempty"`
	MaxStartDelay string   `json:"max_start_delay,omitempty"`

	// Runtime state (e.g. ACTIVE, INACTIVE, COMPLETED).
	State string `json:"state,omitempty"`
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
	EventType          string    `json:"event_type"` // node_up, node_down, job_start, job_end, etc.
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
	Clusters     []Cluster     `json:"clusters"`
	Partitions   []Partition   `json:"partitions"`
	Nodes        []Node        `json:"nodes"`
	Jobs         []Job         `json:"jobs"`
	Reservations []Reservation `json:"reservations"`
	Users        []User        `json:"users"`
	Accounts     []Account     `json:"accounts"`
	LastUpdated  time.Time     `json:"last_updated"`
	Clusters     []Cluster
	Partitions   []Partition
	Nodes        []Node
	Jobs         []Job
	JobSteps     []JobStep
	Reservations []Reservation
	Users        []User
	Accounts     []Account
	SlurmDB      SlurmDBSnapshot
	SlurmLogs    []SlurmLog
	Events       []Event
	LastUpdated  time.Time
}
