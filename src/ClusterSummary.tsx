import React, { useEffect, useState } from "react";
import { Card, CardBody, Title } from "@patternfly/react-core";
import cockpit from "cockpit";

const ClusterSummary = () => {
  const [nodes, setNodes] = useState(0);
  const [jobs, setJobs] = useState(0);
  const [cpuUsage, setCpuUsage] = useState({ allocated: 0, idle: 0, total: 0 });

  useEffect(() => {
    cockpit.script("sinfo -h -o '%D'").then(output => setNodes(parseInt(output)));
    cockpit.script("squeue -h -t RUNNING | wc -l").then(output => setJobs(parseInt(output)));
    cockpit.script("sinfo -h -o '%C'").then(output => {
      const [allocated, idle, , total] = output.trim().split("/").map(Number);
      setCpuUsage({ allocated, idle, total });
    });
  }, []);

  return (
    <div className="grid grid-cols-3 gap-4">
      <Card>
        <CardBody>
          <Title headingLevel="h2">Total Nodes</Title>
          <p>{nodes}</p>
        </CardBody>
      </Card>
      <Card>
        <CardBody>
          <Title headingLevel="h2">Running Jobs</Title>
          <p>{jobs}</p>
        </CardBody>
      </Card>
      <Card>
        <CardBody>
          <Title headingLevel="h2">CPU Usage</Title>
          <p>{cpuUsage.allocated} / {cpuUsage.total} Allocated</p>
        </CardBody>
      </Card>
    </div>
  );
};

export default ClusterSummary;
