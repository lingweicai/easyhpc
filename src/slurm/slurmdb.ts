/*
 * SPDX-License-Identifier: LGPL-2.1-or-later
 */

export interface SlurmDBCommandProfileRecord {
    object: string;
    fields: string[];
    parsable2: boolean;
    noheader: boolean;
}

export interface SlurmDBMetaRecord {
    source?: string;
    sacctmgr_version?: string;
    partial?: boolean;
    errors?: string[];
    command_profiles?: Record<string, SlurmDBCommandProfileRecord>;
}

export interface SlurmDBClusterRecord {
    name: string;
    control_host?: string;
    classification?: string;
    organization?: string;
    federation?: string;
}

export interface SlurmDBAccountRecord {
    name: string;
    description?: string;
    organization?: string;
    parent_account?: string;
    default_qos?: string;
    qos_list?: string[];
    flags?: string[];
}

export interface SlurmDBUserRecord {
    name: string;
    default_account?: string;
    default_wckey?: string;
    admin_level?: string;
    coordinator_accounts?: string[];
}

export interface SlurmDBAssociationRecord {
    id?: number;
    cluster?: string;
    account?: string;
    user?: string;
    partition?: string;
    is_default?: boolean;
    parent_id?: number;
    max_jobs?: number;
    max_submit_jobs?: number;
    max_wall?: string;
    grp_jobs?: number;
    grp_submit_jobs?: number;
    grp_tres?: string;
    max_tres_per_job?: string;
    max_tres_per_node?: string;
    max_tres_mins_per_job?: string;
    priority?: number;
    shares_raw?: number;
    qos_list?: string[];
}

export interface SlurmDBQOSRecord {
    name: string;
    description?: string;
    priority?: number;
    flags?: string[];
    max_jobs_per_user?: number;
    max_submit_jobs_per_user?: number;
    grp_jobs?: number;
    grp_submit_jobs?: number;
    grp_tres?: string;
    max_tres_per_job?: string;
    max_tres_per_user?: string;
    max_wall?: string;
}

export interface SlurmDBWckeyRecord {
    name: string;
    cluster?: string;
    user?: string;
    is_default?: boolean;
}

export interface SlurmDBTRESRecord {
    type?: string;
    name?: string;
    id?: number;
    billing_weight?: number;
}

export interface SlurmDBSnapshotRecord {
    schema_version: string;
    collected_at?: string;
    clusters?: SlurmDBClusterRecord[];
    accounts?: SlurmDBAccountRecord[];
    users?: SlurmDBUserRecord[];
    associations?: SlurmDBAssociationRecord[];
    qos?: SlurmDBQOSRecord[];
    wckeys?: SlurmDBWckeyRecord[];
    tres?: SlurmDBTRESRecord[];
    partial?: boolean;
    errors?: string[];
    meta?: SlurmDBMetaRecord;
}

export interface SlurmDBRecordsResourceRecord<T> {
    schema_version: string;
    collected_at?: string;
    records?: T[];
    partial?: boolean;
    errors?: string[];
    meta?: SlurmDBMetaRecord;
}

export interface SlurmDBCommandProfile {
    object: string;
    fields: string[];
    parsable2: boolean;
    noHeader: boolean;
}

export interface SlurmDBMeta {
    source?: string;
    sacctmgrVersion?: string;
    partial: boolean;
    errors: string[];
    commandProfiles: Record<string, SlurmDBCommandProfile>;
}

export interface SlurmDBCluster {
    name: string;
    controlHost?: string;
    classification?: string;
    organization?: string;
    federation?: string;
}

export interface SlurmDBAccount {
    name: string;
    description?: string;
    organization?: string;
    parentAccount?: string;
    defaultQos?: string;
    qosList: string[];
    flags: string[];
}

export interface SlurmDBUser {
    name: string;
    defaultAccount?: string;
    defaultWckey?: string;
    adminLevel?: string;
    coordinatorAccounts: string[];
}

export interface SlurmDBAssociation {
    id?: number;
    cluster?: string;
    account?: string;
    user?: string;
    partition?: string;
    isDefault: boolean;
    parentId?: number;
    maxJobs?: number;
    maxSubmitJobs?: number;
    maxWall?: string;
    grpJobs?: number;
    grpSubmitJobs?: number;
    grpTres?: string;
    maxTresPerJob?: string;
    maxTresPerNode?: string;
    maxTresMinsPerJob?: string;
    priority?: number;
    sharesRaw?: number;
    qosList: string[];
}

export interface SlurmDBQOS {
    name: string;
    description?: string;
    priority?: number;
    flags: string[];
    maxJobsPerUser?: number;
    maxSubmitJobsPerUser?: number;
    grpJobs?: number;
    grpSubmitJobs?: number;
    grpTres?: string;
    maxTresPerJob?: string;
    maxTresPerUser?: string;
    maxWall?: string;
}

