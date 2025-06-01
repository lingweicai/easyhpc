import React, { useEffect, useState } from "react";
import {
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td
} from "@patternfly/react-table";
import {
    Spinner,
    Bullseye,
    Label,
    LabelProps,
    Button,
    ButtonVariant,
    Tooltip,
    Modal,
    ModalFooter,
    ModalBody,
} from "@patternfly/react-core";
import { PauseIcon, PlayIcon, TrashIcon } from "@patternfly/react-icons";
import cockpit from "cockpit";
import { useSinfoContext } from "../SinfoContext"; // Adjust path as needed

const _ = cockpit.gettext;

interface SqueueRow {
  jobid: string;
  partition: string;
  name: string;
  user: string;
  state: string;
  time: string;
  nodes: string;
  nodelist: string;
}

const getJobStateColor = (state: string): NonNullable<LabelProps["color"]> => {
    switch (state.toUpperCase()) {
    case "R":
        return "green"; // Running
    case "PD":
        return "orange"; // Pending
    case "CG":
        return "blue"; // Completing
    case "CD":
        return "grey"; // Completed
    case "F":
        return "red"; // Failed
    default:
        return "purple"; // Fallback color
    }
};

const SqueueTable: React.FC = () => {
    const [rows, setRows] = useState<SqueueRow[]>([]);
    const [loading, setLoading] = useState(true);
    const [currentUser, setCurrentUser] = useState<string>("");
    const [isAdmin, setIsAdmin] = useState<boolean>(false);

    const [confirmOpen, setConfirmOpen] = useState(false);
    const [confirmAction, setConfirmAction] = useState<"suspend" | "resume" | "cancel" | null>(null);
    const [confirmJobid, setConfirmJobid] = useState<string | null>(null);

    const { refresh: refreshSinfo } = useSinfoContext();
    const fetchData = () => {
        setLoading(true);
        cockpit.spawn(["squeue"])
                .then((output: string) => {
                    const lines = output.trim().split("\n");
                    const parsed: SqueueRow[] = [];

                    lines.forEach(line => {
                        const match = line.match(/^ *(\d+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\d+)\s+(.+)$/);
                        if (match) {
                            const [, jobid, partition, name, user, state, time, nodes, nodelist] = match;
                            parsed.push({ jobid, partition, name, user, state, time, nodes, nodelist });
                        }
                    });

                    setRows(parsed);
                })
                .catch((err) => console.error("Failed to run squeue:", err))
                .finally(() => setLoading(false));
    };

    const openConfirmDialog = (action: "suspend" | "resume" | "cancel", jobid: string) => {
        setConfirmAction(action);
        setConfirmJobid(jobid);
        setConfirmOpen(true);
    };

    const handleConfirm = () => {
        if (!confirmAction || !confirmJobid) return;
        const command = confirmAction === "cancel"
            ? `scancel ${confirmJobid}`
            : `sudo scontrol ${confirmAction} ${confirmJobid}`;
        cockpit.spawn(["bash", "-c", command])
                .then(() => {
                    console.log(`${confirmAction} succeeded for job ${confirmJobid}`);
                    fetchData();
                    refreshSinfo();
                })
                .catch((err) => {
                    console.error(`${confirmAction} failed for job ${confirmJobid}:`, err);
                })
                .finally(() => {
                    setConfirmOpen(false);
                    setConfirmAction(null);
                    setConfirmJobid(null);
                });
    };

    const { refreshCounter } = useSinfoContext();
    useEffect(() => {
        // Get current user
        cockpit.spawn(["whoami"])
                .then((output: string) => setCurrentUser(output.trim()))
                .catch((err) => console.error("Failed to get current user:", err));

        // Check admin groups
        cockpit.spawn(["id", "-Gn"])
                .then((groups: string) => {
                    const groupList = groups.trim().split(/\s+/);
                    setIsAdmin(groupList.includes("slurm") || groupList.includes("wheel") || groupList.includes("root"));
                })
                .catch((err) => console.error("Failed to get groups:", err));

        fetchData();
    }, [refreshCounter]);

    if (loading) {
        return (
            <Bullseye>
                <Spinner />
            </Bullseye>
        );
    }

    return (
        <>
            <Table aria-label="squeue table" variant="compact">
                <Thead>
                    <Tr>
                        <Th>{_("JOBID")}</Th>
                        <Th>{_("PARTITION")}</Th>
                        <Th>{_("NAME")}</Th>
                        <Th>{_("USER")}</Th>
                        <Th>{_("ST")}</Th>
                        <Th>{_("TIME")}</Th>
                        <Th>{_("NODES")}</Th>
                        <Th>{_("NODELIST")}</Th>
                        <Th>{_("ACTIONS")}</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {rows.map((row, idx) => {
                        const jobState = row.state.toUpperCase();
                        const canControl = row.user === currentUser || isAdmin;

                        return (
                            <Tr key={idx}>
                                <Td dataLabel="JOBID">{row.jobid}</Td>
                                <Td dataLabel="PARTITION">{row.partition}</Td>
                                <Td dataLabel="NAME">{row.name}</Td>
                                <Td dataLabel="USER">{row.user}</Td>
                                <Td dataLabel="ST">
                                    <Label color={getJobStateColor(jobState)}>{jobState}</Label>
                                </Td>
                                <Td dataLabel="TIME">{row.time}</Td>
                                <Td dataLabel="NODES">{row.nodes}</Td>
                                <Td dataLabel="NODELIST(REASON)">{row.nodelist}</Td>
                                <Td dataLabel="Actions">
                                    {canControl && (
                                        <>
                                            <Tooltip content="Suspend job">
                                                <Button
                                                    icon={<PauseIcon />}
                                                    variant={ButtonVariant.plain}
                                                    size="sm"
                                                    onClick={() => openConfirmDialog("suspend", row.jobid)}
                                                    isDisabled={jobState !== "R"}
                                                />
                                            </Tooltip>
                                            <Tooltip content="Resume job">
                                                <Button
        icon={<PlayIcon />}
        variant={ButtonVariant.plain}
        size="sm"
        onClick={() => openConfirmDialog("resume", row.jobid)}
        isDisabled={jobState !== "S"} // Enable only if Suspended
                                                />
                                            </Tooltip>

                                            <Tooltip content="Cancel job">
                                                <Button
                                                    icon={<TrashIcon />}
                                                    variant={ButtonVariant.plain}
                                                    size="sm"
                                                    onClick={() => openConfirmDialog("cancel", row.jobid)}
                                                />
                                            </Tooltip>
                                        </>
                                    )}
                                </Td>
                            </Tr>
                        );
                    })}
                </Tbody>
            </Table>

            <Modal
                variant="small"
                title="Confirm Job Action"
                isOpen={confirmOpen}
                onClose={() => setConfirmOpen(false)}
            >
                <ModalBody>Are you sure you want to {confirmAction} job {confirmJobid}?</ModalBody>
                <ModalFooter>
                    <Button variant="danger" onClick={handleConfirm}>
                        Confirm
                    </Button>
                    <Button variant="link" onClick={() => setConfirmOpen(false)}>
                        Cancel
                    </Button>
                </ModalFooter>
            </Modal>

        </>
    );
};

export default SqueueTable;
