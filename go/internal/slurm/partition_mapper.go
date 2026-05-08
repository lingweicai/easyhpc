// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

import "fmt"

// ---------------------------------------------------------------------------
// Raw Slurm JSON shapes (scontrol show partitions --json)
//
// The Slurm data_parser plugin (≥23.11) serialises partition_info_t fields
// as a nested JSON object whose structure matches the keys used in
// src/plugins/data_parser/v0.0.42/parsers.c.  Paths that contain "/" become
// nested sub-objects in the JSON output (e.g. "nodes/total" → nodes.total).
//
// These types are used only for decoding; callers should convert to the
// normalised Partition model via MapSlurmPartitionRaw.
// ---------------------------------------------------------------------------

// SlurmPartitionNodes corresponds to the "nodes" sub-object in the partition
// JSON, which aggregates several partition_info_t fields:
//
//	nodes/allowed_allocation → AllowedAllocation (AllowNodes/AllocNodes)
//	nodes/configured         → Configured        (Nodes hostlist)
//	nodes/total              → Total             (TotalNodes)
type SlurmPartitionNodes struct {
	AllowedAllocation string `json:"allowed_allocation"`
	Configured        string `json:"configured"`
	Total             int    `json:"total"`
}

// SlurmPartitionAccounts holds the accounts sub-object (AllowAccounts /
// DenyAccounts).
type SlurmPartitionAccounts struct {
	Allowed string `json:"allowed"`
	Deny    string `json:"deny"`
}

// SlurmPartitionGroups holds the groups sub-object (AllowGroups).
type SlurmPartitionGroups struct {
	Allowed string `json:"allowed"`
}

// SlurmPartitionQOS holds the qos sub-object (AllowQOS / DenyQOS / QOS).
type SlurmPartitionQOS struct {
	Allowed  string `json:"allowed"`
	Deny     string `json:"deny"`
	Assigned string `json:"assigned"`
}

// SlurmPartitionDefaults holds the defaults sub-object.
type SlurmPartitionDefaults struct {
	// Time is DefaultTime in minutes (UINT32_NO_VAL envelope).
	Time SlurmOptionalInt `json:"time"`
	// Job is the JobDefaults string.
	Job string `json:"job"`
}

// SlurmPartitionOversubscribe holds the oversubscribe sub-object.
type SlurmPartitionOversubscribe struct {
	// Jobs is the maximum number of jobs allowed to share a resource.
	Jobs int `json:"jobs"`
	// Flags contains "force" when OverSubscribe is FORCE.
	Flags []string `json:"flags"`
}

// SlurmPartitionMaximums holds the maximums sub-object for a partition.
type SlurmPartitionMaximums struct {
	Nodes         SlurmOptionalInt            `json:"nodes"`
	Time          SlurmOptionalInt            `json:"time"`
	CPUsPerNode   SlurmOptionalInt            `json:"cpus_per_node"`
	CPUsPerSocket SlurmOptionalInt            `json:"cpus_per_socket"`
	OverTimeLimit SlurmOptionalInt            `json:"over_time_limit"`
	Oversubscribe SlurmPartitionOversubscribe `json:"oversubscribe"`
}

// SlurmPartitionMinimums holds the minimums sub-object.
type SlurmPartitionMinimums struct {
	Nodes int `json:"nodes"`
}

// SlurmPartitionPriority holds the priority sub-object.
type SlurmPartitionPriority struct {
	JobFactor int `json:"job_factor"`
	Tier      int `json:"tier"`
}

// SlurmPartitionCPUs holds the cpus sub-object.
type SlurmPartitionCPUs struct {
	Total int `json:"total"`
}

// SlurmPartitionState holds the partition sub-object which carries the
// partition state flags (e.g. ["UP"], ["DOWN"], ["DRAIN"]).
type SlurmPartitionState struct {
	State []string `json:"state"`
}

// SlurmPartitionTRES holds the tres sub-object.
type SlurmPartitionTRES struct {
	// BillingWeights is the TRESBillingWeights string.
	BillingWeights string `json:"billing_weights"`
	// Configured is the TRES string (cpu=N,mem=M,...).
	Configured string `json:"configured"`
}

// SlurmPartitionRaw models the JSON shape of a single partition entry from
// "scontrol show partitions --json" (Slurm data_parser ≥23.11).
//
// Only the fields required for the canonical Partition model are decoded; the
// struct can be extended as the REST API coverage grows.
type SlurmPartitionRaw struct {
	Name      string                  `json:"name"`
	Alternate string                  `json:"alternate"`
	NodeSets  string                  `json:"node_sets"`
	Nodes     SlurmPartitionNodes     `json:"nodes"`
	Accounts  SlurmPartitionAccounts  `json:"accounts"`
	Groups    SlurmPartitionGroups    `json:"groups"`
	QOS       SlurmPartitionQOS       `json:"qos"`
	Defaults  SlurmPartitionDefaults  `json:"defaults"`
	Maximums  SlurmPartitionMaximums  `json:"maximums"`
	Minimums  SlurmPartitionMinimums  `json:"minimums"`
	Priority  SlurmPartitionPriority  `json:"priority"`
	CPUs      SlurmPartitionCPUs      `json:"cpus"`
	Partition SlurmPartitionState     `json:"partition"`
	TRES      SlurmPartitionTRES      `json:"tres"`
	GraceTime int                     `json:"grace_time"`
}

