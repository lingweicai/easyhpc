// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Raw Slurm JSON shapes (scontrol show reservations --json)
//
// These types mirror reservation JSON fields exposed by Slurm's data_parser.
// Not all scontrol create/update specification fields are always emitted by
// current Slurm JSON output; fields may therefore be empty after mapping.
// ---------------------------------------------------------------------------

// SlurmReservationsResponse models the top-level JSON object returned by
// "scontrol show reservations --json".
type SlurmReservationsResponse struct {
	Reservations []SlurmReservationRaw `json:"reservations"`
	Reservation  []SlurmReservationRaw `json:"reservation"`
}

// Items returns the reservation list, supporting both "reservations" and
// legacy/singular "reservation" top-level keys.
func (r SlurmReservationsResponse) Items() []SlurmReservationRaw {
	if len(r.Reservations) > 0 {
		return r.Reservations
	}
	return r.Reservation
}

// SlurmReservationRaw models one reservation entry from Slurm JSON.
// Flexible fields are retained as raw JSON to support version differences
// (string vs optional-number envelope vs scalar).
type SlurmReservationRaw struct {
	Name            string          `json:"name"`
	ReservationName string          `json:"reservation_name"`
	StartTime       json.RawMessage `json:"start_time"`
	EndTime         json.RawMessage `json:"end_time"`
	Duration        string          `json:"duration"`
	Users           string          `json:"users"`
	Accounts        string          `json:"accounts"`
	Groups          string          `json:"groups"`
	NodeList        string          `json:"node_list"`
	Nodes           string          `json:"nodes"`
	NodeCnt         json.RawMessage `json:"node_cnt"`
	NodeCount       json.RawMessage `json:"node_count"`
	CoreCnt         json.RawMessage `json:"core_cnt"`
	CoreCount       json.RawMessage `json:"core_count"`
	PartitionName   string          `json:"partition_name"`
	Partition       string          `json:"partition"`
	Features        string          `json:"features"`
	Licenses        string          `json:"licenses"`
	TRES            string          `json:"tres"`
	BurstBuffer     string          `json:"burst_buffer"`
	MaxStartDelay   json.RawMessage `json:"max_start_delay"`
	Flags           json.RawMessage `json:"flags"`
	State           json.RawMessage `json:"state"`
}

// MapSlurmReservationRaw converts a SlurmReservationRaw decoded from
// "scontrol show reservations --json" into the canonical Reservation model.
func MapSlurmReservationRaw(raw SlurmReservationRaw) Reservation {
	name := firstNonEmpty(raw.Name, raw.ReservationName)
	nodeList := firstNonEmpty(raw.NodeList, raw.Nodes)
	partition := firstNonEmpty(raw.PartitionName, raw.Partition)
	nodeCount := firstNonZero(parseSlurmIntRaw(raw.NodeCnt), parseSlurmIntRaw(raw.NodeCount))
	coreCount := firstNonZero(parseSlurmIntRaw(raw.CoreCnt), parseSlurmIntRaw(raw.CoreCount))

	return Reservation{
		ReservationID:   name,
		ReservationName: name,

		StartTime: parseSlurmTimeRaw(raw.StartTime),
		EndTime:   parseSlurmTimeRaw(raw.EndTime),
		Duration:  raw.Duration,

		Users:    raw.Users,
		Accounts: raw.Accounts,
		Groups:   raw.Groups,

		Nodes:         []Node{},
		NodeList:      nodeList,
		NodeCount:     nodeCount,
		CoreCount:     coreCount,
		CPUs:          extractCPUFromTRES(raw.TRES),
		PartitionName: partition,
		Features:      raw.Features,
		Licenses:      raw.Licenses,
		TRES:          raw.TRES,
		BurstBuffer:   raw.BurstBuffer,

		Flags:         parseSlurmStringListRaw(raw.Flags),
		MaxStartDelay: parseSlurmStringRaw(raw.MaxStartDelay),

		State: firstState(parseSlurmStringListRaw(raw.State)),
	}
}

