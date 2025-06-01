import React, { useEffect, useState } from 'react';
import {
    Chart,
    ChartAxis,
    ChartGroup,
    ChartBar,
    ChartScatter,
    ChartTooltip
} from '@patternfly/react-charts/victory';
import cockpit from 'cockpit';
import { Card, CardBody, CardTitle, Icon } from '@patternfly/react-core';
import { ChartPieIcon } from "@patternfly/react-icons";

// Add gettext
const _ = cockpit.gettext;

interface PartitionData {
  name: string;
  isDefault: boolean;
  totalNodes: number;
  totalCPUs: number;
  state: string;
  nodes: string;
  billing: number;
}

const parsePartitionOutput = (output: string): PartitionData[] => {
    const blocks = output.split(/(?=PartitionName=)/).filter(Boolean);
    return blocks.map(block => {
        const get = (pattern: RegExp, fallback: any = null) => {
            const match = block.match(pattern);
            return match ? match[1] : fallback;
        };

        const name = get(/PartitionName=(\S+)/, 'unknown');
        const isDefault = get(/Default=(\S+)/, 'NO') === 'YES';
        const totalNodes = parseInt(get(/TotalNodes=(\d+)/, '0'));
        const totalCPUs = parseInt(get(/TotalCPUs=(\d+)/, '0'));
        const state = get(/State=(\S+)/, 'UNKNOWN');
        const nodes = get(/\bNodes=(\S+)/, 'N/A');
        const billing = parseInt(get(/billing=(\d+)/, '0'));

        return { name, isDefault, totalNodes, totalCPUs, state, nodes, billing };
    });
};

const getStateColor = (state: string) => {
    switch (state) {
    case 'UP': return '#4cb140'; // green
    case 'DOWN': return '#8a8d90'; // grey
    default: return '#4cb140'; // green
    }
};

const PartitionsChart: React.FC = () => {
    const [data, setData] = useState<PartitionData[]>([]);

    useEffect(() => {
        cockpit.spawn(['scontrol', 'show', 'partitions']).then((output: string) => {
            setData(parsePartitionOutput(output));
        });
    }, []);

    return (
        <Card style={{ width: '100%' }}>
            <CardTitle>
                <Icon>
                    <ChartPieIcon />
                </Icon>{'  '}
                {/* Translate title */}
                {_('Partitions')}
            </CardTitle>
            <CardBody style={{ paddingLeft: '5px', paddingRight: '5px' }}>
                <Chart
                    domainPadding={{ x: 30 }}
                    height={250}
                    width={300}
                    padding={{ top: 20, bottom: 30, left: 42, right: 42 }}
                >
                    {/* X Axis */}
                    <ChartAxis />

                    {/* Left Y Axis for TotalNodes (scaled ×2, labels ÷2) */}
                    <ChartAxis
                        dependentAxis
                        tickFormat={(t) => t / 2}
                        label={_('Total Nodes')}
                        style={{
                            axisLabel: { padding: 24 }
                        }}
                    />

                    {/* Right Y Axis for TotalCPUs */}
                    <ChartAxis
                        dependentAxis
                        orientation="right"
                        label={_('Total CPUs')}
                        style={{
                            axisLabel: { padding: 24 }
                        }}
                    />

                    <ChartGroup offset={20}>
                        {/* Bar chart for TotalNodes */}
                        <ChartBar
                            data={data.map(d => ({
                                x: d.name + (d.isDefault ? ` ${_('(default)')}` : ''),
                                y: d.totalNodes * 2,
                                label: `${_('Total Nodes')}: ${d.totalNodes}\n${_('Nodes')}: ${d.nodes}\n${_('State')}: ${d.state}`,
                                fill: getStateColor(d.state),
                            }))}
                            labels={({ datum }) => datum.label}
                            labelComponent={<ChartTooltip constrainToVisibleArea />}
                        />

                        {/* Scatter chart for TotalCPUs */}
                        <ChartScatter
                            data={data.map(d => ({
                                x: d.name + (d.isDefault ? ` ${_('(default)')}` : ''),
                                y: d.totalCPUs,
                                label: `${_('Total CPUs')}: ${d.totalCPUs}`,
                                symbol: 'star',
                            }))}
                            size={8}
                            style={{
                                data: {
                                    fill: ({ datum }) => getStateColor(
                                        data.find(
                                            p => p.name + (p.isDefault ? ` ${_('(default)')}` : '') === datum.x
                                        )?.state || 'UNKNOWN'
                                    ),
                                },
                            }}
                            labels={({ datum }) => datum.label}
                            labelComponent={<ChartTooltip constrainToVisibleArea />}
                        />
                    </ChartGroup>
                </Chart>
            </CardBody>
        </Card>
    );
};

export default PartitionsChart;