export interface SlurmDBWckey {
    name: string;
    cluster?: string;
    user?: string;
    isDefault: boolean;
}

export interface SlurmDBTRES {
    type?: string;
    name?: string;
    id?: number;
    billingWeight?: number;
}

export interface SlurmDBSnapshot {
    schemaVersion: string;
    collectedAt?: string;
    clusters: SlurmDBCluster[];
    accounts: SlurmDBAccount[];
    users: SlurmDBUser[];
    associations: SlurmDBAssociation[];
    qos: SlurmDBQOS[];
    wckeys: SlurmDBWckey[];
    tres: SlurmDBTRES[];
    partial: boolean;
    errors: string[];
    meta: SlurmDBMeta;
}

export interface SlurmDBRecordsResource<T> {
    schemaVersion: string;
    collectedAt?: string;
    records: T[];
    partial: boolean;
    errors: string[];
    meta: SlurmDBMeta;
}

const optionalString = (value?: string): string | undefined => value ?? undefined;
const optionalNumber = (value?: number): number | undefined => value ?? undefined;
const stringArray = (value?: string[]): string[] => Array.isArray(value) ? [...value] : [];

function normalizeCommandProfile(record: SlurmDBCommandProfileRecord): SlurmDBCommandProfile {
    return {
        object: record.object,
        fields: stringArray(record.fields),
        parsable2: record.parsable2,
        noHeader: record.noheader,
    };
}

export function normalizeSlurmDBMeta(record?: SlurmDBMetaRecord): SlurmDBMeta {
    const commandProfiles = Object.fromEntries(
        Object.entries(record?.command_profiles ?? {}).map(([name, profile]) => [name, normalizeCommandProfile(profile)])
    );

    return {
        ...(record?.source !== undefined ? { source: record.source } : {}),
        ...(record?.sacctmgr_version !== undefined ? { sacctmgrVersion: record.sacctmgr_version } : {}),
        partial: record?.partial ?? false,
        errors: stringArray(record?.errors),
        commandProfiles,
    };
}

export function normalizeSlurmDBCluster(record: SlurmDBClusterRecord): SlurmDBCluster {
    return {
        name: record.name,
        ...(optionalString(record.control_host) !== undefined ? { controlHost: record.control_host } : {}),
        ...(optionalString(record.classification) !== undefined ? { classification: record.classification } : {}),
        ...(optionalString(record.organization) !== undefined ? { organization: record.organization } : {}),
        ...(optionalString(record.federation) !== undefined ? { federation: record.federation } : {}),
    };
}

export function normalizeSlurmDBAccount(record: SlurmDBAccountRecord): SlurmDBAccount {
    return {
        name: record.name,
        ...(optionalString(record.description) !== undefined ? { description: record.description } : {}),
        ...(optionalString(record.organization) !== undefined ? { organization: record.organization } : {}),
        ...(optionalString(record.parent_account) !== undefined ? { parentAccount: record.parent_account } : {}),
        ...(optionalString(record.default_qos) !== undefined ? { defaultQos: record.default_qos } : {}),
        qosList: stringArray(record.qos_list),
        flags: stringArray(record.flags),
    };
}

export function normalizeSlurmDBUser(record: SlurmDBUserRecord): SlurmDBUser {
    return {
        name: record.name,
        ...(optionalString(record.default_account) !== undefined ? { defaultAccount: record.default_account } : {}),
        ...(optionalString(record.default_wckey) !== undefined ? { defaultWckey: record.default_wckey } : {}),
        ...(optionalString(record.admin_level) !== undefined ? { adminLevel: record.admin_level } : {}),
        coordinatorAccounts: stringArray(record.coordinator_accounts),
    };
}

export function normalizeSlurmDBAssociation(record: SlurmDBAssociationRecord): SlurmDBAssociation {
    return {
        ...(optionalNumber(record.id) !== undefined ? { id: record.id } : {}),
        ...(optionalString(record.cluster) !== undefined ? { cluster: record.cluster } : {}),
        ...(optionalString(record.account) !== undefined ? { account: record.account } : {}),
        ...(optionalString(record.user) !== undefined ? { user: record.user } : {}),
        ...(optionalString(record.partition) !== undefined ? { partition: record.partition } : {}),
        isDefault: record.is_default ?? false,
        ...(optionalNumber(record.parent_id) !== undefined ? { parentId: record.parent_id } : {}),
        ...(optionalNumber(record.max_jobs) !== undefined ? { maxJobs: record.max_jobs } : {}),
        ...(optionalNumber(record.max_submit_jobs) !== undefined ? { maxSubmitJobs: record.max_submit_jobs } : {}),
        ...(optionalString(record.max_wall) !== undefined ? { maxWall: record.max_wall } : {}),
        ...(optionalNumber(record.grp_jobs) !== undefined ? { grpJobs: record.grp_jobs } : {}),
        ...(optionalNumber(record.grp_submit_jobs) !== undefined ? { grpSubmitJobs: record.grp_submit_jobs } : {}),
        ...(optionalString(record.grp_tres) !== undefined ? { grpTres: record.grp_tres } : {}),
        ...(optionalString(record.max_tres_per_job) !== undefined ? { maxTresPerJob: record.max_tres_per_job } : {}),
        ...(optionalString(record.max_tres_per_node) !== undefined ? { maxTresPerNode: record.max_tres_per_node } : {}),
        ...(optionalString(record.max_tres_mins_per_job) !== undefined ? { maxTresMinsPerJob: record.max_tres_mins_per_job } : {}),
        ...(optionalNumber(record.priority) !== undefined ? { priority: record.priority } : {}),
        ...(optionalNumber(record.shares_raw) !== undefined ? { sharesRaw: record.shares_raw } : {}),
        qosList: stringArray(record.qos_list),
    };
}

