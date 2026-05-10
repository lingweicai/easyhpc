// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm_test

import (
	"testing"

	"github.com/lingweicai/easyhpc/bridge/internal/slurm"
)

func TestNewCacheIsEmpty(t *testing.T) {
	c := slurm.NewCache()
	snap := c.Get()

	resources := []string{"clusters", "partitions", "nodes", "jobs", "reservations", "users", "accounts"}
	resources := []string{"clusters", "partitions", "nodes", "jobs", "reservations", "users", "accounts", "slurmdb"}
	for _, r := range resources {
		v, ok := snap[r]
		if !ok {
			t.Errorf("missing resource %q in initial snapshot", r)
			continue
		}
		// Each value must be a non-nil slice (may be empty).
		if v == nil {
			t.Errorf("resource %q is nil, want empty slice", r)
		}
	}
}

func TestNewCacheHasExtendedResources(t *testing.T) {
	c := slurm.NewCache()
	snap := c.Get()

	extended := []string{"job_steps", "slurm_logs", "events"}
	for _, r := range extended {
		v, ok := snap[r]
		if !ok {
			t.Errorf("missing extended resource %q in initial snapshot", r)
			continue
		}
		if v == nil {
			t.Errorf("extended resource %q is nil, want empty slice", r)
		}
	}
}

func TestGetResourceReturnsCorrectType(t *testing.T) {
	c := slurm.NewCache()

	cases := []struct {
		resource string
		wantNil  bool
	}{
		{"clusters", false},
		{"nodes", false},
		{"partitions", false},
		{"jobs", false},
		{"reservations", false},
		{"users", false},
		{"accounts", false},
		{"job_steps", false},
		{"reservations", false},
		{"users", false},
		{"accounts", false},
		{"slurmdb", false},
		{"slurmdb_clusters", false},
		{"slurmdb_accounts", false},
		{"slurmdb_users", false},
		{"slurmdb_associations", false},
		{"slurmdb_qos", false},
		{"slurmdb_wckeys", false},
		{"slurmdb_tres", false},
		{"slurm_logs", false},
		{"events", false},
		{"unknown", true},
	}

	for _, tc := range cases {
		v := c.GetResource(tc.resource)
		isNil := v == nil
		if isNil != tc.wantNil {
			t.Errorf("GetResource(%q): got nil=%v, want nil=%v", tc.resource, isNil, tc.wantNil)
		}
	}
}

// TestRefreshWithoutSlurm verifies that Refresh returns an error (because
// Slurm binaries are not installed in CI) but does NOT panic and still
// leaves the cache in a valid state.
func TestRefreshWithoutSlurm(t *testing.T) {
	c := slurm.NewCache()
	// We expect an error here because sinfo/sacctmgr are not available.
	_ = c.Refresh()

	snap := c.Get()
	if len(snap) == 0 {
		t.Error("Get() returned empty map after Refresh(), want 7 resource entries")
		t.Error("Get() returned empty map after Refresh(), want resource entries")
	}
}

func TestDefaultLogPaths(t *testing.T) {
	paths := slurm.DefaultLogPaths()
	if len(paths) == 0 {
		t.Error("DefaultLogPaths() returned empty slice")
	}
	for _, p := range paths {
		if p == "" {
			t.Error("DefaultLogPaths() contains an empty string")
		}
	}
}

func TestLogWatcherNoFile(t *testing.T) {
	// When no log file exists the watcher should return immediately
	// without calling the callback.
	w := slurm.NewLogWatcher([]string{"/nonexistent/slurmctld.log"})
	called := false
	w.Watch(func(slurm.LogEvent) { called = true })
	if called {
		t.Error("Watch callback was called despite no log file existing")
	}
}
