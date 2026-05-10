/*
 * SPDX-License-Identifier: LGPL-2.1-or-later
 *
 * Copyright (C) 2024 Red Hat, Inc.
 */

import React, { useCallback, useEffect, useRef, useState } from 'react';

import { Alert } from "@patternfly/react-core/dist/esm/components/Alert/index.js";
import { Badge } from "@patternfly/react-core/dist/esm/components/Badge/index.js";
import { Card, CardBody, CardTitle } from "@patternfly/react-core/dist/esm/components/Card/index.js";
import { EmptyState, EmptyStateBody } from "@patternfly/react-core/dist/esm/components/EmptyState/index.js";
import { Label } from "@patternfly/react-core/dist/esm/components/Label/index.js";
import { Spinner } from "@patternfly/react-core/dist/esm/components/Spinner/index.js";
import { Tab, TabContent, TabTitleText, Tabs } from "@patternfly/react-core/dist/esm/components/Tabs/index.js";
import { Grid, GridItem } from "@patternfly/react-core/dist/esm/layouts/Grid/index.js";

import cockpit from 'cockpit';

const _ = cockpit.gettext;

// ── Types ──────────────────────────────────────────────────────────────────

interface SlurmNode {
    name: string;
    state: string;
    partition: string;
    cpus: string;
    memory: string;
    reason?: string;
}

interface SlurmPartition {
    name: string;
    state: string;
    total_cpus: string;
    total_nodes: string;
}

interface SlurmCluster {
    name: string;
    control_host: string;
    control_port: string;
}

interface SlurmJob {
    job_id: string;
    name: string;
    user: string;
    account: string;
    state: string;
    partition: string;
    nodes: string;
    time: string;
    time_limit: string;
    reason?: string;
}

interface SlurmReservation {
    name: string;
    state: string;
    nodes: string;
    start_time: string;
    end_time: string;
    duration: string;
    users: string;
    accounts: string;
}

interface SlurmUser {
    name: string;
    default_account: string;
    admin: string;
}

interface SlurmAccount {
    name: string;
    description: string;
    organization: string;
}

interface LogEntry {
    id: number;
    level: string;
    message: string;
    timestamp: string;
}

interface BridgeMessage {
    type: string;
    resource?: string;
    data?: unknown;
    message?: string;
    level?: string;
    timestamp?: string;
    version?: string;
}

type ConnectionStatus = 'connecting' | 'ready' | 'error' | 'closed';

// ── State colour helpers ───────────────────────────────────────────────────

function nodeStateColor(state: string): 'green' | 'red' | 'orange' | 'grey' {
    const s = state.toLowerCase();
    if (s.startsWith('idle')) return 'green';
    if (s.startsWith('down') || s.startsWith('fail')) return 'red';
    if (s.startsWith('drain')) return 'orange';
    return 'grey';
}

function jobStateColor(state: string): 'green' | 'blue' | 'orange' | 'grey' {
    const s = state.toUpperCase();
    if (s === 'RUNNING') return 'green';
    if (s === 'PENDING') return 'blue';
    if (s === 'COMPLETING') return 'orange';
    return 'grey';
}

// ── Simple table helper ────────────────────────────────────────────────────

interface TableProps {
    columns: string[];
    rows: React.ReactNode[][];
}

const SimpleTable = ({ columns, rows }: TableProps) => (
    <table className="pf-v6-c-table pf-m-compact pf-m-grid-sm">
        <thead className="pf-v6-c-table__thead">
            <tr className="pf-v6-c-table__tr">
                {columns.map(col => (
                    <th key={col} className="pf-v6-c-table__th">{col}</th>
                ))}
            </tr>
        </thead>
        <tbody className="pf-v6-c-table__tbody">
            {rows.length === 0
                ? (
                    <tr className="pf-v6-c-table__tr">
                        <td className="pf-v6-c-table__td" colSpan={columns.length}>
                            <EmptyState>
                                <EmptyStateBody>{_("No data available")}</EmptyStateBody>
                            </EmptyState>
                        </td>
                    </tr>
                )
                : rows.map((row, i) => (
                    <tr key={i} className="pf-v6-c-table__tr">
                        {row.map((cell, j) => (
                            <td key={j} className="pf-v6-c-table__td">{cell}</td>
                        ))}
                    </tr>
                ))}
        </tbody>
    </table>
);

