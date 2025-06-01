import { ChartDonut, ChartThemeColor } from '@patternfly/react-charts/victory';
import { Card, CardTitle, CardBody, Icon } from '@patternfly/react-core';
import { ServerIcon } from '@patternfly/react-icons';
import cockpit from 'cockpit';
import React, { useEffect, useState } from 'react';

const _ = cockpit.gettext;

interface NodeStateData {
  [state: string]: number;
}

function parseNodeStates(output: string): NodeStateData {
    const stateCounts: NodeStateData = {};
    const lines = output.split('\n');
    for (const line of lines) {
        const stateMatch = line.match(/State=([A-Z_]+)/);
        if (stateMatch) {
            const state = stateMatch[1];
            stateCounts[state] = (stateCounts[state] || 0) + 1;
        }
    }
    return stateCounts;
}

const NodeDonutChart: React.FC = () => {
    const [nodeStates, setNodeStates] = useState<NodeStateData>({});

    useEffect(() => {
        cockpit.spawn(['scontrol', 'show', 'nodes'], { superuser: 'try' })
                .then((output: string) => {
                    const parsed = parseNodeStates(output);
                    setNodeStates(parsed);
                })
                .catch((error: any) => {
                    console.error('Failed to fetch node info:', error);
                });
    }, []);

    const chartData = Object.entries(nodeStates).map(([state, count]) => ({
        x: state,
        y: count,
    }));

    return (
        <Card style={{ width: '100%' }} isFullHeight>
            <CardTitle>
                <Icon><ServerIcon /></Icon>{'  '}
                {_("Nodes")}
            </CardTitle>
            <CardBody>
                <ChartDonut
          ariaDesc={_("Node state distribution")}
          ariaTitle={_("Slurm node state donut chart")}
          constrainToVisibleArea
          data={chartData}
          height={250}
          width={300}
          labels={({ datum }) => `${_(datum.x)}: ${datum.y}`}
          legendData={chartData.map(d => ({ name: `${_(d.x)}: ${d.y}` }))}
          legendOrientation="horizontal"
          legendPosition="bottom"
          name="node-states-chart"
          padding={{ top: 20, bottom: 40, left: 10, right: 10 }}
          subTitle={_("Total Nodes")}
          title={Object.values(nodeStates).reduce((sum, count) => sum + count, 0)
                  .toString()}
          themeColor={ChartThemeColor.multiOrdered}
                />
            </CardBody>
        </Card>
    );
};

export default NodeDonutChart;
