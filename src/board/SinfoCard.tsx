import React, { useState } from "react";
import { Card, CardTitle, CardBody, Button, Tooltip } from "@patternfly/react-core";
import { SyncAltIcon } from "@patternfly/react-icons";
import SinfoTable from "./SinfoTable";
import cockpit from 'cockpit';
const _ = cockpit.gettext;

export default function SinfoCard() {
    const [refreshKey, setRefreshKey] = useState(0);

    const handleRefresh = () => {
        setRefreshKey(prev => prev + 1);
    };

    return (
        <Card style={{ marginTop: '8px' }}>
            <CardTitle>{_("Sinfo")}
                <Tooltip content="Refresh">
                    <Button variant="plain" onClick={handleRefresh} aria-label="Refresh" style={{ backgroundColor: "transparent" }}>
                        <SyncAltIcon />
                    </Button>
                </Tooltip>
            </CardTitle>
            <CardBody>
                <SinfoTable key={refreshKey} />
            </CardBody>
        </Card>
    );
}
