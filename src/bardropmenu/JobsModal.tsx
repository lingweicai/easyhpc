import React, { useEffect, useState } from "react";
import {
    Card,
    CardTitle,
    CardBody,
    Spinner
} from "@patternfly/react-core";
import {
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td
} from "@patternfly/react-table";
import cockpit from "cockpit";

const _ = cockpit.gettext;

interface JobData {
  JobID: string;
  JobName: string;
  Partition: string;
  Account: string;
  AllocCPUS: string;
  State: string;
  ExitCode: string;
}

const JobsModal: React.FC = () => {
    const [jobs, setJobs] = useState<JobData[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        setLoading(true);
        cockpit
                .spawn(["sacct", "--noheader", "-o", "JobID,JobName,Partition,Account,AllocCPUS,State,ExitCode"])
                .then(output => {
                    const lines = output
                            .split("\n")
                            .map(line => line.trim())
                            .filter(line => line.length > 0);

                    const parsedJobs = lines.map(line => {
                        const parts = line.split(/\s+/);
                        return {
                            JobID: parts[0] || "",
                            JobName: parts[1] || "",
                            Partition: parts[2] || "",
                            Account: parts[3] || "",
                            AllocCPUS: parts[4] || "",
                            State: parts[5] || "",
                            ExitCode: parts[6] || ""
                        };
                    });

                    setJobs(parsedJobs);
                    setLoading(false);
                })
                .catch(err => {
                    console.error(_("Failed to fetch jobs:"), err);
                    setLoading(false);
                });
    }, []);

    return (
        <Card>
            <CardTitle>{_("SLURM Jobs")}</CardTitle>
            <CardBody>
                {loading
                    ? (
                        <Spinner isSVG />
                    )
                    : (
                        <Table aria-label={_("SLURM Jobs Table")} variant="compact">
                            <Thead>
                                <Tr>
                                    <Th>{_("JobID")}</Th>
                                    <Th>{_("JobName")}</Th>
                                    <Th>{_("Partition")}</Th>
                                    <Th>{_("Account")}</Th>
                                    <Th>{_("AllocCPUS")}</Th>
                                    <Th>{_("State")}</Th>
                                    <Th>{_("ExitCode")}</Th>
                                </Tr>
                            </Thead>
                            <Tbody>
                                {jobs.map((job, index) => (
                                    <Tr key={index}>
                                        <Td>{job.JobID}</Td>
                                        <Td>{job.JobName}</Td>
                                        <Td>{job.Partition}</Td>
                                        <Td>{job.Account}</Td>
                                        <Td>{job.AllocCPUS}</Td>
                                        <Td>{job.State}</Td>
                                        <Td>{job.ExitCode}</Td>
                                    </Tr>
                                ))}
                            </Tbody>
                        </Table>
                    )}
            </CardBody>
        </Card>
    );
};

export default JobsModal;