// ── Bridge hook ────────────────────────────────────────────────────────────

// The bridge binary is installed alongside the other package files.
// Users can also place it anywhere on $PATH.
const BRIDGE_BINARY = "easyhpc-bridge";

function useSlurmBridge() {
    const [status, setStatus] = useState<ConnectionStatus>('connecting');
    const [error, setError] = useState<string | null>(null);
    const [clusters, setClusters] = useState<SlurmCluster[]>([]);
    const [partitions, setPartitions] = useState<SlurmPartition[]>([]);
    const [nodes, setNodes] = useState<SlurmNode[]>([]);
    const [jobs, setJobs] = useState<SlurmJob[]>([]);
    const [reservations, setReservations] = useState<SlurmReservation[]>([]);
    const [users, setUsers] = useState<SlurmUser[]>([]);
    const [accounts, setAccounts] = useState<SlurmAccount[]>([]);
    const [events, setEvents] = useState<LogEntry[]>([]);
    const [lastUpdated, setLastUpdated] = useState<string | null>(null);

    const bridgeRef = useRef<ReturnType<typeof cockpit.spawn> | null>(null);
    const bufferRef = useRef<string>('');
    const eventIdRef = useRef<number>(0);

    const handleMessage = useCallback((msg: BridgeMessage) => {
        switch (msg.type) {
        case 'ready':
            setStatus('ready');
            break;
        case 'data':
            setLastUpdated(msg.timestamp ?? null);
            switch (msg.resource) {
            case 'clusters':
                setClusters((msg.data as SlurmCluster[]) ?? []);
                break;
            case 'partitions':
                setPartitions((msg.data as SlurmPartition[]) ?? []);
                break;
            case 'nodes':
                setNodes((msg.data as SlurmNode[]) ?? []);
                break;
            case 'jobs':
                setJobs((msg.data as SlurmJob[]) ?? []);
                break;
            case 'reservations':
                setReservations((msg.data as SlurmReservation[]) ?? []);
                break;
            case 'users':
                setUsers((msg.data as SlurmUser[]) ?? []);
                break;
            case 'accounts':
                setAccounts((msg.data as SlurmAccount[]) ?? []);
                break;
            }
            break;
        case 'event':
            setEvents(prev => {
                const entry: LogEntry = {
                    id: ++eventIdRef.current,
                    level: msg.level ?? 'info',
                    message: msg.message ?? '',
                    timestamp: msg.timestamp ?? new Date().toISOString(),
                };
                return [entry, ...prev].slice(0, 200);
            });
            break;
        case 'error':
            setError(msg.message ?? 'Unknown error');
            break;
        }
    }, []);

    useEffect(() => {
        const bridge = cockpit.spawn([BRIDGE_BINARY], {
            err: "message",
            superuser: "try",
        });
        bridgeRef.current = bridge;

        bridge.stream((chunk: string) => {
            bufferRef.current += chunk;
            const lines = bufferRef.current.split('\n');
            bufferRef.current = lines.pop() ?? '';
            for (const line of lines) {
                const trimmed = line.trim();
                if (!trimmed) continue;
                try {
                    const msg: BridgeMessage = JSON.parse(trimmed);
                    handleMessage(msg);
                } catch {
                    /* ignore malformed lines */
                }
            }
        });

        bridge.fail((err: Error) => {
            setStatus('error');
            setError(_("Bridge process failed: ") + err.message);
        });

        bridge.done(() => {
            setStatus('closed');
        });

        // Ask for an immediate full snapshot.
        bridge.input(JSON.stringify({ type: 'poll', resource: 'all' }) + '\n');

        return () => {
            bridge.close('disconnecting');
            bridgeRef.current = null;
        };
    }, [handleMessage]);

    const refresh = useCallback(() => {
        bridgeRef.current?.input(JSON.stringify({ type: 'refresh' }) + '\n');
    }, []);

    return {
        status,
        error,
        clusters,
        partitions,
        nodes,
        jobs,
        reservations,
        users,
        accounts,
        events,
        lastUpdated,
        refresh,
    };
}

