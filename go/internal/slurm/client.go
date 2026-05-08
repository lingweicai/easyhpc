// SPDX-License-Identifier: LGPL-2.1-or-later

// Package slurm provides utilities for querying Slurm HPC resources and
// watching the slurmctld log for real-time events.
package slurm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CacheManager holds a thread-safe snapshot of all Slurm resources.
type CacheManager struct {
	mu    sync.RWMutex
	cache Cache
}

// NewCache returns an initialised, empty CacheManager.
func NewCache() *CacheManager {
	return &CacheManager{
		cache: Cache{
			Clusters:     []Cluster{},
			Partitions:   []Partition{},
			Nodes:        []Node{},
			Jobs:         []Job{},
			JobSteps:     []JobStep{},
			Reservations: []Reservation{},
			Users:        []User{},
			Accounts:     []Account{},
			SlurmLogs:    []SlurmLog{},
			Events:       []Event{},
		},
	}
}

// Refresh runs all Slurm commands and updates the cache atomically.
// It returns the first error encountered but still stores whatever data
// could be collected.
func (c *CacheManager) Refresh() error {
	clusters, errC := getClusters()
	nodes, errN := getNodes()
	partitions, errP := getPartitions()
	jobs, errJ := getJobs()
	reservations, errR := getReservations()
	users, errU := getUsers()
	accounts, errA := getAccounts()

	c.mu.Lock()
	c.cache.Clusters = clusters
	c.cache.Nodes = nodes
	c.cache.Partitions = partitions
	c.cache.Jobs = jobs
	c.cache.Reservations = reservations
	c.cache.Users = users
	c.cache.Accounts = accounts
	c.cache.LastUpdated = time.Now()
	c.mu.Unlock()

	for _, err := range []error{errC, errN, errP, errJ, errR, errU, errA} {
		if err != nil {
			return err
		}
	}
	return nil
}

// Get returns a copy of all cached resources keyed by resource name.
func (c *CacheManager) Get() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return map[string]interface{}{
		"clusters":     c.cache.Clusters,
		"partitions":   c.cache.Partitions,
		"nodes":        c.cache.Nodes,
		"jobs":         c.cache.Jobs,
		"job_steps":    c.cache.JobSteps,
		"reservations": c.cache.Reservations,
		"users":        c.cache.Users,
		"accounts":     c.cache.Accounts,
		"slurm_logs":   c.cache.SlurmLogs,
		"events":       c.cache.Events,
	}
}

// GetResource returns cached data for a single named resource.
func (c *CacheManager) GetResource(resource string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	switch resource {
	case "clusters":
		return c.cache.Clusters
	case "partitions":
		return c.cache.Partitions
	case "nodes":
		return c.cache.Nodes
	case "jobs":
		return c.cache.Jobs
	case "job_steps":
		return c.cache.JobSteps
	case "reservations":
		return c.cache.Reservations
	case "users":
		return c.cache.Users
	case "accounts":
		return c.cache.Accounts
	case "slurm_logs":
		return c.cache.SlurmLogs
	case "events":
		return c.cache.Events
	}
	return nil
}

// runCommandOutput executes a binary with args and returns the raw stdout
// bytes (suitable for JSON decoding).
func runCommandOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// runCommand executes a binary with args and returns non-empty output lines.
func runCommand(name string, args ...string) ([]string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var lines []string
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// parseInt converts a string to int, returning 0 on error.
func parseInt(s string) int {
	s = strings.TrimSpace(s)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

// parseFloat converts a string to float64, returning 0.0 on error.
func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return v
}

// parseTime converts a Slurm timestamp string (2006-01-02T15:04:05) to
// time.Time, returning zero time on error or for "N/A"/"Unknown" values.
func parseTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" || s == "N/A" || s == "Unknown" || s == "None" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// parseCPUTotal extracts the total CPU count from sinfo's A/I/O/T format.
// e.g. "2/4/0/6" → 6.
func parseCPUTotal(s string) int {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, "/")
	if len(parts) == 4 {
		return parseInt(parts[3])
	}
	return parseInt(s)
}

// getClusters queries sacctmgr for cluster information.
func getClusters() ([]Cluster, error) {
	lines, err := runCommand("sacctmgr", "list", "cluster",
		"--noheader", "--parsable2")
	if err != nil {
		return []Cluster{}, err
	}
	clusters := make([]Cluster, 0, len(lines))
	for _, line := range lines {
		f := strings.Split(line, "|")
		if len(f) < 2 {
			continue
		}
		clusters = append(clusters, Cluster{
			ClusterName: f[0],
			ControlHost: f[1],
		})
	}
	return clusters, nil
}

// getNodes queries scontrol for per-node information using JSON output
// (scontrol show nodes --json).  Falls back gracefully to an empty slice
// on error so the bridge continues operating when Slurm is not available.
func getNodes() ([]Node, error) {
	out, err := runCommandOutput("scontrol", "show", "nodes", "--json")
	if err != nil {
		return []Node{}, err
	}
	var resp SlurmNodesResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return []Node{}, fmt.Errorf("parsing scontrol nodes JSON: %w", err)
	}
	nodes := make([]Node, 0, len(resp.Nodes))
	for _, raw := range resp.Nodes {
		nodes = append(nodes, MapSlurmNodeRaw(raw))
	}
	return nodes, nil
}

