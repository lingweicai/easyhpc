import React, { useState } from "react";
import { Card, CardTitle, CardBody, Button, Tooltip } from "@patternfly/react-core";
import { SyncAltIcon } from "@patternfly/react-icons";
import SqueueTable from "./SqueueTable";
import cockpit from 'cockpit';
const _ = cockpit.gettext;

export default function SqueueCard() {
    const [refreshKey, setRefreshKey] = useState(0);

    const handleRefresh = () => {
        setRefreshKey(prev => prev + 1);
    };

    return (
        <Card style={{ marginTop: '8px' }}>
            <CardTitle>{_("Squeue")}
                <Tooltip content="Refresh">
                    <Button variant="plain" onClick={handleRefresh} aria-label="Refresh" style={{ backgroundColor: "transparent" }}>
                        <SyncAltIcon />
                    </Button>
                </Tooltip>
            </CardTitle>
            <CardBody>
                <SqueueTable key={refreshKey} />
            </CardBody>
        </Card>
    );
}
