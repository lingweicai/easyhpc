import React, { useState, useEffect } from 'react';
import {
    Card,
    CardHeader,
    CardTitle,
    CardBody,
    Dropdown,
    DropdownItem,
    DropdownList,
    MenuToggle,
    Divider
} from '@patternfly/react-core';
import { Chart, ChartAxis, ChartLine, ChartGroup } from '@patternfly/react-charts/victory';

const _ = (window as any).cockpit.gettext;

// Define a type for the trend data points
type TrendPoint = { x: string; y: number };

const ReportModal: React.FC = () => {
    const [trendData, setTrendData] = useState<TrendPoint[]>([]);
    const [range, setRange] = useState<'day' | 'week' | 'month'>('day');
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);

    const dateLabels = (range: string): string[] => {
        const today = new Date();
        const labels: string[] = [];
        for (let i = 6; i >= 0; i--) {
            const d = new Date(today);
            if (range === 'day') d.setDate(today.getDate() - i);
            else if (range === 'week') d.setDate(today.getDate() - i * 7);
            else if (range === 'month') d.setMonth(today.getMonth() - i);
            labels.push(
                range === 'month'
                    ? `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
                    : `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
            );
        }
        return labels;
    };

    useEffect(() => {
        const fetchData = async () => {
            const labels = dateLabels(range);
            const data: TrendPoint[] = [];
            for (const label of labels) {
                const start = range === 'month' ? `${label}-01` : label;
                const cmd = [
                    'bash',
                    '-c',
                    `sreport user TopUsage start=${start} end=now format=Login,Used`
                ];

                try {
                    const output = await (window as any).cockpit.spawn(cmd, { superuser: false });
                    const match = output.match(/\bUsed\s+\n[-\s]+\n.*\s+(\d+(?:\.\d+)?)\b/);
                    const used = match ? parseFloat(match[1]) : 0;
                    data.push({ x: label, y: used });
                } catch (error) {
                    console.error(_('Failed to fetch sreport data'), error);
                    data.push({ x: label, y: 0 });
                }
            }
            setTrendData(data);
        };
        fetchData();
    }, [range]);

    const dropdownItems = (
        <>
            <DropdownItem key="day" onClick={() => { setRange('day'); setIsDropdownOpen(false) }}>{_('Daily')}</DropdownItem>
            <DropdownItem key="week" onClick={() => { setRange('week'); setIsDropdownOpen(false) }}>{_('Weekly')}</DropdownItem>
            <DropdownItem key="month" onClick={() => { setRange('month'); setIsDropdownOpen(false) }}>{_('Monthly')}</DropdownItem>
            <Divider component="li" key="separator" />
            <DropdownItem isDisabled key="year">{_('Yearly (Coming Soon)')}</DropdownItem>
        </>
    );

    const headerActions = (
        <Dropdown
            onSelect={() => setIsDropdownOpen(false)}
            isOpen={isDropdownOpen}
            onOpenChange={(isOpen: boolean) => setIsDropdownOpen(isOpen)}
            toggle={(toggleRef) => (
                <MenuToggle
                    ref={toggleRef}
                    isExpanded={isDropdownOpen}
                    onClick={() => setIsDropdownOpen((prev) => !prev)}
                    variant='secondary'
                    aria-label={_('Trend range toggle')}
                >
                    {range === 'day' ? _('Daily') : range === 'week' ? _('Weekly') : _('Monthly')}
                </MenuToggle>
            )}
        >
            <DropdownList>{dropdownItems}</DropdownList>
        </Dropdown>
    );

    return (
        <Card>
            <CardHeader actions={{ actions: headerActions }}>
                <CardTitle>{_('User Usage Trend')}</CardTitle>
            </CardHeader>
            <CardBody>
                <div style={{ height: '250px', width: '100%' }}>
                    <Chart
                        ariaDesc={_('User usage trend chart')}
                        ariaTitle={_('User usage trend')}
                        height={250}
                        width={600}
                        padding={{ top: 20, bottom: 60, left: 60, right: 20 }}
                        domainPadding={{ x: [30, 25] }}
                    >
                        <ChartAxis tickFormat={(t) => t} fixLabelOverlap />
                        <ChartAxis dependentAxis tickFormat={(y) => `${y}`} />
                        <ChartGroup>
                            <ChartLine data={trendData} />
                        </ChartGroup>
                    </Chart>
                </div>
            </CardBody>
        </Card>
    );
};

export default ReportModal;
