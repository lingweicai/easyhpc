// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm_test

import (
	"testing"

	"github.com/lingweicai/easyhpc/bridge/internal/slurm"
)

func int64Value(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func float64Value(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func TestMapSlurmDBClusterRecord(t *testing.T) {
	record := slurm.MapSlurmDBClusterRecord([]string{"alpha", "ctl01", "production", "science", "fed-a"})

	if record.Name != "alpha" {
		t.Errorf("Name: got %q, want alpha", record.Name)
	}
	if record.ControlHost != "ctl01" {
		t.Errorf("ControlHost: got %q, want ctl01", record.ControlHost)
	}
	if record.Classification != "production" {
		t.Errorf("Classification: got %q, want production", record.Classification)
	}
	if record.Organization != "science" {
		t.Errorf("Organization: got %q, want science", record.Organization)
	}
	if record.Federation != "fed-a" {
		t.Errorf("Federation: got %q, want fed-a", record.Federation)
	}
}

func TestMapSlurmDBAccountRecord(t *testing.T) {
	record := slurm.MapSlurmDBAccountRecord([]string{"research", "PI group", "science", "root", "normal", "normal,burst", "denynew,usagefactor"})

	if record.Name != "research" {
		t.Errorf("Name: got %q, want research", record.Name)
	}
	if record.ParentAccount != "root" {
		t.Errorf("ParentAccount: got %q, want root", record.ParentAccount)
	}
	if record.DefaultQOS != "normal" {
		t.Errorf("DefaultQOS: got %q, want normal", record.DefaultQOS)
	}
	if len(record.QOSList) != 2 || record.QOSList[0] != "normal" || record.QOSList[1] != "burst" {
		t.Errorf("QOSList: got %v, want [normal burst]", record.QOSList)
	}
	if len(record.Flags) != 2 || record.Flags[0] != "denynew" || record.Flags[1] != "usagefactor" {
		t.Errorf("Flags: got %v, want [denynew usagefactor]", record.Flags)
	}
}

func TestMapSlurmDBUserRecord(t *testing.T) {
	record := slurm.MapSlurmDBUserRecord([]string{"alice", "research", "project-x", "Admin", "research,teaching"})

	if record.Name != "alice" {
		t.Errorf("Name: got %q, want alice", record.Name)
	}
	if record.DefaultAccount != "research" {
		t.Errorf("DefaultAccount: got %q, want research", record.DefaultAccount)
	}
	if record.DefaultWckey != "project-x" {
		t.Errorf("DefaultWckey: got %q, want project-x", record.DefaultWckey)
	}
	if record.AdminLevel != "Admin" {
		t.Errorf("AdminLevel: got %q, want Admin", record.AdminLevel)
	}
	if len(record.CoordinatorAccounts) != 2 || record.CoordinatorAccounts[0] != "research" || record.CoordinatorAccounts[1] != "teaching" {
		t.Errorf("CoordinatorAccounts: got %v, want [research teaching]", record.CoordinatorAccounts)
	}
}

func TestMapSlurmDBAssociationRecord(t *testing.T) {
	record := slurm.MapSlurmDBAssociationRecord([]string{
		"42", "alpha", "research", "alice", "normal", "normal",
		"7", "100", "120", "2-00:00:00", "80", "90",
		"cpu=1000,mem=2T", "cpu=64", "gpu=4", "cpu=1024", "10", "5", "normal,burst",
	})

	if int64Value(record.ID) != 42 {
		t.Errorf("ID: got %d, want 42", int64Value(record.ID))
	}
	if record.Cluster != "alpha" || record.Account != "research" || record.User != "alice" {
		t.Errorf("identity: got cluster=%q account=%q user=%q", record.Cluster, record.Account, record.User)
	}
	if !record.IsDefault {
		t.Error("IsDefault: got false, want true")
	}
	if int64Value(record.ParentID) != 7 {
		t.Errorf("ParentID: got %d, want 7", int64Value(record.ParentID))
	}
	if int64Value(record.MaxJobs) != 100 {
		t.Errorf("MaxJobs: got %d, want 100", int64Value(record.MaxJobs))
	}
	if record.MaxWall != "2-00:00:00" {
		t.Errorf("MaxWall: got %q, want 2-00:00:00", record.MaxWall)
	}
	if record.GrpTRES != "cpu=1000,mem=2T" {
		t.Errorf("GrpTRES: got %q, want cpu=1000,mem=2T", record.GrpTRES)
	}
	if len(record.QOSList) != 2 || record.QOSList[0] != "normal" || record.QOSList[1] != "burst" {
		t.Errorf("QOSList: got %v, want [normal burst]", record.QOSList)
	}
}

func TestMapSlurmDBAssociationRecord_DefaultQOSPrependedWhenMissingFromList(t *testing.T) {
	record := slurm.MapSlurmDBAssociationRecord([]string{
		"42", "alpha", "research", "alice", "normal", "normal",
		"", "", "", "", "", "",
		"", "", "", "", "", "", "burst",
	})

	if !record.IsDefault {
		t.Error("IsDefault: got false, want true")
	}
	if len(record.QOSList) != 2 || record.QOSList[0] != "normal" || record.QOSList[1] != "burst" {
		t.Errorf("QOSList: got %v, want [normal burst]", record.QOSList)
	}
}

func TestMapSlurmDBQOSRecord(t *testing.T) {
	record := slurm.MapSlurmDBQOSRecord([]string{
		"burst", "short queue", "1000", "preempt_exempt,relative", "10", "20", "30", "40",
		"cpu=400", "cpu=8", "cpu=128", "01:00:00",
	})

	if record.Name != "burst" {
		t.Errorf("Name: got %q, want burst", record.Name)
	}
	if int64Value(record.Priority) != 1000 {
		t.Errorf("Priority: got %d, want 1000", int64Value(record.Priority))
	}
	if len(record.Flags) != 2 || record.Flags[0] != "preempt_exempt" || record.Flags[1] != "relative" {
		t.Errorf("Flags: got %v", record.Flags)
	}
	if int64Value(record.MaxJobsPerUser) != 10 {
		t.Errorf("MaxJobsPerUser: got %d, want 10", int64Value(record.MaxJobsPerUser))
	}
	if record.MaxWall != "01:00:00" {
		t.Errorf("MaxWall: got %q, want 01:00:00", record.MaxWall)
	}
}

func TestMapSlurmDBWckeyRecord(t *testing.T) {
	record := slurm.MapSlurmDBWckeyRecord([]string{"project-x", "alpha", "alice", "Y"})

	if record.Name != "project-x" {
		t.Errorf("Name: got %q, want project-x", record.Name)
	}
	if !record.IsDefault {
		t.Error("IsDefault: got false, want true")
	}
}

func TestMapSlurmDBTRESRecord(t *testing.T) {
	record := slurm.MapSlurmDBTRESRecord([]string{"1", "cpu", "cpu", "1.25"})

	if int64Value(record.ID) != 1 {
		t.Errorf("ID: got %d, want 1", int64Value(record.ID))
	}
	if record.Type != "cpu" || record.Name != "cpu" {
		t.Errorf("identity: got type=%q name=%q", record.Type, record.Name)
	}
	if float64Value(record.BillingWeight) != 1.25 {
		t.Errorf("BillingWeight: got %v, want 1.25", float64Value(record.BillingWeight))
	}
}

func TestSlurmDBResourceWrappers(t *testing.T) {
	cache := slurm.NewCache()
	resource, ok := cache.GetResource("slurmdb_accounts").(slurm.SlurmDBRecordsResource)
	if !ok {
		t.Fatal("slurmdb_accounts: unexpected resource type")
	}

	if resource.SchemaVersion == "" {
		t.Error("SchemaVersion: expected non-empty value")
	}
	if resource.Meta.Source != "sacctmgr" {
		t.Errorf("Meta.Source: got %q, want sacctmgr", resource.Meta.Source)
	}
	records, ok := resource.Records.([]slurm.SlurmDBAccount)
	if !ok {
		t.Fatal("Records: unexpected concrete type")
	}
	if records == nil {
		t.Error("Records: got nil, want empty slice")
	}
}
