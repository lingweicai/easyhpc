import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';
import ts from 'typescript';

async function loadSlurmDBModule() {
    const rootDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
    const filePath = path.join(rootDir, 'src/slurm/slurmdb.ts');
    const source = fs.readFileSync(filePath, 'utf8');
    const transpiled = ts.transpileModule(source, {
        compilerOptions: {
            module: ts.ModuleKind.ES2020,
            target: ts.ScriptTarget.ES2020,
        },
        fileName: filePath,
    });

    const url = `data:text/javascript;base64,${Buffer.from(transpiled.outputText, 'utf8').toString('base64')}`;
    return import(url);
}

test('normalizeSlurmDBSnapshot maps snake_case records to camelCase models', async () => {
    const slurmdb = await loadSlurmDBModule();

    const snapshot = slurmdb.normalizeSlurmDBSnapshot({
        schema_version: '1',
        collected_at: '2026-05-10T00:00:00Z',
        clusters: [{ name: 'alpha', control_host: 'ctl01', classification: 'production' }],
        accounts: [{ name: 'research', parent_account: 'root', default_qos: 'normal', qos_list: ['normal'], flags: ['relative'] }],
        users: [{ name: 'alice', default_account: 'research', coordinator_accounts: ['research'] }],
        associations: [{ id: 42, cluster: 'alpha', account: 'research', user: 'alice', is_default: true, qos_list: ['normal'] }],
        qos: [{ name: 'normal', priority: 100, flags: ['usage_factor_safe'] }],
        wckeys: [{ name: 'project-x', cluster: 'alpha', user: 'alice', is_default: true }],
        tres: [{ id: 1, type: 'cpu', name: 'cpu', billing_weight: 1.25 }],
        partial: true,
        errors: ['tres: unsupported field'],
        meta: {
            source: 'sacctmgr',
            sacctmgr_version: 'slurm 24.05.0',
            partial: true,
            errors: ['tres: unsupported field'],
            command_profiles: {
                clusters: {
                    object: 'cluster',
                    fields: ['Cluster', 'ControlHost'],
                    parsable2: true,
                    noheader: true,
                },
            },
        },
    });

    assert.equal(snapshot.schemaVersion, '1');
    assert.equal(snapshot.collectedAt, '2026-05-10T00:00:00Z');
    assert.equal(snapshot.clusters[0].controlHost, 'ctl01');
    assert.equal(snapshot.accounts[0].parentAccount, 'root');
    assert.equal(snapshot.accounts[0].defaultQos, 'normal');
    assert.deepEqual(snapshot.users[0].coordinatorAccounts, ['research']);
    assert.equal(snapshot.associations[0].isDefault, true);
    assert.equal(snapshot.qos[0].priority, 100);
    assert.equal(snapshot.wckeys[0].isDefault, true);
    assert.equal(snapshot.tres[0].billingWeight, 1.25);
    assert.equal(snapshot.partial, true);
    assert.deepEqual(snapshot.errors, ['tres: unsupported field']);
    assert.equal(snapshot.meta.sacctmgrVersion, 'slurm 24.05.0');
    assert.equal(snapshot.meta.commandProfiles.clusters.noHeader, true);
});

test('normalizeSlurmDBSnapshot remains backward-compatible when newer fields are absent', async () => {
    const slurmdb = await loadSlurmDBModule();

    const snapshot = slurmdb.normalizeSlurmDBSnapshot({
        schema_version: '1',
        clusters: [{ name: 'alpha' }],
    });

    assert.deepEqual(snapshot.accounts, []);
    assert.deepEqual(snapshot.users, []);
    assert.deepEqual(snapshot.associations, []);
    assert.deepEqual(snapshot.qos, []);
    assert.deepEqual(snapshot.wckeys, []);
    assert.deepEqual(snapshot.tres, []);
    assert.equal(snapshot.partial, false);
    assert.deepEqual(snapshot.errors, []);
    assert.deepEqual(snapshot.meta.errors, []);
    assert.deepEqual(snapshot.meta.commandProfiles, {});
});

test('normalizeSlurmDBRecordsResource normalizes targeted records envelopes', async () => {
    const slurmdb = await loadSlurmDBModule();

    const resource = slurmdb.normalizeSlurmDBRecordsResource({
        schema_version: '1',
        records: [
            { name: 'research', default_qos: 'normal', qos_list: ['normal', 'burst'], flags: ['relative'] },
        ],
        meta: {
            source: 'sacctmgr',
            command_profiles: {
                accounts: {
                    object: 'account',
                    fields: ['Account', 'DefQOS', 'QOS'],
                    parsable2: true,
                    noheader: true,
                },
            },
        },
    }, slurmdb.normalizeSlurmDBAccount);

    assert.equal(resource.schemaVersion, '1');
    assert.equal(resource.records[0].defaultQos, 'normal');
    assert.deepEqual(resource.records[0].qosList, ['normal', 'burst']);
    assert.deepEqual(resource.records[0].flags, ['relative']);
    assert.equal(resource.meta.commandProfiles.accounts.object, 'account');
});
