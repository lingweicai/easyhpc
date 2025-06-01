import React, { useEffect, useState } from "react";
import cockpit from 'cockpit';

export function HomeStorageUsage() {
    const [usage, setUsage] = useState<string>("Loading...");

    useEffect(() => {
        cockpit
                .spawn(["df", "--output=pcent", "/home"], { superuser: false })
                .then((output: string) => {
                    const match = output.match(/(\d+)%/);
                    if (match) {
                        setUsage(`${match[1]}%`);
                    } else {
                        setUsage("N/A");
                    }
                })
                .catch(() => {
                    setUsage("Error");
                });
    }, []);

    // eslint-disable-next-line react/jsx-no-useless-fragment
    return <>{ usage }</>;
}

export function NumberJobs(): JSX.Element {
    const [jobCount, setJobCount] = useState<string>('Loading...');

    useEffect(() => {
        cockpit
                .spawn(['sacct', '--noheader', '--format=JobID'])
                .then((output: string) => {
                    const lines = output
                            .trim()
                            .split('\n')
                            .map(line => line.trim())
                            .filter(line => line !== '' && !line.includes('.')); // Exclude sub-jobs

                    setJobCount(lines.length.toString());
                })
                .catch(() => {
                    setJobCount('Error');
                });
    }, []);

    // eslint-disable-next-line react/jsx-no-useless-fragment
    return <>{jobCount}</>;
}
