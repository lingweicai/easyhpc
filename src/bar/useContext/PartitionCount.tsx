// PartitionCount.tsx
import React from 'react';
import { useSinfoContext } from '../../SinfoContext';
import { Spinner, Alert } from '@patternfly/react-core';

export const PartitionCount: React.FC = () => {
    const { rawOutput, loading, error } = useSinfoContext();

    if (loading) return <Spinner />;
    if (error) return <Alert variant="danger" title="Error fetching sinfo">{error}</Alert>;

    const count = countPartitions(rawOutput);

    return (
        <strong>{count}</strong>
    );
};

const countPartitions = (raw: string): number => {
    const lines = raw.trim().split('\n');
    if (lines.length < 2) return 0; // no data
    return lines.length - 1; // exclude header line
};
