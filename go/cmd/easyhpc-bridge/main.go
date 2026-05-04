// SPDX-License-Identifier: LGPL-2.1-or-later
//
// easyhpc-bridge – Cockpit event bridge for Slurm HPC clusters.
//
// The bridge is spawned by the Cockpit frontend (cockpit.spawn).  It
// communicates over stdin/stdout using newline-delimited JSON:
//
//	Frontend → bridge (stdin):
//	  {"type":"poll","resource":"all"}       poll all resources
//	  {"type":"poll","resource":"nodes"}     poll a single resource
//	  {"type":"refresh"}                     force an immediate refresh
//	  {"type":"ping"}                        keep-alive
//
//	Bridge → frontend (stdout):
//	  {"type":"ready","version":"1.0"}
//	  {"type":"data","resource":"nodes","data":[...],"timestamp":"..."}
//	  {"type":"event","level":"info","message":"...","timestamp":"..."}
//	  {"type":"error","message":"..."}
//	  {"type":"pong"}

package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/lingweicai/easyhpc/bridge/internal/slurm"
)

// inMessage is a command received from the frontend via stdin.
type inMessage struct {
	Type     string `json:"type"`
	Resource string `json:"resource,omitempty"`
}

// outMessage is a message sent to the frontend via stdout.
type outMessage struct {
	Type      string      `json:"type"`
	Resource  string      `json:"resource,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Level     string      `json:"level,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Version   string      `json:"version,omitempty"`
}

var (
	logger  = log.New(os.Stderr, "easyhpc-bridge: ", log.LstdFlags)
	writeMu sync.Mutex
	enc     = json.NewEncoder(os.Stdout)
)

// writeJSON serialises msg to stdout in a thread-safe manner.
func writeJSON(msg outMessage) {
	writeMu.Lock()
	defer writeMu.Unlock()
	if err := enc.Encode(msg); err != nil {
		logger.Printf("encode error: %v", err)
	}
}

// sendResource emits a data frame for a single resource.
func sendResource(resource string, data interface{}) {
	writeJSON(outMessage{
		Type:      "data",
		Resource:  resource,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// sendSnapshot emits data frames for every resource in the cache snapshot.
func sendSnapshot(snapshot map[string]interface{}) {
	for resource, data := range snapshot {
		sendResource(resource, data)
	}
}

func main() {
	// Signal readiness to the frontend.
	writeJSON(outMessage{Type: "ready", Version: "1.0"})

	cache := slurm.NewCache()

	// Perform an initial refresh; errors are non-fatal (Slurm may not be
	// installed in development environments).
	if err := cache.Refresh(); err != nil {
		logger.Printf("initial refresh: %v", err)
		writeJSON(outMessage{
			Type:    "error",
			Message: "Slurm commands unavailable – ensure sinfo, squeue, sacctmgr and scontrol are installed and in PATH: " + err.Error(),
		})
	}
	sendSnapshot(cache.Get())

	// Background goroutine: refresh cached state every 10 seconds.
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := cache.Refresh(); err != nil {
				logger.Printf("refresh: %v", err)
			}
			sendSnapshot(cache.Get())
		}
	}()

	// Background goroutine: tail slurmctld log and push events.
	go func() {
		watcher := slurm.NewLogWatcher(slurm.DefaultLogPaths())
		watcher.Watch(func(event slurm.LogEvent) {
			writeJSON(outMessage{
				Type:      "event",
				Level:     event.Level,
				Message:   event.Message,
				Timestamp: event.Timestamp,
			})
		})
	}()

	// Main loop: read JSON commands from the frontend via stdin.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var msg inMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			logger.Printf("parse error: %v", err)
			continue
		}

		switch msg.Type {
		case "poll":
			if msg.Resource == "" || msg.Resource == "all" {
				sendSnapshot(cache.Get())
			} else {
				sendResource(msg.Resource, cache.GetResource(msg.Resource))
			}
		case "refresh":
			if err := cache.Refresh(); err != nil {
				writeJSON(outMessage{
					Type:    "error",
					Message: "Refresh error: " + err.Error(),
				})
			} else {
				sendSnapshot(cache.Get())
			}
		case "ping":
			writeJSON(outMessage{Type: "pong"})
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Printf("stdin: %v", err)
	}
}
