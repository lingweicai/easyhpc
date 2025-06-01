import React, { useState } from 'react';
import {
    Card,
    CardHeader,
    CardTitle,
    CardBody,
    CardExpandableContent,
    Level,
    LabelGroup,
    Label,
    Grid,
    Flex,
    Dropdown,
    DropdownList,
    DropdownItem,
    MenuToggle,
    MenuToggleElement,
    LevelItem,
    FlexItem,
    Icon,
    Modal,
    Button,
    ModalFooter,
    ModalBody
} from '@patternfly/react-core';
import EllipsisVIcon from '@patternfly/react-icons/dist/esm/icons/ellipsis-v-icon';
import {
    ChartPieIcon,
    ServerIcon,
    TachometerAltIcon,
    HistoryIcon,
    PlayIcon,
    TerminalIcon,
    FileCodeIcon
} from "@patternfly/react-icons";

import PartitionsChart from './bar/PartitionsChart';
import NodeDonutChart from './bar/NodeDonutChart';
import ClusterWorkloadCard from './bar/ClusterWorkloadCard';
import ReportLineChartCard from './bar/ReportLineChartCard';
import { PartitionCount } from './bar/useContext/PartitionCount';
import { NodeCount } from './bar/useContext/NodeCount';
import { HomeStorageUsage, NumberJobs } from './bar/BarHeaderInfo';
import cockpit from 'cockpit';

import PartitionsModal from './bardropmenu/PartitionsModal';
import NodesModal from './bardropmenu/NodesModal';
import JobsModal from './bardropmenu/JobsModal';
import ReportModal from './bardropmenu/ReportModal';
import AccountModal from './bardropmenu/AccountModal';

const _ = cockpit.gettext;

