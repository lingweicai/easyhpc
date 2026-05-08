// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

// ---------------------------------------------------------------------------
// Raw Slurm JSON shapes (scontrol show nodes --json)
//
// These types mirror the JSON produced by Slurm's data_parser plugin (≥23.11).
// They are used only for decoding; callers should convert to the normalised
// Node model via MapSlurmNodeRaw.
// ---------------------------------------------------------------------------

// SlurmNodeRaw models the JSON shape of a single node entry from
// "scontrol show nodes --json".  Only the fields needed for the canonical
// Node model are decoded here; the struct can be extended as required.
type SlurmNodeRaw struct {
	Name         string           `json:"name"`
	Architecture string           `json:"architecture"`
	CPUs         int              `json:"cpus"`
	Sockets      int              `json:"sockets"`
	RealMemory   int              `json:"real_memory"`
	CPULoad      float64          `json:"cpu_load"`
	State        []string         `json:"state"`
	Partitions   []string         `json:"partitions"`
	FreeMem      SlurmOptionalInt `json:"free_mem"`
}

// SlurmNodesResponse models the top-level JSON object returned by
// "scontrol show nodes --json".
type SlurmNodesResponse struct {
	Nodes []SlurmNodeRaw `json:"nodes"`
}

// ---------------------------------------------------------------------------
// Mapping
// ---------------------------------------------------------------------------

// MapSlurmNodeRaw converts a SlurmNodeRaw (decoded from scontrol show nodes
// --json) into the normalised Node model used by the API and the frontend.
//
// FreeMem is derived from the optional-number envelope: if Set && !Infinite
// the integer value is used; otherwise it defaults to 0.
func MapSlurmNodeRaw(raw SlurmNodeRaw) Node {
	// Use the first element of state as the canonical state string.
	state := ""
	if len(raw.State) > 0 {
		state = raw.State[0]
	}

	// Derive free memory from the optional envelope.
	freeMem := 0
	if raw.FreeMem.Set && !raw.FreeMem.Infinite {
		freeMem = int(raw.FreeMem.Number)
	}

	return Node{
		NodeName:   raw.Name,
		Arch:       raw.Architecture,
		CPUs:       raw.CPUs,
		Mem:        raw.RealMemory,
		State:      state,
		Partitions: raw.Partitions,
		Sockets:    raw.Sockets,
		FreeMem:    freeMem,
		CPULoad:    raw.CPULoad,
	}
}
