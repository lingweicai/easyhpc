// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func slurmDBCommandProfiles() map[string]SlurmDBCommandProfile {
	return map[string]SlurmDBCommandProfile{
		"clusters": {
			Object:    "cluster",
			Fields:    []string{"Cluster", "ControlHost", "Classification", "Organization", "Federation"},
			Parsable2: true,
			NoHeader:  true,
		},
		"accounts": {
			Object:    "account",
			Fields:    []string{"Account", "Description", "Organization", "ParentName", "DefQOS", "QOS", "Flags"},
			Parsable2: true,
			NoHeader:  true,
		},
		"users": {
			Object:    "user",
			Fields:    []string{"User", "DefaultAccount", "DefaultWCKey", "AdminLevel", "CoordinatorAccounts"},
			Parsable2: true,
			NoHeader:  true,
		},
		"associations": {
			Object:    "assoc",
			Fields:    []string{"ID", "Cluster", "Account", "User", "Partition", "DefaultQOS", "ParentID", "MaxJobs", "MaxSubmitJobs", "MaxWall", "GrpJobs", "GrpSubmitJobs", "GrpTRES", "MaxTRESPerJob", "MaxTRESPerNode", "MaxTRESMinsPerJob", "Priority", "Shares", "QOS"},
			Parsable2: true,
			NoHeader:  true,
		},
		"qos": {
			Object:    "qos",
			Fields:    []string{"Name", "Description", "Priority", "Flags", "MaxJobsPU", "MaxSubmitPU", "GrpJobs", "GrpSubmitJobs", "GrpTRES", "MaxTRESPerJob", "MaxTRESPU", "MaxWall"},
			Parsable2: true,
			NoHeader:  true,
		},
		"wckeys": {
			Object:    "wckey",
			Fields:    []string{"Wckey", "Cluster", "User", "Default"},
			Parsable2: true,
			NoHeader:  true,
		},
		"tres": {
			Object:    "tres",
			Fields:    []string{"ID", "Type", "Name", "TRESBillingWeights"},
			Parsable2: true,
			NoHeader:  true,
		},
	}
}

func slurmDBVersion() (string, error) {
	lines, err := runCommand("sacctmgr", "--version")
	if err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", nil
	}
	return strings.TrimSpace(lines[0]), nil
}

func runSacctMgrList(profile SlurmDBCommandProfile) ([]string, error) {
	args := []string{"list", profile.Object}
	if profile.NoHeader {
		args = append(args, "--noheader")
	}
	if profile.Parsable2 {
		args = append(args, "--parsable2")
	}
	if len(profile.Fields) > 0 {
		args = append(args, "format="+strings.Join(profile.Fields, ","))
	}
	return runCommand("sacctmgr", args...)
}

func parseSlurmDBText(value string) string {
	v := strings.TrimSpace(value)
	switch strings.ToLower(v) {
	case "", "(null)", "n/a", "none", "null":
		return ""
	default:
		return v
	}
}

func parseSlurmDBList(value string) []string {
	raw := parseSlurmDBText(value)
	if raw == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = parseSlurmDBText(part)
		if part != "" {
			values = append(values, part)
		}
	}
	if len(values) == 0 {
		return []string{}
	}
	return values
}

func parseSlurmDBBool(value string) bool {
	switch strings.ToLower(parseSlurmDBText(value)) {
	case "y", "yes", "true", "1", "default":
		return true
	default:
		return false
	}
}

func parseSlurmDBInt64(value string) *int64 {
	v := parseSlurmDBText(value)
	if v == "" {
		return nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil
	}
	return &n
}

func parseSlurmDBFloat64(value string) *float64 {
	v := parseSlurmDBText(value)
	if v == "" {
		return nil
	}
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &n
}

func slurmDBField(fields []string, index int) string {
	if index < 0 || index >= len(fields) {
		return ""
	}
	return fields[index]
}

// MapSlurmDBClusterRecord converts a fixed-format sacctmgr cluster record to
// the canonical bridge/frontend model.
func MapSlurmDBClusterRecord(fields []string) SlurmDBCluster {
	return SlurmDBCluster{
		Name:           parseSlurmDBText(slurmDBField(fields, 0)),
		ControlHost:    parseSlurmDBText(slurmDBField(fields, 1)),
		Classification: parseSlurmDBText(slurmDBField(fields, 2)),
		Organization:   parseSlurmDBText(slurmDBField(fields, 3)),
		Federation:     parseSlurmDBText(slurmDBField(fields, 4)),
	}
}

// MapSlurmDBAccountRecord converts a fixed-format sacctmgr account record to
// the canonical bridge/frontend model.
func MapSlurmDBAccountRecord(fields []string) SlurmDBAccount {
	return SlurmDBAccount{
		Name:          parseSlurmDBText(slurmDBField(fields, 0)),
		Description:   parseSlurmDBText(slurmDBField(fields, 1)),
		Organization:  parseSlurmDBText(slurmDBField(fields, 2)),
		ParentAccount: parseSlurmDBText(slurmDBField(fields, 3)),
		DefaultQOS:    parseSlurmDBText(slurmDBField(fields, 4)),
		QOSList:       parseSlurmDBList(slurmDBField(fields, 5)),
		Flags:         parseSlurmDBList(slurmDBField(fields, 6)),
	}
}

