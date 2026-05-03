// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"
)

// LogWatcher tails slurmctld log files and emits LogEvents via a callback.
type LogWatcher struct {
	paths []string
}

// NewLogWatcher creates a LogWatcher that will search the given paths in order.
func NewLogWatcher(paths []string) *LogWatcher {
	return &LogWatcher{paths: paths}
}

// DefaultLogPaths returns the common locations for the slurmctld log file.
func DefaultLogPaths() []string {
	return []string{
		"/var/log/slurm/slurmctld.log",
		"/var/log/slurmctld.log",
		"/var/log/slurm-llnl/slurmctld.log",
	}
}

// Watch tails the first log file found and invokes callback for every new
// line.  It returns immediately if no log file exists.
func (w *LogWatcher) Watch(callback func(LogEvent)) {
	var logPath string
	for _, p := range w.paths {
		if _, err := os.Stat(p); err == nil {
			logPath = p
			break
		}
	}
	if logPath == "" {
		return
	}

	f, err := os.Open(logPath)
	if err != nil {
		return
	}
	defer f.Close()

	// Start reading from end of file so we only emit new events.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return
	}

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return
			}
			// No new data yet; poll again shortly.
			time.Sleep(500 * time.Millisecond)
			continue
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}
		callback(parseLogLine(line))
	}
}

// parseLogLine converts a raw slurmctld log line into a LogEvent.
// slurmctld lines look like: [2024-01-01T00:00:00.000] error: ...
func parseLogLine(line string) LogEvent {
	level := "info"
	lower := strings.ToLower(line)
	switch {
	case strings.Contains(lower, "error"):
		level = "error"
	case strings.Contains(lower, "warn"):
		level = "warning"
	case strings.Contains(lower, "debug"):
		level = "debug"
	}

	// Extract timestamp from the log line if present ([YYYY-MM-DDTHH:MM:SS.mmm]).
	ts := time.Now().Format(time.RFC3339)
	if strings.HasPrefix(line, "[") {
		end := strings.Index(line, "]")
		if end > 1 {
			raw := line[1:end]
			// Trim sub-second precision before parsing.
			if dot := strings.Index(raw, "."); dot > 0 {
				raw = raw[:dot]
			}
			if t, err := time.Parse("2006-01-02T15:04:05", raw); err == nil {
				ts = t.Format(time.RFC3339)
			}
		}
	}

	return LogEvent{
		Level:     level,
		Message:   line,
		Timestamp: ts,
	}
}