func parseSlurmTimeRaw(raw json.RawMessage) *time.Time {
	raw = compactRawJSON(raw)
	if len(raw) == 0 {
		return nil
	}

	// Optional-number envelope ({"set":true,"infinite":false,"number":...}).
	if raw[0] == '{' {
		var opt SlurmOptionalInt
		if err := json.Unmarshal(raw, &opt); err == nil {
			if !opt.Set || opt.Infinite || opt.Number == 0 {
				return nil
			}
			t := time.Unix(int64(opt.Number), 0).UTC()
			return &t
		}
	}

	// Unix epoch numeric value.
	var unix float64
	if err := json.Unmarshal(raw, &unix); err == nil && unix > 0 {
		t := time.Unix(int64(unix), 0).UTC()
		return &t
	}

	// Textual timestamps.
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil
	}
	s = strings.TrimSpace(s)
	if s == "" || s == "N/A" || s == "Unknown" || s == "None" {
		return nil
	}

	// Numeric string epoch.
	if n, err := strconv.ParseInt(s, 10, 64); err == nil && n > 0 {
		t := time.Unix(n, 0).UTC()
		return &t
	}

	// Slurm oneliner format.
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		u := t.UTC()
		return &u
	}
	// RFC3339-compatible textual format.
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		u := t.UTC()
		return &u
	}
	return nil
}

func parseSlurmIntRaw(raw json.RawMessage) int {
	raw = compactRawJSON(raw)
	if len(raw) == 0 {
		return 0
	}

	// Optional-number envelope.
	if raw[0] == '{' {
		var opt SlurmOptionalInt
		if err := json.Unmarshal(raw, &opt); err == nil {
			if !opt.Set || opt.Infinite {
				return 0
			}
			return int(opt.Number)
		}
	}

	var n float64
	if err := json.Unmarshal(raw, &n); err == nil {
		return int(n)
	}

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return parseInt(s)
	}

	return 0
}

func parseSlurmStringRaw(raw json.RawMessage) string {
	raw = compactRawJSON(raw)
	if len(raw) == 0 {
		return ""
	}

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return strings.TrimSpace(s)
	}

	// Optional-number envelope.
	if raw[0] == '{' {
		var opt SlurmOptionalInt
		if err := json.Unmarshal(raw, &opt); err == nil {
			if !opt.Set || opt.Infinite {
				return ""
			}
			return strconv.Itoa(int(opt.Number))
		}
	}

	var n float64
	if err := json.Unmarshal(raw, &n); err == nil && n != 0 {
		return strconv.Itoa(int(n))
	}

	return ""
}

func parseSlurmStringListRaw(raw json.RawMessage) []string {
	raw = compactRawJSON(raw)
	if len(raw) == 0 {
		return nil
	}

	var list []string
	if err := json.Unmarshal(raw, &list); err == nil {
		return list
	}

	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		single = strings.TrimSpace(single)
		if single == "" {
			return nil
		}
		if strings.Contains(single, ",") {
			parts := strings.Split(single, ",")
			out := make([]string, 0, len(parts))
			for _, p := range parts {
				if v := strings.TrimSpace(p); v != "" {
					out = append(out, v)
				}
			}
			if len(out) == 0 {
				return nil
			}
			return out
		}
		return []string{single}
	}

	return nil
}

func compactRawJSON(raw json.RawMessage) json.RawMessage {
	return json.RawMessage(strings.TrimSpace(string(raw)))
}

func firstState(states []string) string {
	if len(states) == 0 {
		return ""
	}
	return states[0]
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func firstNonZero(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

func extractCPUFromTRES(tres string) int {
	for _, part := range strings.Split(tres, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "cpu=") {
			return parseInt(strings.TrimPrefix(part, "cpu="))
		}
	}
	return 0
}