// MapSlurmDBUserRecord converts a fixed-format sacctmgr user record to the
// canonical bridge/frontend model.
func MapSlurmDBUserRecord(fields []string) SlurmDBUser {
	return SlurmDBUser{
		Name:                parseSlurmDBText(slurmDBField(fields, 0)),
		DefaultAccount:      parseSlurmDBText(slurmDBField(fields, 1)),
		DefaultWckey:        parseSlurmDBText(slurmDBField(fields, 2)),
		AdminLevel:          parseSlurmDBText(slurmDBField(fields, 3)),
		CoordinatorAccounts: parseSlurmDBList(slurmDBField(fields, 4)),
	}
}

// MapSlurmDBAssociationRecord converts a fixed-format sacctmgr assoc record to
// the canonical bridge/frontend model.
func MapSlurmDBAssociationRecord(fields []string) SlurmDBAssociation {
	defaultQOS := parseSlurmDBText(slurmDBField(fields, 5))
	isDefault := defaultQOS != ""
	qosList := parseSlurmDBList(slurmDBField(fields, 18))
	if isDefault && defaultQOS != "" {
		found := false
		for _, qos := range qosList {
			if qos == defaultQOS {
				found = true
				break
			}
		}
		if !found {
			qosList = append([]string{defaultQOS}, qosList...)
		}
	}

	return SlurmDBAssociation{
		ID:                parseSlurmDBInt64(slurmDBField(fields, 0)),
		Cluster:           parseSlurmDBText(slurmDBField(fields, 1)),
		Account:           parseSlurmDBText(slurmDBField(fields, 2)),
		User:              parseSlurmDBText(slurmDBField(fields, 3)),
		Partition:         parseSlurmDBText(slurmDBField(fields, 4)),
		IsDefault:         isDefault,
		ParentID:          parseSlurmDBInt64(slurmDBField(fields, 6)),
		MaxJobs:           parseSlurmDBInt64(slurmDBField(fields, 7)),
		MaxSubmitJobs:     parseSlurmDBInt64(slurmDBField(fields, 8)),
		MaxWall:           parseSlurmDBText(slurmDBField(fields, 9)),
		GrpJobs:           parseSlurmDBInt64(slurmDBField(fields, 10)),
		GrpSubmitJobs:     parseSlurmDBInt64(slurmDBField(fields, 11)),
		GrpTRES:           parseSlurmDBText(slurmDBField(fields, 12)),
		MaxTRESPerJob:     parseSlurmDBText(slurmDBField(fields, 13)),
		MaxTRESPerNode:    parseSlurmDBText(slurmDBField(fields, 14)),
		MaxTRESMinsPerJob: parseSlurmDBText(slurmDBField(fields, 15)),
		Priority:          parseSlurmDBInt64(slurmDBField(fields, 16)),
		SharesRaw:         parseSlurmDBInt64(slurmDBField(fields, 17)),
		QOSList:           qosList,
	}
}

// MapSlurmDBQOSRecord converts a fixed-format sacctmgr qos record to the
// canonical bridge/frontend model.
func MapSlurmDBQOSRecord(fields []string) SlurmDBQOS {
	return SlurmDBQOS{
		Name:                 parseSlurmDBText(slurmDBField(fields, 0)),
		Description:          parseSlurmDBText(slurmDBField(fields, 1)),
		Priority:             parseSlurmDBInt64(slurmDBField(fields, 2)),
		Flags:                parseSlurmDBList(slurmDBField(fields, 3)),
		MaxJobsPerUser:       parseSlurmDBInt64(slurmDBField(fields, 4)),
		MaxSubmitJobsPerUser: parseSlurmDBInt64(slurmDBField(fields, 5)),
		GrpJobs:              parseSlurmDBInt64(slurmDBField(fields, 6)),
		GrpSubmitJobs:        parseSlurmDBInt64(slurmDBField(fields, 7)),
		GrpTRES:              parseSlurmDBText(slurmDBField(fields, 8)),
		MaxTRESPerJob:        parseSlurmDBText(slurmDBField(fields, 9)),
		MaxTRESPerUser:       parseSlurmDBText(slurmDBField(fields, 10)),
		MaxWall:              parseSlurmDBText(slurmDBField(fields, 11)),
	}
}