// ── Tab content components ─────────────────────────────────────────────────

const OverviewTab = ({ clusters, partitions, nodes, jobs, reservations, users, accounts }: {
    clusters: SlurmCluster[];
    partitions: SlurmPartition[];
    nodes: SlurmNode[];
    jobs: SlurmJob[];
    reservations: SlurmReservation[];
    users: SlurmUser[];
    accounts: SlurmAccount[];
}) => {
    const stats = [
        { label: _("Clusters"), value: clusters.length },
        { label: _("Partitions"), value: partitions.length },
        { label: _("Nodes"), value: nodes.length },
        { label: _("Jobs"), value: jobs.length },
        { label: _("Reservations"), value: reservations.length },
        { label: _("Users"), value: users.length },
        { label: _("Accounts"), value: accounts.length },
    ];

    return (
        <Grid hasGutter>
            {stats.map(s => (
                <GridItem key={s.label} span={3}>
                    <Card>
                        <CardTitle>{s.label}</CardTitle>
                        <CardBody>
                            <Badge screenReaderText="">{s.value}</Badge>
                        </CardBody>
                    </Card>
                </GridItem>
            ))}
        </Grid>
    );
};

const ClustersTab = ({ clusters }: { clusters: SlurmCluster[] }) => (
    <SimpleTable
        columns={[_("Name"), _("Control Host"), _("Control Port")]}
        rows={clusters.map(c => [c.name, c.control_host, c.control_port])}
    />
);

const PartitionsTab = ({ partitions }: { partitions: SlurmPartition[] }) => (
    <SimpleTable
        columns={[_("Name"), _("State"), _("Total CPUs"), _("Total Nodes")]}
        rows={partitions.map(p => [
            p.name,
            <Label key={p.name} color={p.state === 'up' ? 'green' : 'red'}>{p.state}</Label>,
            p.total_cpus,
            p.total_nodes,
        ])}
    />
);

const NodesTab = ({ nodes }: { nodes: SlurmNode[] }) => (
    <SimpleTable
        columns={[_("Name"), _("State"), _("Partition"), _("CPUs"), _("Memory (MB)"), _("Reason")]}
        rows={nodes.map(n => [
            n.name,
            <Label key={n.name} color={nodeStateColor(n.state)}>{n.state}</Label>,
            n.partition,
            n.cpus,
            n.memory,
            n.reason ?? '',
        ])}
    />
);

const JobsTab = ({ jobs }: { jobs: SlurmJob[] }) => (
    <SimpleTable
        columns={[_("Job ID"), _("Name"), _("User"), _("Account"), _("State"), _("Partition"), _("Nodes"), _("Time"), _("Time Limit")]}
        rows={jobs.map(j => [
            j.job_id,
            j.name,
            j.user,
            j.account,
            <Label key={j.job_id} color={jobStateColor(j.state)}>{j.state}</Label>,
            j.partition,
            j.nodes,
            j.time,
            j.time_limit,
        ])}
    />
);

const ReservationsTab = ({ reservations }: { reservations: SlurmReservation[] }) => (
    <SimpleTable
        columns={[_("Name"), _("State"), _("Nodes"), _("Start Time"), _("End Time"), _("Users"), _("Accounts")]}
        rows={reservations.map(r => [
            r.name,
            r.state,
            r.nodes,
            r.start_time,
            r.end_time,
            r.users,
            r.accounts,
        ])}
    />
);

const UsersAccountsTab = ({ users, accounts }: { users: SlurmUser[]; accounts: SlurmAccount[] }) => (
    <Grid hasGutter>
        <GridItem span={6}>
            <Card>
                <CardTitle>{_("Users")}</CardTitle>
                <CardBody>
                    <SimpleTable
                        columns={[_("Name"), _("Default Account"), _("Admin")]}
                        rows={users.map(u => [u.name, u.default_account, u.admin])}
                    />
                </CardBody>
            </Card>
        </GridItem>
        <GridItem span={6}>
            <Card>
                <CardTitle>{_("Accounts")}</CardTitle>
                <CardBody>
                    <SimpleTable
                        columns={[_("Name"), _("Description"), _("Organization")]}
                        rows={accounts.map(a => [a.name, a.description, a.organization])}
                    />
                </CardBody>
            </Card>
        </GridItem>
    </Grid>
);