const Dashbar: React.FunctionComponent = () => {
    const [isCardExpanded, setIsCardExpanded] = useState(false);
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [openModal, setOpenModal] = useState<string | null>(null);

    const modalContentMap: { [key: string]: React.ReactNode } = {
        Partitions: <PartitionsModal />,
        Nodes: <NodesModal />,
        Jobs: <JobsModal />,
        Report: <ReportModal />,
        Account: <AccountModal />
    };

    const onCardExpand = () => {
        setIsCardExpanded(!isCardExpanded);
    };

    const onActionToggle = () => {
        setIsDropdownOpen(!isDropdownOpen);
    };

    const onActionSelect = (_event?: any, itemId?: string | number) => {
        setIsDropdownOpen(false);
        if (typeof itemId === 'string') {
            setOpenModal(itemId); // Open the corresponding modal
        }
    };

    const dropdownItems = (
        <>
            <DropdownItem key="action1" itemId="Partitions">{_("Partitions")}</DropdownItem>
            <DropdownItem key="action2" itemId="Nodes">{_("Nodes")}</DropdownItem>
            <DropdownItem key="action3" itemId="Jobs">{_("Jobs")}</DropdownItem>
            <DropdownItem key="action4" itemId="Report">{_("Report")}</DropdownItem>
            <DropdownItem key="action5" itemId="Account">{_("Account")}</DropdownItem>
        </>
    );

    const headerActions = (
        <Dropdown
            onSelect={onActionSelect}
            isOpen={isDropdownOpen}
            popperProps={{ position: 'right' }}
            onOpenChange={(isOpen: boolean) => setIsDropdownOpen(isOpen)}
            toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                <MenuToggle
                    ref={toggleRef}
                    isExpanded={isDropdownOpen}
                    onClick={onActionToggle}
                    variant="plain"
                    aria-label="Horizontal card grid demo kebab toggle"
                    icon={<EllipsisVIcon />}
                />
            )}
        >
            <DropdownList>{dropdownItems}</DropdownList>
        </Dropdown>
    );

    return (
        <>
            <Card id="dashbar card" isExpanded={isCardExpanded}>
                <CardHeader
                    actions={{ actions: headerActions }}
                    onExpand={onCardExpand}
                    toggleButtonProps={{
                        id: 'toggle-button',
                        'aria-label': 'Actions',
                        'aria-labelledby': 'titleId toggle-button',
                        'aria-expanded': isCardExpanded
                    }}
                >
                    {isCardExpanded && <CardTitle id="titleId">{_("Charts")}</CardTitle>}
                    {!isCardExpanded && (
                        <Level hasGutter style={{ justifyContent: 'space-between', width: '100%' }}>
                            <LevelItem>
                                <Flex alignItems={{ default: 'alignItemsCenter' }}>
                                    <FlexItem>
                                        <CardTitle id="titleId">{_("Summary")}</CardTitle>
                                    </FlexItem>
                                    <FlexItem>
                                        <LabelGroup numLabels={4} isCompact>
                                            <Label isCompact icon={<ChartPieIcon />} color="blue" style={{ marginRight: '1rem' }}>
                                                {_("Partitions")}: <span className='dashbar-header-summary'><PartitionCount /></span>
                                            </Label>
                                            <Label isCompact icon={<ServerIcon />} color="purple" style={{ marginRight: '1rem' }}>
                                                {_("Nodes")}: <span className='dashbar-header-summary'><NodeCount /></span>
                                            </Label>
                                            <Label isCompact icon={<TachometerAltIcon />} color="green" style={{ marginRight: '1rem' }}>
                                                {_("Home Storage")}: <span className='dashbar-header-summary'><HomeStorageUsage /></span>
                                            </Label>
                                            <Label isCompact icon={<HistoryIcon />} color="orange">
                                                {_("Today Jobs")}: <span className='dashbar-header-summary'><NumberJobs /></span>
                                            </Label>
                                        </LabelGroup>
                                    </FlexItem>
                                </Flex>
                            </LevelItem>
                            <LevelItem>
                                <Icon style={{ marginRight: "1rem" }}><TerminalIcon /></Icon>
                                <Icon style={{ marginRight: "1rem" }}><PlayIcon /></Icon>
                                <Icon style={{ marginRight: "1rem" }}><FileCodeIcon /></Icon>
                            </LevelItem>
                        </Level>
                    )}
                </CardHeader>
                <CardExpandableContent>
                    <CardBody>
                        <Grid md={6} lg={3} hasGutter>
                            <Flex style={{ width: '100%' }} spaceItems={{ default: 'spaceItemsLg' }} alignItems={{ default: 'alignItemsFlexStart' }} direction={{ default: 'column' }}>
                                <Flex style={{ width: '100%' }} spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsFlexStart' }} direction={{ default: 'column' }} grow={{ default: 'grow' }}>
                                    <PartitionsChart />
                                </Flex>
                            </Flex>
                            <Flex style={{ width: '100%' }} spaceItems={{ default: 'spaceItemsLg' }} alignItems={{ default: 'alignItemsFlexStart' }} direction={{ default: 'column' }}>
                                <Flex style={{ width: '100%' }} spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsFlexStart' }} direction={{ default: 'column' }} grow={{ default: 'grow' }}>
                                    <NodeDonutChart />
                                </Flex>
                            </Flex>
                            <Flex style={{ width: '100%' }} spaceItems={{ default: 'spaceItemsLg' }} alignItems={{ default: 'alignItemsFlexStart' }} direction={{ default: 'column' }}>
                                <Flex style={{ width: '100%' }} spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsFlexStart' }} direction={{ default: 'column' }} grow={{ default: 'grow' }}>
                                    <ClusterWorkloadCard />
                                </Flex>
                            </Flex>
                            <Flex style={{ width: '100%' }} spaceItems={{ default: 'spaceItemsLg' }} alignItems={{ default: 'alignItemsFlexStart' }} direction={{ default: 'column' }}>
                                <Flex style={{ width: '100%' }} spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsFlexStart' }} direction={{ default: 'column' }} grow={{ default: 'grow' }}>
                                    <ReportLineChartCard />
                                </Flex>
                            </Flex>
                        </Grid>
                    </CardBody>
                </CardExpandableContent>
            </Card>

            {/* Modals */}

            {Object.keys(modalContentMap).map((key) => (
                <Modal
        key={key}
        title={_(key)}
        isOpen={openModal === key}
        onClose={() => setOpenModal(null)}
                >
                    <ModalBody>
                        {modalContentMap[key]}
                    </ModalBody>
                    <ModalFooter>
                        <Button key="close" variant="primary" onClick={() => setOpenModal(null)}>
                            {_("Close")}
                        </Button>
                    </ModalFooter>
                </Modal>
            ))}
        </>
    );
};

export default Dashbar;