// MapSlurmDBWckeyRecord converts a fixed-format sacctmgr wckey record to the
// canonical bridge/frontend model.
func MapSlurmDBWckeyRecord(fields []string) SlurmDBWckey {
	return SlurmDBWckey{
		Name:      parseSlurmDBText(slurmDBField(fields, 0)),
		Cluster:   parseSlurmDBText(slurmDBField(fields, 1)),
		User:      parseSlurmDBText(slurmDBField(fields, 2)),
		IsDefault: parseSlurmDBBool(slurmDBField(fields, 3)),
	}
}

// MapSlurmDBTRESRecord converts a fixed-format sacctmgr tres record to the
// canonical bridge/frontend model.
func MapSlurmDBTRESRecord(fields []string) SlurmDBTRES {
	return SlurmDBTRES{
		ID:            parseSlurmDBInt64(slurmDBField(fields, 0)),
		Type:          parseSlurmDBText(slurmDBField(fields, 1)),
		Name:          parseSlurmDBText(slurmDBField(fields, 2)),
		BillingWeight: parseSlurmDBFloat64(slurmDBField(fields, 3)),
	}
}

func legacyClustersFromSlurmDB(records []SlurmDBCluster) []Cluster {
	clusters := make([]Cluster, 0, len(records))
	for _, record := range records {
		clusters = append(clusters, Cluster{
			ClusterName: record.Name,
			ControlHost: record.ControlHost,
		})
	}
	return clusters
}

func legacyUsersFromSlurmDB(records []SlurmDBUser) []User {
	users := make([]User, 0, len(records))
	for _, record := range records {
		users = append(users, User{
			UserID:   record.Name,
			UserName: record.Name,
		})
	}
	return users
}

func legacyAccountsFromSlurmDB(records []SlurmDBAccount) []Account {
	accounts := make([]Account, 0, len(records))
	for _, record := range records {
		accounts = append(accounts, Account{
			AccountID:   record.Name,
			AccountName: record.Name,
			UserNames:   []string{},
		})
	}
	return accounts
}

func collectSlurmDBRecords[T any](profile SlurmDBCommandProfile, mapper func([]string) T) ([]T, error) {
	lines, err := runSacctMgrList(profile)
	if err != nil {
		return []T{}, err
	}

	records := make([]T, 0, len(lines))
	for _, line := range lines {
		records = append(records, mapper(strings.Split(line, "|")))
	}
	return records, nil
}

func getSlurmDBSnapshot() (SlurmDBSnapshot, error) {
	snapshot := newEmptySlurmDBSnapshot()
	collectedAt := time.Now().UTC()
	snapshot.CollectedAt = &collectedAt

	if version, err := slurmDBVersion(); err == nil {
		snapshot.Meta.SacctmgrVersion = version
	} else {
		snapshot.Partial = true
		snapshot.Errors = append(snapshot.Errors, fmt.Sprintf("sacctmgr version: %v", err))
		snapshot.Meta.Errors = append(snapshot.Meta.Errors, snapshot.Errors[len(snapshot.Errors)-1])
	}

	var firstErr error
	collectors := []struct {
		name   string
		target func() error
	}{
		{
			name: "clusters",
			target: func() error {
				records, err := collectSlurmDBRecords(slurmDBCommandProfiles()["clusters"], MapSlurmDBClusterRecord)
				snapshot.Clusters = records
				return err
			},
		},
		{
			name: "accounts",
			target: func() error {
				records, err := collectSlurmDBRecords(slurmDBCommandProfiles()["accounts"], MapSlurmDBAccountRecord)
				snapshot.Accounts = records
				return err
			},
		},
		{
			name: "users",
			target: func() error {
				records, err := collectSlurmDBRecords(slurmDBCommandProfiles()["users"], MapSlurmDBUserRecord)
				snapshot.Users = records
				return err
			},
		},
		{
			name: "associations",
			target: func() error {
				records, err := collectSlurmDBRecords(slurmDBCommandProfiles()["associations"], MapSlurmDBAssociationRecord)
				snapshot.Associations = records
				return err
			},
		},
		{
			name: "qos",
			target: func() error {
				records, err := collectSlurmDBRecords(slurmDBCommandProfiles()["qos"], MapSlurmDBQOSRecord)
				snapshot.QOS = records
				return err
			},
		},
		{
			name: "wckeys",
			target: func() error {
				records, err := collectSlurmDBRecords(slurmDBCommandProfiles()["wckeys"], MapSlurmDBWckeyRecord)
				snapshot.Wckeys = records
				return err
			},
		},
		{
			name: "tres",
			target: func() error {
				records, err := collectSlurmDBRecords(slurmDBCommandProfiles()["tres"], MapSlurmDBTRESRecord)
				snapshot.TRES = records
				return err
			},
		},
	}

	for _, collector := range collectors {
		if err := collector.target(); err != nil {
			msg := fmt.Sprintf("%s: %v", collector.name, err)
			snapshot.Partial = true
			snapshot.Errors = append(snapshot.Errors, msg)
			snapshot.Meta.Errors = append(snapshot.Meta.Errors, msg)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	snapshot.Meta.Partial = snapshot.Partial
	return snapshot, firstErr
}