const EventsTab = ({ events }: { events: LogEntry[] }) => {
    if (events.length === 0) {
        return (
            <EmptyState>
                <EmptyStateBody>
                    {_("No log events received yet.  Events are streamed from the slurmctld log.")}
                </EmptyStateBody>
            </EmptyState>
        );
    }
    return (
        <div className="ehpc-events-list">
            {events.map(e => (
                <Alert
                    key={e.id}
                    variant={e.level === 'error' ? 'danger' : e.level === 'warning' ? 'warning' : 'info'}
                    title={e.timestamp}
                    isInline
                    className="ehpc-event-alert"
                >
                    {e.message}
                </Alert>
            ))}
        </div>
    );
};

// ── Main Application ───────────────────────────────────────────────────────

export const Application = () => {
    const [activeTab, setActiveTab] = useState<string | number>('overview');
    const {
        status, error,
        clusters, partitions, nodes, jobs, reservations, users, accounts, events,
        lastUpdated,
        refresh,
    } = useSlurmBridge();

    const tabs = [
        { key: 'overview', title: _("Overview") },
        { key: 'clusters', title: _("Clusters") },
        { key: 'partitions', title: _("Partitions") },
        { key: 'nodes', title: _("Nodes") },
        { key: 'jobs', title: _("Jobs") },
        { key: 'reservations', title: _("Reservations") },
        { key: 'users-accounts', title: _("Users & Accounts") },
        { key: 'events', title: _("Log Events") },
    ];

    return (
        <div className="ehpc-app">
            <div className="ehpc-header">
                <div className="ehpc-header-title">
                    <h1>{_("EasyHPC – Slurm Dashboard")}</h1>
                    {lastUpdated && (
                        <span className="ehpc-last-updated">
                            {cockpit.format(_("Last updated: $0"), new Date(lastUpdated).toLocaleTimeString())}
                        </span>
                    )}
                </div>
                <div className="ehpc-header-status">
                    {status === 'connecting' && <><Spinner size="sm" />&nbsp;{_("Connecting…")}</>}
                    {status === 'ready' && (
                        <button className="pf-v6-c-button pf-m-secondary pf-m-small" onClick={refresh}>
                            {_("Refresh")}
                        </button>
                    )}
                    {status === 'error' && <Label color="red">{_("Bridge error")}</Label>}
                    {status === 'closed' && <Label color="orange">{_("Bridge stopped")}</Label>}
                </div>
            </div>

            {error && (
                <Alert
                    variant="warning"
                    title={_("Slurm bridge warning")}
                    isInline
                    className="ehpc-bridge-alert"
                >
                    {error}
                </Alert>
            )}

            <Tabs
                activeKey={activeTab}
                onSelect={(_ev, key) => setActiveTab(key)}
                aria-label={_("Slurm resource tabs")}
            >
                {tabs.map(t => (
                    <Tab key={t.key} eventKey={t.key} title={<TabTitleText>{t.title}</TabTitleText>} />
                ))}
            </Tabs>

            <TabContent className="ehpc-tab-content">
                {activeTab === 'overview' && (
                    <OverviewTab {...{ clusters, partitions, nodes, jobs, reservations, users, accounts }} />
                )}
                {activeTab === 'clusters' && <ClustersTab clusters={clusters} />}
                {activeTab === 'partitions' && <PartitionsTab partitions={partitions} />}
                {activeTab === 'nodes' && <NodesTab nodes={nodes} />}
                {activeTab === 'jobs' && <JobsTab jobs={jobs} />}
                {activeTab === 'reservations' && <ReservationsTab reservations={reservations} />}
                {activeTab === 'users-accounts' && <UsersAccountsTab users={users} accounts={accounts} />}
                {activeTab === 'events' && <EventsTab events={events} />}
            </TabContent>
        </div>
    );
};
