import React, { useEffect, useState } from "react";
import { Card, CardTitle, CardBody, Progress, Spinner, Icon } from "@patternfly/react-core";
import cockpit from "cockpit";
import { TachometerAltIcon } from "@patternfly/react-icons";

const _ = cockpit.gettext;

interface WorkloadData {
  cpuUsage: number;
  memoryUsage: number;
  storageUsage: number;
}

export default function ClusterWorkloadCard() {
    const [workload, setWorkload] = useState<WorkloadData | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchWorkload = async () => {
            try {
                const [cpuData, memoryData, storageData] = await Promise.all([
                    cockpit.file("/proc/stat").read(),
                    cockpit.file("/proc/meminfo").read(),
                    cockpit.spawn(["df", "--output=pcent", "/"])
                ]);

                const cpuUsage = parseCpuUsage(cpuData);
                const memoryUsage = parseMemoryUsage(memoryData);
                const storageUsage = parseStorageUsage(storageData);

                setWorkload({ cpuUsage, memoryUsage, storageUsage });
            } catch (error) {
                console.error("Failed to fetch workload:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchWorkload();
    }, []);

    const parseCpuUsage = (data: string): number => {
        const lines = data.split("\n");
        const cpuLine = lines.find(line => line.startsWith("cpu "));
        if (!cpuLine) return 0;

        const parts = cpuLine.trim().split(/\s+/)
                .slice(1)
                .map(Number);
        const idle = parts[3] + (parts[4] || 0); // idle + iowait
        const total = parts.reduce((a, b) => a + b, 0);
        return total ? Math.round(((total - idle) / total) * 100) : 0;
    };

    const parseMemoryUsage = (data: string): number => {
        const memTotalLine = data.split("\n").find(line => line.startsWith("MemTotal"));
        const memAvailableLine = data.split("\n").find(line => line.startsWith("MemAvailable"));
        if (!memTotalLine || !memAvailableLine) return 0;

        const total = parseInt(memTotalLine.match(/\d+/)?.[0] || "1", 10);
        const available = parseInt(memAvailableLine.match(/\d+/)?.[0] || "0", 10);

        return total ? Math.round(((total - available) / total) * 100) : 0;
    };

    const parseStorageUsage = (data: string): number => {
        const lines = data.trim().split("\n");
        if (lines.length < 2) return 0;
        const usageStr = lines[1].trim().replace("%", "");
        return parseInt(usageStr, 10) || 0;
    };

    if (loading) {
        return (
            <Card style={{ width: '100%' }}>
                <CardTitle>{_("Head Node Workload")}</CardTitle>
                <CardBody>
                    <Spinner size="lg" />
                </CardBody>
            </Card>
        );
    }

    if (!workload) {
        return (
            <Card className="p-4">
                <CardTitle>{_("Head Node Workload")}</CardTitle>
                <CardBody>
                    <p className="text-red-500">{_("Failed to load data.")}</p>
                </CardBody>
            </Card>
        );
    }

    return (
        <Card style={{ width: '100%' }} isFullHeight>
            <CardTitle>
                <Icon><TachometerAltIcon /></Icon>{'  '}
                {_("Head Node Workload")}
            </CardTitle>
            <CardBody>
                <div style={{ marginTop:"1.8em" }}>
                    <div>{_("CPU Usage")}</div>
                    <Progress value={workload.cpuUsage} />
                </div>

                <div style={{ marginTop:"1.2em" }}>
                    <div>{_("Memory Usage")}</div>
                    <Progress value={workload.memoryUsage} />
                </div>

                <div style={{ marginTop:"1.2em" }}>
                    <div>{_("Storage Usage")}</div>
                    <Progress value={workload.storageUsage} />
                </div>
            </CardBody>
        </Card>
    );
}
