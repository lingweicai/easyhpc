import React, { useEffect, useState } from "react";
import {
    Card,
    CardHeader,
    CardTitle,
    CardBody,
    Dropdown,
    DropdownItem,
    DropdownList,
    MenuToggle,
    Icon
} from "@patternfly/react-core";
import {
    Chart,
    ChartAxis,
    ChartGroup,
    ChartLine,
    ChartThemeColor,
    ChartVoronoiContainer,
} from "@patternfly/react-charts/victory";
import cockpit from "cockpit";
import { HistoryIcon } from "@patternfly/react-icons";

const _ = cockpit.gettext;

const ReportLineChartCard: React.FC = () => {
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [selectedMode, setSelectedMode] = useState("daily");
    const [chartData, setChartData] = useState<{ x: string; y: number }[]>([]);

    const modes = ["daily", "weekly", "monthly"];

    const fetchData = async (mode: string) => {
        try {
            const result = await cockpit.spawn(
                ["/home/dev/easyhpc/src/sreport.sh", mode],
                { superuser: false }
            );
            const parsed = result
                .trim()
                .split("\n")
                .map((line) => {
                    const [dateStr, valueStr] = line.trim().split(/\s+/);
                    const y = parseInt(valueStr, 10);
                    let x = dateStr;
                    if (mode === "daily" || mode === "weekly") {
                        const [yyyy, mm, dd] = dateStr.split("-");
                        x = `${mm}/${dd}`;
                    } else if (mode === "monthly") {
                        const [yyyy, mm] = dateStr.split("-");
                        x = `${yyyy}/${mm}`;
                    }
                    return { x, y };
                });
            setChartData(parsed);
        } catch (error) {
            console.error(_("Failed to fetch chart data:"), error);
        }
    };

    useEffect(() => {
        fetchData(selectedMode);
    }, [selectedMode]);

    const onSelect = (
        _event?: React.MouseEvent<Element>,
        value?: string | number
    ) => {
        if (typeof value === "string") {
            setSelectedMode(value);
            setIsDropdownOpen(false);
        }
    };

    return (
        <Card style={{ width: '100%' }} isFullHeight>
            <CardHeader
                actions={{
                    actions: (
                        <Dropdown
                            isOpen={isDropdownOpen}
                            onSelect={onSelect}
                            onOpenChange={setIsDropdownOpen}
                            toggle={(toggleRef) => (
                                <MenuToggle
                                    ref={toggleRef}
                                    onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                                    isExpanded={isDropdownOpen}
                                    style={{ fontSize: '0.8em' }}
                                >
                                    {_(selectedMode.charAt(0).toUpperCase() + selectedMode.slice(1))}
                                </MenuToggle>
                            )}
                        >
                            <DropdownList style={{ fontSize: '0.8em' }}>
                                {modes.map((mode) => (
                                    <DropdownItem key={mode} value={mode}>
                                        {_(mode)}
                                    </DropdownItem>
                                ))}
                            </DropdownList>
                        </Dropdown>
                    ),
                }}
            >
                <CardTitle>
                    <Icon>
                        <HistoryIcon />
                    </Icon>{'  '}
                    {_("Utilization Report")}
                </CardTitle>
            </CardHeader>
            <CardBody style={{ paddingLeft: '4px' }}>
                <Chart
                    ariaDesc={_("Utilization over time")}
                    ariaTitle={_("Utilization Line Chart")}
                    containerComponent={
                        <ChartVoronoiContainer
                            labels={({ datum }) => `${datum.x}: ${datum.y}`}
                        />
                    }
                    height={250}
                    width={450}
                    padding={{ top: 20, bottom: 50, left: 70, right: 20 }}
                    themeColor={ChartThemeColor.blue}
                    domainPadding={{ x: 10 }}
                >
                    <ChartAxis label={_("Date")} fixLabelOverlap />
                    <ChartAxis
                        dependentAxis
                        label={_("Usage")}
                        style={{
                            axisLabel: { padding: 52 }
                        }}
                    />
                    <ChartGroup>
                        <ChartLine data={chartData} />
                    </ChartGroup>
                </Chart>
            </CardBody>
        </Card>
    );
};

export default ReportLineChartCard;