// SlurmPartitionsResponse models the top-level JSON object returned by
// "scontrol show partitions --json".
type SlurmPartitionsResponse struct {
	Partitions []SlurmPartitionRaw `json:"partitions"`
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// formatPartitionTimeLimit converts a Slurm UINT32_NO_VAL optional-integer
// (in minutes) to a human-readable time string using Slurm's D-HH:MM:SS
// notation.
//
//   - infinite=true             → "INFINITE"
//   - set=false, infinite=false → "NONE" (unset / use global default)
//   - set=true                  → formatted minute value
func formatPartitionTimeLimit(opt SlurmOptionalInt) string {
	if opt.Infinite {
		return "INFINITE"
	}
	if !opt.Set {
		return "NONE"
	}
	minutes := int(opt.Number)
	days := minutes / (60 * 24)
	rem := minutes % (60 * 24)
	hours := rem / 60
	mins := rem % 60
	if days > 0 {
		return fmt.Sprintf("%d-%02d:%02d:00", days, hours, mins)
	}
	return fmt.Sprintf("%02d:%02d:00", hours, mins)
}

// formatOverSubscribe converts the oversubscribe count and flags into the
// Slurm OverSubscribe string:
//
//   - jobs ≤ 1, no force → "NO"
//   - jobs > 1, force     → "FORCE:N"
//   - jobs > 1, no force  → "YES:N"
func formatOverSubscribe(os SlurmPartitionOversubscribe) string {
	force := false
	for _, f := range os.Flags {
		if f == "force" {
			force = true
			break
		}
	}
	if os.Jobs <= 1 && !force {
		return "NO"
	}
	if force {
		return fmt.Sprintf("FORCE:%d", os.Jobs)
	}
	return fmt.Sprintf("YES:%d", os.Jobs)
}

// optionalIntValue returns the integer value of a SlurmOptionalInt when set
// and not infinite; otherwise it returns 0.
func optionalIntValue(opt SlurmOptionalInt) int {
	if opt.Set && !opt.Infinite {
		return int(opt.Number)
	}
	return 0
}

// ---------------------------------------------------------------------------
// Mapping
// ---------------------------------------------------------------------------

// MapSlurmPartitionRaw converts a SlurmPartitionRaw (decoded from
// "scontrol show partitions --json") into the canonical Partition model.
//
// The Default field is not present in the Slurm JSON output (the flags field
// is not yet exported by the data_parser); callers should populate it
// separately from "sinfo --noheader --format=%P" output.
func MapSlurmPartitionRaw(raw SlurmPartitionRaw) Partition {
	state := ""
	if len(raw.Partition.State) > 0 {
		state = raw.Partition.State[0]
	}

	overTimeLimit := formatPartitionTimeLimit(raw.Maximums.OverTimeLimit)
	if overTimeLimit == "NONE" {
		overTimeLimit = ""
	}

	return Partition{
		PartitionName:      raw.Name,
		State:              state,
		MaxTime:            formatPartitionTimeLimit(raw.Maximums.Time),
		DefaultTime:        formatPartitionTimeLimit(raw.Defaults.Time),
		TotalNodes:         raw.Nodes.Total,
		TotalCPUs:          raw.CPUs.Total,
		MinNodes:           raw.Minimums.Nodes,
		MaxNodes:           optionalIntValue(raw.Maximums.Nodes),
		AllowGroups:        raw.Groups.Allowed,
		AllowAccounts:      raw.Accounts.Allowed,
		AllowQOS:           raw.QOS.Allowed,
		DenyAccounts:       raw.Accounts.Deny,
		DenyQOS:            raw.QOS.Deny,
		AllocNodes:         raw.Nodes.AllowedAllocation,
		MaxCPUsPerNode:     optionalIntValue(raw.Maximums.CPUsPerNode),
		MaxCPUsPerSocket:   optionalIntValue(raw.Maximums.CPUsPerSocket),
		PriorityJobFactor:  raw.Priority.JobFactor,
		PriorityTier:       raw.Priority.Tier,
		OverSubscribe:      formatOverSubscribe(raw.Maximums.Oversubscribe),
		OverTimeLimit:      overTimeLimit,
		TRESBillingWeights: raw.TRES.BillingWeights,
		TRES:               raw.TRES.Configured,
		QOS:                raw.QOS.Assigned,
		GraceTime:          raw.GraceTime,
		Alternate:          raw.Alternate,
		NodeList:           raw.Nodes.Configured,
		Nodes:              []Node{},
	}
}
