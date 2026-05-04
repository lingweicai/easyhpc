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

// Job represents a user-submitted job.
// Source: sacctmgr, scontrol show jobs, slurmdbd.
type Job struct {
	JobID          string    `json:"job_id"`
	UserName       string    `json:"user_name"`
	Partition      string    `json:"partition"`
	State          string    `json:"state"` // pending, running, completed, failed, canceled, etc.
	SubmitTime     time.Time `json:"submit_time"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	AllocatedCPUs  int       `json:"allocated_cpus"`
	AllocatedNodes int       `json:"allocated_nodes"`
	ExitCode       int       `json:"exit_code"`
	CPUUsage       float64   `json:"cpu_usage"`
	MemUsage       float64   `json:"mem_usage"`
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
