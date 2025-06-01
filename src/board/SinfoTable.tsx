import React from 'react';
import {
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td,
} from '@patternfly/react-table';
import { Spinner, Alert, Label } from '@patternfly/react-core';
import { useSinfoContext } from '../SinfoContext';
import cockpit from 'cockpit';

const _ = cockpit.gettext;

const SinfoTable: React.FC = () => {
    const { rawOutput, loading, error } = useSinfoContext();

    if (loading) return <Spinner />;
    if (error) return <Alert variant="danger" title="Failed to load sinfo">{error}</Alert>;

    const rows = parseSinfoOutput(rawOutput);
    type LabelColor = | "blue" | "teal" | "green" | "orange"
  | "purple" | "red" | "orangered" | "grey" | "yellow";

    const getStateColor = (state: string): LabelColor => {
        switch (state.toLowerCase()) {
        case "idle":
            return "green";
        case "alloc":
            return "blue";
        case "mixed":
            return "orange";
        case "down":
            return "red";
        case "drain":
            return "purple";
        case "maint":
            return "teal";
        default:
            return "grey"; // fallback is still valid color
        }
    };

    return (
        <Table aria-label="sinfo output table" variant="compact">
            <Thead>
                <Tr>
                    <Th>{_("PARTITION")}</Th>
                    <Th>{_("AVAIL")}</Th>
                    <Th>{_("TIMELIMIT")}</Th>
                    <Th>{_("NODES")}</Th>
                    <Th>{_("STATE")}</Th>
                    <Th>{_("NODELIST")}</Th>
                </Tr>
            </Thead>
            <Tbody>
                {rows.map((row, index) => (
                    <Tr key={index}>
                        <Td>{row.partition}</Td>
                        <Td>{row.avail}</Td>
                        <Td>{row.timelimit}</Td>
                        <Td>{row.nodes}</Td>
                        <Td>
                            <Label color={getStateColor(row.state)}>{row.state}</Label>
                        </Td>
                        <Td>{row.nodelist}</Td>
                    </Tr>
                ))}
            </Tbody>
        </Table>
    );
};

// Helper: parse sinfo output into structured data
function parseSinfoOutput(raw: string): {
  partition: string;
  avail: string;
  timelimit: string;
  nodes: string;
  state: string;
  nodelist: string;
}[] {
    const lines = raw.trim().split('\n');
    if (lines.length < 2) return [];

    return lines.slice(1).map(line => {
        const [partition, avail, timelimit, nodes, state, ...nodelistParts] = line.trim().split(/\s+/);
        return {
            partition,
            avail,
            timelimit,
            nodes,
            state,
            nodelist: nodelistParts.join(' ')
        };
    });
}

export default SinfoTable;