// getDefaultPartitionNames queries sinfo for partition names and returns a set
// of those marked as the default partition (names ending with '*' in sinfo
// output).  The '*' suffix is the only reliable way to detect the default
// partition because the Slurm JSON data_parser does not yet export the
// PART_FLAG_DEFAULT bit from the flags field.
func getDefaultPartitionNames() map[string]bool {
	lines, err := runCommand("sinfo", "--noheader", "--format=%P")
	if err != nil {
		return map[string]bool{}
	}
	defaults := make(map[string]bool)
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if strings.HasSuffix(name, "*") {
			defaults[strings.TrimSuffix(name, "*")] = true
		}
	}
	return defaults
}

// getPartitions queries scontrol for partition information using JSON output
// (scontrol show partitions --json, Slurm ≥23.11 data_parser).  Falls back
// gracefully to an empty slice on error so the bridge continues operating
// when Slurm is not available.
//
// The Default flag is determined separately via sinfo because Slurm's JSON
// parser does not yet expose the partition flags bitmask.
func getPartitions() ([]Partition, error) {
	// Identify the default partition(s) before parsing the full JSON payload.
	defaultNames := getDefaultPartitionNames()

	out, err := runCommandOutput("scontrol", "show", "partitions", "--json")
	if err != nil {
		return []Partition{}, err
	}
	var resp SlurmPartitionsResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return []Partition{}, fmt.Errorf("parsing scontrol partitions JSON: %w", err)
	}
	partitions := make([]Partition, 0, len(resp.Partitions))
	for _, raw := range resp.Partitions {
		p := MapSlurmPartitionRaw(raw)
		p.Default = defaultNames[p.PartitionName]
		partitions = append(partitions, p)
	}
	return partitions, nil
}

// getJobs queries scontrol for job information using JSON output
// (scontrol show jobs --json).  Falls back gracefully to an empty slice
// on error so the bridge continues operating when Slurm is not available.
func getJobs() ([]Job, error) {
	out, err := runCommandOutput("scontrol", "show", "jobs", "--json")
	if err != nil {
		return []Job{}, err
	}
	var resp SlurmJobsResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return []Job{}, fmt.Errorf("parsing scontrol jobs JSON: %w", err)
	}
	jobs := make([]Job, 0, len(resp.Jobs))
	for _, raw := range resp.Jobs {
		jobs = append(jobs, MapSlurmJobRaw(raw))
	}
	return jobs, nil
}

// getReservations queries scontrol for advance reservation information.
func getReservations() ([]Reservation, error) {
	lines, err := runCommand("scontrol", "show", "reservation", "--oneliner")
	if err != nil {
		return []Reservation{}, err
	}
	reservations := make([]Reservation, 0)
	for _, line := range lines {
		if !strings.Contains(line, "ReservationName=") {
			continue
		}
		var r Reservation
		r.Nodes = []Node{}
		for _, field := range strings.Fields(line) {
			kv := strings.SplitN(field, "=", 2)
			if len(kv) != 2 {
				continue
			}
			switch kv[0] {
			case "ReservationName":
				r.ReservationID = kv[1]
			case "State":
				r.State = kv[1]
			case "StartTime":
				r.StartTime = parseTime(kv[1])
			case "EndTime":
				r.EndTime = parseTime(kv[1])
			case "TRES":
				// TRES field may contain cpu=N; extract CPU count if present.
				for _, part := range strings.Split(kv[1], ",") {
					if strings.HasPrefix(part, "cpu=") {
						r.CPUs = parseInt(strings.TrimPrefix(part, "cpu="))
					}
				}
			}
		}
		if r.ReservationID != "" {
			reservations = append(reservations, r)
		}
	}
	return reservations, nil
}

// getUsers queries sacctmgr for user information.
func getUsers() ([]User, error) {
	lines, err := runCommand("sacctmgr", "list", "user",
		"--noheader", "--parsable2")
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(lines))
	for _, line := range lines {
		f := strings.Split(line, "|")
		if len(f) < 1 {
			continue
		}
		name := f[0]
		users = append(users, User{
			UserID:   name,
			UserName: name,
		})
	}
	return users, nil
}

// getAccounts queries sacctmgr for account information.
func getAccounts() ([]Account, error) {
	lines, err := runCommand("sacctmgr", "list", "account",
		"--noheader", "--parsable2")
	if err != nil {
		return []Account{}, err
	}
	accounts := make([]Account, 0, len(lines))
	for _, line := range lines {
		f := strings.Split(line, "|")
		if len(f) < 1 {
			continue
		}
		name := f[0]
		accounts = append(accounts, Account{
			AccountID:   name,
			AccountName: name,
			UserNames:   []string{},
		})
	}
	return accounts, nil
}
