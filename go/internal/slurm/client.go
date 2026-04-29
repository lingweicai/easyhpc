// SPDX-License-Identifier: LGPL-2.1-or-later

// Package slurm provides utilities for querying Slurm HPC resources and
// watching the slurmctld log for real-time events.
package slurm

import (
	"bufio"
	"bytes"
	"os/exec"
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
			Reservations: []Reservation{},
			Users:        []User{},
			Accounts:     []Account{},
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
		"reservations": c.cache.Reservations,
		"users":        c.cache.Users,
		"accounts":     c.cache.Accounts,
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
	case "reservations":
		return c.cache.Reservations
	case "users":
		return c.cache.Users
	case "accounts":
		return c.cache.Accounts
	}
	return nil
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
		if len(f) < 3 {
			continue
		}
		clusters = append(clusters, Cluster{
			Name:        f[0],
			ControlHost: f[1],
			ControlPort: f[2],
		})
	}
	return clusters, nil
}

// getNodes queries sinfo for per-node information.
func getNodes() ([]Node, error) {
	lines, err := runCommand("sinfo",
		"--Node", "--noheader",
		"--format=%N|%P|%t|%c|%m|%r")
	if err != nil {
		return []Node{}, err
	}
	nodeMap := make(map[string]*Node)
	for _, line := range lines {
		f := strings.Split(line, "|")
		if len(f) < 6 {
			continue
		}
		name := strings.TrimSpace(f[0])
		partition := strings.TrimSuffix(strings.TrimSpace(f[1]), "*")
		if existing, ok := nodeMap[name]; ok {
			existing.Partition = existing.Partition + "," + partition
		} else {
			nodeMap[name] = &Node{
				Name:      name,
				Partition: partition,
				State:     strings.TrimSpace(f[2]),
				CPUs:      strings.TrimSpace(f[3]),
				Memory:    strings.TrimSpace(f[4]),
				Reason:    strings.TrimSpace(f[5]),
			}
		}
	}
	nodes := make([]Node, 0, len(nodeMap))
	for _, n := range nodeMap {
		nodes = append(nodes, *n)
	}
	return nodes, nil
}

// getPartitions queries sinfo for partition-level information.
func getPartitions() ([]Partition, error) {
	lines, err := runCommand("sinfo",
		"--noheader",
		"--format=%P|%a|%F|%D")
	if err != nil {
		return []Partition{}, err
	}
	seen := make(map[string]bool)
	partitions := make([]Partition, 0)
	for _, line := range lines {
		f := strings.Split(line, "|")
		if len(f) < 4 {
			continue
		}
		name := strings.TrimSuffix(strings.TrimSpace(f[0]), "*")
		if seen[name] {
			continue
		}
		seen[name] = true
		partitions = append(partitions, Partition{
			Name:       name,
			State:      strings.TrimSpace(f[1]),
			TotalCPUs:  strings.TrimSpace(f[2]),
			TotalNodes: strings.TrimSpace(f[3]),
		})
	}
	return partitions, nil
}

// getJobs queries squeue for job information.
func getJobs() ([]Job, error) {
	lines, err := runCommand("squeue",
		"--noheader",
		"--format=%i|%j|%u|%a|%T|%P|%D|%M|%l|%r")
	if err != nil {
		return []Job{}, err
	}
	jobs := make([]Job, 0, len(lines))
	for _, line := range lines {
		f := strings.Split(line, "|")
		if len(f) < 10 {
			continue
		}
		jobs = append(jobs, Job{
			JobID:     strings.TrimSpace(f[0]),
			Name:      strings.TrimSpace(f[1]),
			User:      strings.TrimSpace(f[2]),
			Account:   strings.TrimSpace(f[3]),
			State:     strings.TrimSpace(f[4]),
			Partition: strings.TrimSpace(f[5]),
			Nodes:     strings.TrimSpace(f[6]),
			Time:      strings.TrimSpace(f[7]),
			TimeLimit: strings.TrimSpace(f[8]),
			Reason:    strings.TrimSpace(f[9]),
		})
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
		for _, field := range strings.Fields(line) {
			kv := strings.SplitN(field, "=", 2)
			if len(kv) != 2 {
				continue
			}
			switch kv[0] {
			case "ReservationName":
				r.Name = kv[1]
			case "State":
				r.State = kv[1]
			case "Nodes":
				r.Nodes = kv[1]
			case "StartTime":
				r.StartTime = kv[1]
			case "EndTime":
				r.EndTime = kv[1]
			case "Duration":
				r.Duration = kv[1]
			case "Users":
				r.Users = kv[1]
			case "Accounts":
				r.Accounts = kv[1]
			}
		}
		if r.Name != "" {
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
		if len(f) < 3 {
			continue
		}
		users = append(users, User{
			Name:           f[0],
			DefaultAccount: f[1],
			Admin:          f[2],
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
		if len(f) < 3 {
			continue
		}
		accounts = append(accounts, Account{
			Name:         f[0],
			Description:  f[1],
			Organization: f[2],
		})
	}
	return accounts, nil
}
