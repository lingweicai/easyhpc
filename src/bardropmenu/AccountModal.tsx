import React, { useEffect, useState } from "react";
import {
    Card,
    CardTitle,
    CardBody,
} from "@patternfly/react-core";
import {
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td,
} from "@patternfly/react-table";
import cockpit from "cockpit";

const _ = cockpit.gettext;

interface UserInfo {
  User: string;
  Account: string;
  DefaultAccount: string;
}

const AccountModal : React.FC = () => {
    const [users, setUsers] = useState<UserInfo[]>([]);

    useEffect(() => {
        cockpit
                .spawn(["bash", "-c", "sacctmgr list user where user=$USER format=User,Account,DefaultAccount"], { superuser: false })
                .then((output) => {
                    const lines = output.trim().split("\n");
                    const dataLines = lines.slice(2); // skip header lines
                    const parsed = dataLines.map((line) => {
                        const match = line.trim().match(/^(\S+)?\s*(\S+)?\s*(\S+)?$/);
                        if (!match) return null;
                        return {
                            User: match[1] || "",
                            Account: match[2] || "",
                            DefaultAccount: match[3] || "",
                        };
                    }).filter(Boolean) as UserInfo[];
                    setUsers(parsed);
                });
    }, []);

    return (
        <Card>
            <CardTitle>{_("User Account Info")}</CardTitle>
            <CardBody>
                <Table aria-label={_("User Info Table")}>
                    <Thead>
                        <Tr>
                            <Th>{_("User")}</Th>
                            <Th>{_("Account")}</Th>
                            <Th>{_("Default Account")}</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        {users.map((user, idx) => (
                            <Tr key={idx}>
                                <Td>{user.User}</Td>
                                <Td>{user.Account}</Td>
                                <Td>{user.DefaultAccount}</Td>
                            </Tr>
                        ))}
                    </Tbody>
                </Table>
            </CardBody>
        </Card>
    );
};

export default AccountModal;
