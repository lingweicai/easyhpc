import React, { useEffect, useState } from "react";
import {
    Table,
    Thead,
    Tr,
    Th,
    Tbody,
    Td,
    ExpandableRowContent
} from "@patternfly/react-table";
import { Spinner } from "@patternfly/react-core";
import cockpit from "cockpit";

const _ = cockpit.gettext;

interface NodeInfo {
  NodeName: string;
  Arch: string;
  CoresPerSocket: string;
  CPUTot: string;
  OS: string;
  Sockets: string;
  FreeMem: string;
  State: string;
  extra: string;
}

const NodesModal: React.FC = () => {
    const [nodes, setNodes] = useState<NodeInfo[]>([]);
    const [expanded, setExpanded] = useState<Set<number>>(new Set());
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        cockpit
                .spawn(["scontrol", "show", "nodes"])
                .then(output => {
                    const blocks = output.split(/\n\n+/).filter(Boolean);
                    const parsed: NodeInfo[] = blocks.map(block => {
                        const getVal = (key: string) => {
                            const match = block.match(new RegExp(`${key}=(\\S+)`));
                            return match ? match[1] : "";
                        };

                        const osMatch = block.match(/OS=(.+?)(?= NodeAddr|$)/);
                        const OS = osMatch ? osMatch[1] : "";

                        return {
                            NodeName: getVal("NodeName"),
                            Arch: getVal("Arch"),
                            CoresPerSocket: getVal("CoresPerSocket"),
                            CPUTot: getVal("CPUTot"),
                            OS: getVal("OS"),
                            Sockets: getVal("Sockets"),
                            FreeMem: getVal("FreeMem"),
                            State: getVal("State"),
                            extra: block.trim().replace(/\n/g, "<br/>")
                        };
                    });

                    setNodes(parsed);
                    setLoading(false);
                })
                .catch(err => {
                    console.error(_("Failed to fetch nodes:"), err);
                    setLoading(false);
                });
    }, []);

    const toggleExpand = (index: number) => {
        const newExpanded = new Set(expanded);
        // eslint-disable-next-line @typescript-eslint/no-unused-expressions
        newExpanded.has(index) ? newExpanded.delete(index) : newExpanded.add(index);
        setExpanded(newExpanded);
    };

    if (loading) return <Spinner isSVG />;

    return (
        <Table aria-label={_("Nodes Table")}>
            <Thead>
                <Tr>
                    <Th />
                    <Th>{_("NodeName")}</Th>
                    <Th>{_("Arch")}</Th>
                    <Th>{_("CoresPerSocket")}</Th>
                    <Th>{_("CPUTot")}</Th>
                    <Th>{_("OS")}</Th>
                    <Th>{_("Sockets")}</Th>
                    <Th>{_("FreeMem")}</Th>
                    <Th>{_("State")}</Th>
                </Tr>
            </Thead>
            <Tbody>
                {nodes.map((node, index) => (
                    <React.Fragment key={index}>
                        <Tr>
                            <Td
                                expand={{
                                    rowIndex: index,
                                    isExpanded: expanded.has(index),
                                    onToggle: () => toggleExpand(index)
                                }}
                            />
                            <Td>{node.NodeName}</Td>
                            <Td>{node.Arch}</Td>
                            <Td>{node.CoresPerSocket}</Td>
                            <Td>{node.CPUTot}</Td>
                            <Td>{node.OS}</Td>
                            <Td>{node.Sockets}</Td>
                            <Td>{node.FreeMem}</Td>
                            <Td>{node.State}</Td>
                        </Tr>
                        {expanded.has(index) && (
                            <Tr isExpanded>
                                <Td />
                                <Td colSpan={8}>
                                    <ExpandableRowContent>
                                        <div dangerouslySetInnerHTML={{ __html: node.extra }} />
                                    </ExpandableRowContent>
                                </Td>
                            </Tr>
                        )}
                    </React.Fragment>
                ))}
            </Tbody>
        </Table>
    );
};

export default NodesModal;
