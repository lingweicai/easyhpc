import React, { useEffect, useState } from 'react';
import {
    Card,
    CardTitle,
    CardBody,
} from '@patternfly/react-core';

import {
    Table,
    Thead,
    Tr,
    Th,
    Tbody,
    Td,
    ExpandableRowContent
} from '@patternfly/react-table';
import cockpit from 'cockpit';

const _ = cockpit.gettext;

type PartitionInfo = {
  PartitionName: string;
  Default: string;
  MaxTime: string;
  Nodes: string;
  State: string;
  TotalCPUs: string;
  TotalNodes: string;
  extra: string;
};

const parsePartitionsOutput = (output: string): PartitionInfo[] => {
    const partitions: PartitionInfo[] = [];
    const blocks = output.trim().split(/\n(?=PartitionName=)/);

    blocks.forEach(block => {
        const lines = block.trim().split("\n").map(line => line.trim());

        const getVal = (key: string) => {
            for (const line of lines) {
                const match = line.match(new RegExp(`\\b${key}=(\\S+)`));
                if (match) return match[1];
            }
            return "";
        };

        partitions.push({
            PartitionName: getVal("PartitionName"),
            Default: getVal("Default"),
            MaxTime: getVal("MaxTime"),
            Nodes: getVal("Nodes"),
            State: getVal("State"),
            TotalCPUs: getVal("TotalCPUs"),
            TotalNodes: getVal("TotalNodes"),
            extra: block.trim().replace(/\n/g, '<br/>')
        });
    });

    return partitions;
};

const PartitionsModal: React.FC = () => {
    const [rows, setRows] = useState<PartitionInfo[]>([]);
    const [expandedRows, setExpandedRows] = useState<Set<number>>(new Set());

    useEffect(() => {
        cockpit.spawn(['scontrol', 'show', 'partitions']).then(output => {
            const parsed = parsePartitionsOutput(output);
            setRows(parsed);
        });
    }, []);

    const toggleRow = (index: number) => {
        const newExpanded = new Set(expandedRows);
        // eslint-disable-next-line @typescript-eslint/no-unused-expressions
        newExpanded.has(index) ? newExpanded.delete(index) : newExpanded.add(index);
        setExpandedRows(newExpanded);
    };

    return (
        <Card>
            <CardTitle>{_("SLURM Partitions")}</CardTitle>
            <CardBody>
                <Table variant="compact" aria-label={_("Partitions Table")}>
                    <Thead>
                        <Tr>
                            <Th />
                            <Th>{_("PartitionName")}</Th>
                            <Th>{_("Default")}</Th>
                            <Th>{_("MaxTime")}</Th>
                            <Th>{_("Nodes")}</Th>
                            <Th>{_("State")}</Th>
                            <Th>{_("TotalCPUs")}</Th>
                            <Th>{_("TotalNodes")}</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        {rows.map((row, index) => (
                            <React.Fragment key={index}>
                                <Tr>
                                    <Td
                                        expand={{
                                            rowIndex: index,
                                            isExpanded: expandedRows.has(index),
                                            onToggle: () => toggleRow(index),
                                        }}
                                    />
                                    <Td>{row.PartitionName}</Td>
                                    <Td>{row.Default}</Td>
                                    <Td>{row.MaxTime}</Td>
                                    <Td>{row.Nodes}</Td>
                                    <Td>{row.State}</Td>
                                    <Td>{row.TotalCPUs}</Td>
                                    <Td>{row.TotalNodes}</Td>
                                </Tr>
                                {expandedRows.has(index) && (
                                    <Tr isExpanded>
                                        <Td colSpan={8}>
                                            <ExpandableRowContent>
                                                <div
                                                    dangerouslySetInnerHTML={{ __html: row.extra }}
                                                    style={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap' }}
                                                />
                                            </ExpandableRowContent>
                                        </Td>
                                    </Tr>
                                )}
                            </React.Fragment>
                        ))}
                    </Tbody>
                </Table>
            </CardBody>
        </Card>
    );
};

export default PartitionsModal;
