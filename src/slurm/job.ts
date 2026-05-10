/*
 * SPDX-License-Identifier: LGPL-2.1-or-later
 */

/**
 * Slurm job model for the React frontend.
 *
 * These types mirror the normalised Go `Job` model emitted by the
 * easyhpc-bridge over the Cockpit JSON protocol.  Timestamps are ISO 8601
 * strings (the bridge serialises Go `time.Time` values as RFC 3339).
 *
 * Backend/frontend field mapping
 * ──────────────────────────────
 * Go field              │ JSON key              │ TS field
 * ─────────────────────────────────────────────────────────
 * Job.JobID             │ job_id                │ jobId
 * Job.Cluster           │ cluster               │ cluster
 * Job.Name              │ name                  │ name
 * Job.UserID            │ user_id               │ userId
 * Job.UserName          │ user_name             │ userName
 * Job.GroupID           │ group_id              │ groupId
 * Job.GroupName         │ group_name            │ groupName
 * Job.Account           │ account               │ account
 * Job.QOS               │ qos                   │ qos
 * Job.Partition         │ partition             │ partition
 * Job.State             │ state                 │ state
 * Job.StateReason       │ state_reason          │ stateReason
 * Job.StateDescription  │ state_description     │ stateDescription
 * Job.Flags             │ flags                 │ flags
 * Job.Priority          │ priority              │ priority
 * Job.Hold              │ hold                  │ hold
 * Job.Requeue           │ requeue               │ requeue
 * Job.TimeLimitMinutes  │ time_limit_minutes    │ timeLimitMinutes
 * Job.SubmitTime        │ submit_time           │ submitTime (ISO string)
 * Job.EligibleTime      │ eligible_time         │ eligibleTime (ISO string)
 * Job.StartTime         │ start_time            │ startTime (ISO string)
 * Job.EndTime           │ end_time              │ endTime (ISO string)
 * Job.Requested         │ requested             │ requested
 * Job.Allocated         │ allocated             │ allocated
 * Job.JobResources      │ job_resources         │ jobResources
 * Job.Nodes             │ nodes                 │ nodes
 * Job.AllocatingNode    │ allocating_node       │ allocatingNode
 * Job.BatchHost         │ batch_host            │ batchHost
 * Job.Command           │ command               │ command
 * Job.WorkDir           │ work_dir              │ workDir
 * Job.Output            │ output                │ output
 * Job.Error             │ error                 │ error
 * Job.StdIn             │ stdin                 │ stdIn
 * Job.ExitCode          │ exit_code             │ exitCode
 * Job.DerivedExitCode   │ derived_exit_code     │ derivedExitCode
 * Job.CPUUsageSeconds   │ cpu_usage_seconds     │ cpuUsageSeconds
 * Job.MemUsageMB        │ mem_usage_mb          │ memUsageMb
 */

// ---------------------------------------------------------------------------
// JobState – matches Go JobState constants
// ---------------------------------------------------------------------------

/** All known Slurm job state values. */
export type JobState =
    | "PENDING"
    | "RUNNING"
    | "COMPLETED"
    | "FAILED"
    | "CANCELLED"
    | "TIMEOUT"
    | "SUSPENDED"
    | "PREEMPTED"
    | "NODE_FAIL"
    | "BOOT_FAIL"
    | "REQUEUED"
    | "RESIZING"
    | string; // allow unknown future states without breaking the type

// ---------------------------------------------------------------------------
// Exit code
// ---------------------------------------------------------------------------

/** Terminating signal information within a job exit code. */
export interface ExitCodeSignal {
    id?: number;
    name?: string;
}

/**
 * How a Slurm job or job step terminated.
 * Corresponds to Go ExitCode.
 */
export interface ExitCode {
    /** Human-readable status labels, e.g. ["SUCCESS"] or ["SIGNALED"]. */
    status?: string[];
    returnCode?: number;
    signal?: ExitCodeSignal;
}

// ---------------------------------------------------------------------------
// Resources
// ---------------------------------------------------------------------------

/**
 * CPU/node/task/memory counts for either requested or allocated resources.
 * Corresponds to Go JobResources.
 */
export interface JobResources {
    cpus?: number;
    nodes?: number;
    tasks?: number;
    /** Memory in megabytes. */
    memMb?: number;
    cpusPerTask?: number;
}

/**
 * Per-node resource usage within an allocated job.
 * Corresponds to Go JobAllocatedNode.
 */
export interface JobAllocatedNode {
    nodeName: string;
    cpusUsed?: number;
    /** Memory currently in use on this node, in MB. */
    memoryUsedMb?: number;
    /** Memory allocated to this job on this node, in MB. */
    memoryAllocatedMb?: number;
}

/**
 * Detailed per-node allocation from scontrol job_resources.
 * Corresponds to Go JobNodeResources.
 */
export interface JobNodeResources {
    /** Slurm node expression, e.g. "c[31-32]". */
    nodes?: string;
    allocatedCores?: number;
    allocatedHosts?: number;
    allocatedNodes?: JobAllocatedNode[];
}

// ---------------------------------------------------------------------------
// Job – canonical frontend model
// ---------------------------------------------------------------------------

/**
 * A normalised Slurm job as delivered by the easyhpc-bridge.
 *
 * All timestamp fields are ISO 8601 strings (RFC 3339) or undefined when
 * the timestamp is absent (e.g. endTime for a running job).
 *
 * Corresponds to Go Job.
 */
export interface Job {
    // Identity
    jobId: number;
    cluster?: string;
    name?: string;
    userId?: number;
    userName: string;
    groupId?: number;
    groupName?: string;
    account?: string;
    qos?: string;
    partition: string;

    // Scheduling state
    state: JobState;
    stateReason?: string;
    stateDescription?: string;
    flags?: string[];
    priority?: number;
    hold?: boolean;
    requeue?: boolean;
    timeLimitMinutes?: number;

    // Timing (ISO 8601 strings)
    submitTime?: string;
    eligibleTime?: string;
    startTime?: string;
    endTime?: string;

    // Resources
    requested: JobResources;
    allocated: JobResources;
    jobResources?: JobNodeResources;

    // Placement
    /** Slurm node expression for all allocated nodes, e.g. "c[31-32]". */
    nodes?: string;
    allocatingNode?: string;
    batchHost?: string;

    // I/O
    command?: string;
    workDir?: string;
    output?: string;
    error?: string;
    stdIn?: string;

    // Exit information
    exitCode?: ExitCode;
    derivedExitCode?: ExitCode;

    // Runtime statistics (from sacct/sstat when available)
    cpuUsageSeconds?: number;
    memUsageMb?: number;
}
