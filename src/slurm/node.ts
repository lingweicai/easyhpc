/*
 * SPDX-License-Identifier: LGPL-2.1-or-later
 */

/**
 * Slurm node model for the React frontend.
 *
 * These types mirror the normalised Go `Node` model emitted by the
 * easyhpc-bridge over the Cockpit JSON protocol.
 *
 * Backend/frontend field mapping
 * ──────────────────────────────
 * Go field          │ JSON key    │ TS field
 * ──────────────────────────────────────────
 * Node.NodeName     │ node_name   │ nodeName
 * Node.Arch         │ arch        │ arch
 * Node.CPUs         │ cpus        │ cpus
 * Node.Mem          │ mem         │ mem
 * Node.State        │ state       │ state
 * Node.Partitions   │ partitions  │ partitions
 * Node.Sockets      │ sockets     │ sockets
 * Node.FreeMem      │ free_mem    │ freeMem
 * Node.CPULoad      │ cpu_load    │ cpuLoad
 */

// ---------------------------------------------------------------------------
// Node – canonical frontend model
// ---------------------------------------------------------------------------

/**
 * A normalised Slurm compute node as delivered by the easyhpc-bridge.
 * Corresponds to Go Node (sourced from scontrol show nodes --json).
 */
export interface Node {
    nodeName: string;
    /** CPU architecture reported by the node (e.g. "x86_64"). */
    arch?: string;
    /** Total number of logical CPUs on the node. */
    cpus: number;
    /** Total memory in MB on the node. */
    mem: number;
    /** Current node state (e.g. "IDLE", "ALLOCATED", "DOWN", "DRAINED"). */
    state: string;
    /** Names of the partitions this node belongs to. */
    partitions: string[];
    /** Number of physical processor sockets on the node. */
    sockets?: number;
    /** Free memory in MB as reported by the OS; 0 when unavailable. */
    freeMem: number;
    /** CPU load (1-minute average) as reported by the OS. */
    cpuLoad: number;
}
