import React from 'react';
import { useSinfoContext } from '../../SinfoContext';
import { Spinner, Alert } from '@patternfly/react-core';

export const NodeCount: React.FC = () => {
    const { rawOutput, loading, error } = useSinfoContext();

    if (loading) return <Spinner />;
    if (error) return <Alert variant="danger" title="Error fetching sinfo">{error}</Alert>;

    const totalNodes = calculateTotalNodes(rawOutput);

    return (
        <strong>{ totalNodes }</strong>
    );
};

const calculateTotalNodes = (raw: string): number => {
    const lines = raw.trim().split('\n');
    if (lines.length < 2) return 0;

    // Skip header row
    return lines.slice(1).reduce((sum, line) => {
        const columns = line.trim().split(/\s+/);
        const nodeCount = parseInt(columns[3], 10); // 4th column is index 3
        return sum + (isNaN(nodeCount) ? 0 : nodeCount);
    }, 0);
};
