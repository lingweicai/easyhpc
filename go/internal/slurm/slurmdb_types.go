// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

import "time"

const slurmDBSchemaVersion = "1"

// SlurmDBCommandProfile describes the fixed sacctmgr output profile used to
// collect a resource from slurmdbd.
type SlurmDBCommandProfile struct {
	Object    string   `json:"object"`
	Fields    []string `json:"fields"`
	Parsable2 bool     `json:"parsable2"`
	NoHeader  bool     `json:"noheader"`
}

// SlurmDBMeta contains snapshot collection metadata.
type SlurmDBMeta struct {
	Source          string                          `json:"source"`
	SacctmgrVersion string                          `json:"sacctmgr_version,omitempty"`
	Partial         bool                            `json:"partial"`
	Errors          []string                        `json:"errors"`
	CommandProfiles map[string]SlurmDBCommandProfile `json:"command_profiles"`
}

// SlurmDBCluster is the canonical SlurmDB cluster record.
type SlurmDBCluster struct {
	Name           string `json:"name"`
	ControlHost    string `json:"control_host,omitempty"`
	Classification string `json:"classification,omitempty"`
	Organization   string `json:"organization,omitempty"`
	Federation     string `json:"federation,omitempty"`
}

// SlurmDBAccount is the canonical SlurmDB account record.
type SlurmDBAccount struct {
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Organization  string   `json:"organization,omitempty"`
	ParentAccount string   `json:"parent_account,omitempty"`
	DefaultQOS    string   `json:"default_qos,omitempty"`
	QOSList       []string `json:"qos_list"`
	Flags         []string `json:"flags"`
}

// SlurmDBUser is the canonical SlurmDB user record.
type SlurmDBUser struct {
	Name                string   `json:"name"`
	DefaultAccount      string   `json:"default_account,omitempty"`
	DefaultWckey        string   `json:"default_wckey,omitempty"`
	AdminLevel          string   `json:"admin_level,omitempty"`
	CoordinatorAccounts []string `json:"coordinator_accounts"`
}

// SlurmDBAssociation is the canonical SlurmDB association record.
type SlurmDBAssociation struct {
	ID                  *int64   `json:"id,omitempty"`
	Cluster             string   `json:"cluster,omitempty"`
	Account             string   `json:"account,omitempty"`
	User                string   `json:"user,omitempty"`
	Partition           string   `json:"partition,omitempty"`
	IsDefault           bool     `json:"is_default"`
	ParentID            *int64   `json:"parent_id,omitempty"`
	MaxJobs             *int64   `json:"max_jobs,omitempty"`
	MaxSubmitJobs       *int64   `json:"max_submit_jobs,omitempty"`
	MaxWall             string   `json:"max_wall,omitempty"`
	GrpJobs             *int64   `json:"grp_jobs,omitempty"`
	GrpSubmitJobs       *int64   `json:"grp_submit_jobs,omitempty"`
	GrpTRES             string   `json:"grp_tres,omitempty"`
	MaxTRESPerJob       string   `json:"max_tres_per_job,omitempty"`
	MaxTRESPerNode      string   `json:"max_tres_per_node,omitempty"`
	MaxTRESMinsPerJob   string   `json:"max_tres_mins_per_job,omitempty"`
	Priority            *int64   `json:"priority,omitempty"`
	SharesRaw           *int64   `json:"shares_raw,omitempty"`
	QOSList             []string `json:"qos_list"`
}

// SlurmDBQOS is the canonical SlurmDB QOS record.
type SlurmDBQOS struct {
	Name                 string   `json:"name"`
	Description          string   `json:"description,omitempty"`
	Priority             *int64   `json:"priority,omitempty"`
	Flags                []string `json:"flags"`
	MaxJobsPerUser       *int64   `json:"max_jobs_per_user,omitempty"`
	MaxSubmitJobsPerUser *int64   `json:"max_submit_jobs_per_user,omitempty"`
	GrpJobs              *int64   `json:"grp_jobs,omitempty"`
	GrpSubmitJobs        *int64   `json:"grp_submit_jobs,omitempty"`
	GrpTRES              string   `json:"grp_tres,omitempty"`
	MaxTRESPerJob        string   `json:"max_tres_per_job,omitempty"`
	MaxTRESPerUser       string   `json:"max_tres_per_user,omitempty"`
	MaxWall              string   `json:"max_wall,omitempty"`
}

// SlurmDBWckey is the canonical SlurmDB WCKey record.
type SlurmDBWckey struct {
	Name      string `json:"name"`
	Cluster   string `json:"cluster,omitempty"`
	User      string `json:"user,omitempty"`
	IsDefault bool   `json:"is_default"`
}

// SlurmDBTRES is the canonical SlurmDB TRES record.
type SlurmDBTRES struct {
	Type          string   `json:"type,omitempty"`
	Name          string   `json:"name,omitempty"`
	ID            *int64   `json:"id,omitempty"`
	BillingWeight *float64 `json:"billing_weight,omitempty"`
}

// SlurmDBSnapshot is the canonical SlurmDB snapshot resource emitted by the
// bridge and consumed by Cockpit channel clients.
type SlurmDBSnapshot struct {
	SchemaVersion string              `json:"schema_version"`
	CollectedAt   *time.Time          `json:"collected_at,omitempty"`
	Clusters      []SlurmDBCluster    `json:"clusters"`
	Accounts      []SlurmDBAccount    `json:"accounts"`
	Users         []SlurmDBUser       `json:"users"`
	Associations  []SlurmDBAssociation `json:"associations"`
	QOS           []SlurmDBQOS        `json:"qos"`
	Wckeys        []SlurmDBWckey      `json:"wckeys"`
	TRES          []SlurmDBTRES       `json:"tres"`
	Partial       bool                `json:"partial"`
	Errors        []string            `json:"errors"`
	Meta          SlurmDBMeta         `json:"meta"`
}

// SlurmDBRecordsResource is a targeted SlurmDB resource wrapper for a single
// collection (accounts, users, qos, ...).
type SlurmDBRecordsResource struct {
	SchemaVersion string      `json:"schema_version"`
	CollectedAt   *time.Time  `json:"collected_at,omitempty"`
	Records       interface{} `json:"records"`
	Partial       bool        `json:"partial"`
	Errors        []string    `json:"errors"`
	Meta          SlurmDBMeta `json:"meta"`
}

func newEmptySlurmDBSnapshot() SlurmDBSnapshot {
	return SlurmDBSnapshot{
		SchemaVersion: slurmDBSchemaVersion,
		Clusters:      []SlurmDBCluster{},
		Accounts:      []SlurmDBAccount{},
		Users:         []SlurmDBUser{},
		Associations:  []SlurmDBAssociation{},
		QOS:           []SlurmDBQOS{},
		Wckeys:        []SlurmDBWckey{},
		TRES:          []SlurmDBTRES{},
		Errors:        []string{},
		Meta: SlurmDBMeta{
			Source: "sacctmgr",
			Errors: []string{},
			CommandProfiles: slurmDBCommandProfiles(),
		},
	}
}

func newSlurmDBRecordsResource(snapshot SlurmDBSnapshot, records interface{}) SlurmDBRecordsResource {
	return SlurmDBRecordsResource{
		SchemaVersion: snapshot.SchemaVersion,
		CollectedAt:   snapshot.CollectedAt,
		Records:       records,
		Partial:       snapshot.Partial,
		Errors:        append([]string{}, snapshot.Errors...),
		Meta:          snapshot.Meta,
	}
}

