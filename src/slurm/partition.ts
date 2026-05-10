/*
 * SPDX-License-Identifier: LGPL-2.1-or-later
 */

/**
 * Slurm partition model for the React frontend.
 *
 * These types mirror the normalised Go `Partition` model emitted by the
 * easyhpc-bridge over the Cockpit JSON protocol.
 *
 * Field descriptions follow the Slurm scontrol partition specification:
 * https://slurm.schedmd.com/scontrol.html#SECTION_PARTITIONS---SPECIFICATIONS-FOR-CREATE-AND-UPDATE-COMMANDS
 *
 * Backend/frontend field mapping
 * ──────────────────────────────
 * Go field                   │ JSON key                  │ TS field
 * ───────────────────────────────────────────────────────────────────────
 * Partition.PartitionName    │ partition_name            │ partitionName
 * Partition.Nodes            │ nodes                     │ nodes
 * Partition.State            │ state                     │ state
 * Partition.MaxTime          │ max_time                  │ maxTime
 * Partition.DefaultTime      │ default_time              │ defaultTime
 * Partition.TotalNodes       │ total_nodes               │ totalNodes
 * Partition.TotalCPUs        │ total_cpus                │ totalCpus
 * Partition.MinNodes         │ min_nodes                 │ minNodes
 * Partition.MaxNodes         │ max_nodes                 │ maxNodes
 * Partition.AllowGroups      │ allow_groups              │ allowGroups
 * Partition.AllowAccounts    │ allow_accounts            │ allowAccounts
 * Partition.AllowQOS         │ allow_qos                 │ allowQos
 * Partition.DenyAccounts     │ deny_accounts             │ denyAccounts
 * Partition.DenyQOS          │ deny_qos                  │ denyQos
 * Partition.AllocNodes       │ alloc_nodes               │ allocNodes
 * Partition.MaxCPUsPerNode   │ max_cpus_per_node         │ maxCpusPerNode
 * Partition.MaxCPUsPerSocket │ max_cpus_per_socket       │ maxCpusPerSocket
 * Partition.Default          │ default                   │ default
 * Partition.PriorityJobFactor│ priority_job_factor       │ priorityJobFactor
 * Partition.PriorityTier     │ priority_tier             │ priorityTier
 * Partition.OverSubscribe    │ over_subscribe            │ overSubscribe
 * Partition.PreemptMode      │ preempt_mode              │ preemptMode
 * Partition.OverTimeLimit    │ over_time_limit           │ overTimeLimit
 * Partition.TRESBillingWeights│tres_billing_weights      │ tresBillingWeights
 * Partition.TRES             │ tres                      │ tres
 * Partition.QOS              │ qos                       │ qos
 * Partition.GraceTime        │ grace_time                │ graceTime
 * Partition.Alternate        │ alternate                 │ alternate
 * Partition.NodeList         │ node_list                 │ nodeList
 * Partition.MaxJobs          │ max_jobs                  │ maxJobs
 */

import type { Node } from "./node";

// ---------------------------------------------------------------------------
// PartitionState – matches Slurm partition state values
// ---------------------------------------------------------------------------

/**
 * Current operational state of a Slurm partition.
 * Corresponds to the `state_up` field in partition_info_t.
 */
export type PartitionState =
    | "UP"
    | "DOWN"
    | "DRAIN"
    | "INACTIVE"
    | string; // allow unknown states without breaking the type

// ---------------------------------------------------------------------------
// Partition – canonical frontend model
// ---------------------------------------------------------------------------

/**
 * A normalised Slurm partition as delivered by the easyhpc-bridge.
 *
 * Time-limit fields (maxTime, defaultTime, overTimeLimit) are Slurm-formatted
 * strings such as "INFINITE", "NONE", or "D-HH:MM:SS".
 *
 * Corresponds to Go Partition.
 */
export interface Partition {
    // Identity
    partitionName: string;

    /** Compute nodes belonging to this partition (cross-referenced from node data). */
    nodes: Node[];

    // State
    /** Current partition state: UP, DOWN, DRAIN, or INACTIVE. */
    state: PartitionState;

    // Time limits ("INFINITE" = no limit, "NONE" = unset / use MaxTime)
    /** Maximum wall-clock time for jobs in this partition (Slurm time string). */
    maxTime: string;
    /** Default wall-clock time when a job does not specify --time. */
    defaultTime?: string;

    // Node counts
    totalNodes: number;
    totalCpus: number;
    /** Minimum number of nodes a job must request. */
    minNodes?: number;
    /** Maximum number of nodes a job may request; 0 means unlimited. */
    maxNodes?: number;

    // Access control (empty string or "ALL" means no restriction)
    /** Unix groups allowed to submit jobs; see AllowGroups in slurm.conf. */
    allowGroups?: string;
    /** Accounts allowed to submit jobs; see AllowAccounts in slurm.conf. */
    allowAccounts?: string;
    /** QOS names allowed in this partition. */
    allowQos?: string;
    /** Accounts denied access to this partition. */
    denyAccounts?: string;
    /** QOS names denied in this partition. */
    denyQos?: string;
    /**
     * Nodes on which jobs may be allocated; restricts which submit hosts can
     * run jobs here (AllocNodes in slurm.conf).
     */
    allocNodes?: string;

    // CPU limits per node/socket (0 means unlimited)
    maxCpusPerNode?: number;
    maxCpusPerSocket?: number;

    // Scheduling behaviour
    /** True when this is the default partition (marked with '*' in sinfo). */
    default: boolean;
    /** Relative priority weight for jobs in this partition (PriorityJobFactor). */
    priorityJobFactor?: number;
    /** Tier used to order partitions when calculating job priority (PriorityTier). */
    priorityTier?: number;
    /**
     * Node sharing policy: "NO", "YES:N", or "FORCE:N".
     * Corresponds to OverSubscribe in slurm.conf.
     */
    overSubscribe?: string;
    /**
     * Preemption mode, e.g. "OFF", "REQUEUE", "CANCEL", "SUSPEND".
     * Not yet exported by Slurm's JSON data_parser (always empty).
     */
    preemptMode?: string;
    /**
     * Number of minutes beyond MaxTime a job may run before being killed
     * ("INFINITE" or a Slurm time string).
     */
    overTimeLimit?: string;

    // TRES and billing
    /** TRESBillingWeights string (e.g. "CPU=1.0,Mem=0.25G"). */
    tresBillingWeights?: string;
    /** Configured Trackable RESources string (e.g. "cpu=4,mem=8G,node=2"). */
    tres?: string;

    /** QOS name assigned to this partition. */
    qos?: string;

    // Miscellaneous
    /**
     * Seconds a job is given to clean up after its time limit is reached
     * before being forcibly killed (GraceTime in slurm.conf).
     */
    graceTime?: number;
    /** Alternate partition to use when a job cannot run here. */
    alternate?: string;
    /** Slurm hostlist expression for the nodes configured in this partition. */
    nodeList?: string;
    /** Reserved for future per-partition job count limits. */
    maxJobs?: number;
}
