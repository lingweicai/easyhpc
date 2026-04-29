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
}