export function normalizeSlurmDBQOS(record: SlurmDBQOSRecord): SlurmDBQOS {
    return {
        name: record.name,
        ...(optionalString(record.description) !== undefined ? { description: record.description } : {}),
        ...(optionalNumber(record.priority) !== undefined ? { priority: record.priority } : {}),
        flags: stringArray(record.flags),
        ...(optionalNumber(record.max_jobs_per_user) !== undefined ? { maxJobsPerUser: record.max_jobs_per_user } : {}),
        ...(optionalNumber(record.max_submit_jobs_per_user) !== undefined ? { maxSubmitJobsPerUser: record.max_submit_jobs_per_user } : {}),
        ...(optionalNumber(record.grp_jobs) !== undefined ? { grpJobs: record.grp_jobs } : {}),
        ...(optionalNumber(record.grp_submit_jobs) !== undefined ? { grpSubmitJobs: record.grp_submit_jobs } : {}),
        ...(optionalString(record.grp_tres) !== undefined ? { grpTres: record.grp_tres } : {}),
        ...(optionalString(record.max_tres_per_job) !== undefined ? { maxTresPerJob: record.max_tres_per_job } : {}),
        ...(optionalString(record.max_tres_per_user) !== undefined ? { maxTresPerUser: record.max_tres_per_user } : {}),
        ...(optionalString(record.max_wall) !== undefined ? { maxWall: record.max_wall } : {}),
    };
}

export function normalizeSlurmDBWckey(record: SlurmDBWckeyRecord): SlurmDBWckey {
    return {
        name: record.name,
        ...(optionalString(record.cluster) !== undefined ? { cluster: record.cluster } : {}),
        ...(optionalString(record.user) !== undefined ? { user: record.user } : {}),
        isDefault: record.is_default ?? false,
    };
}

export function normalizeSlurmDBTRES(record: SlurmDBTRESRecord): SlurmDBTRES {
    return {
        ...(optionalString(record.type) !== undefined ? { type: record.type } : {}),
        ...(optionalString(record.name) !== undefined ? { name: record.name } : {}),
        ...(optionalNumber(record.id) !== undefined ? { id: record.id } : {}),
        ...(optionalNumber(record.billing_weight) !== undefined ? { billingWeight: record.billing_weight } : {}),
    };
}

export function normalizeSlurmDBSnapshot(record: SlurmDBSnapshotRecord): SlurmDBSnapshot {
    return {
        schemaVersion: record.schema_version,
        ...(optionalString(record.collected_at) !== undefined ? { collectedAt: record.collected_at } : {}),
        clusters: (record.clusters ?? []).map(normalizeSlurmDBCluster),
        accounts: (record.accounts ?? []).map(normalizeSlurmDBAccount),
        users: (record.users ?? []).map(normalizeSlurmDBUser),
        associations: (record.associations ?? []).map(normalizeSlurmDBAssociation),
        qos: (record.qos ?? []).map(normalizeSlurmDBQOS),
        wckeys: (record.wckeys ?? []).map(normalizeSlurmDBWckey),
        tres: (record.tres ?? []).map(normalizeSlurmDBTRES),
        partial: record.partial ?? record.meta?.partial ?? false,
        errors: stringArray(record.errors ?? record.meta?.errors),
        meta: normalizeSlurmDBMeta(record.meta),
    };
}

export function normalizeSlurmDBRecordsResource<TRecord, TModel>(
    record: SlurmDBRecordsResourceRecord<TRecord>,
    normalizeRecord: (raw: TRecord) => TModel
): SlurmDBRecordsResource<TModel> {
    return {
        schemaVersion: record.schema_version,
        ...(optionalString(record.collected_at) !== undefined ? { collectedAt: record.collected_at } : {}),
        records: (record.records ?? []).map(normalizeRecord),
        partial: record.partial ?? record.meta?.partial ?? false,
        errors: stringArray(record.errors ?? record.meta?.errors),
        meta: normalizeSlurmDBMeta(record.meta),
    };
}
