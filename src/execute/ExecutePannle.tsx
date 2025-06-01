// ExecutePanel.tsx
import React from "react";
import { Card, CardBody, CardTitle, Flex } from "@patternfly/react-core";
import SrunLauncher from "./SrunLauncher";
import SbatchLauncher from "./SbatchLauncher";
import SallocLauncher from "./SallocLauncher";
import cockpit from "cockpit";

const _ = cockpit.gettext;

const ExecutePanel: React.FC = () => {
    return (
        <Card style={{ marginTop: '8px' }}>
            <CardTitle>{_("Running")}</CardTitle>
            <CardBody>
                <Flex>
                    <SallocLauncher />
                    <SrunLauncher />
                    <SbatchLauncher />
                </Flex>
            </CardBody>
        </Card>
    );
};

export default ExecutePanel;
