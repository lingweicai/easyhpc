import React, { useState } from 'react';
import {
    Modal,
    ModalFooter,
    Button,
    Form,
    FormGroup,
    TextArea,
    Select,
    SelectList,
    SelectOption,
    MenuToggle,
    Badge,
    Spinner,
    Alert,
    ModalBody,
} from '@patternfly/react-core';
import cockpit from 'cockpit';

const _ = cockpit.gettext;

function parseNodeList(nodeListStr: string): string[] {
    const match = nodeListStr.match(/^([a-zA-Z\-]+)\[(.+)]$/);
    if (!match) return [nodeListStr];
    const prefix = match[1];
    const inner = match[2];
    const result: string[] = [];

    for (const part of inner.split(',')) {
        if (part.includes('-')) {
            const [start, end] = part.split('-').map(Number);
            for (let i = start; i <= end; i++) {
                result.push(`${prefix}${i}`);
            }
        } else {
            result.push(`${prefix}${part}`);
        }
    }

    return result;
}

interface PdshCommandModalProps {
    nodeList: string;
}

const PdshCommandModal: React.FC<PdshCommandModalProps> = ({ nodeList }) => {
    const initialNodes = nodeList ? parseNodeList(nodeList) : [];
    const [selectedNodes, setSelectedNodes] = useState<string[]>(initialNodes);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isSelectOpen, setIsSelectOpen] = useState(false);
    const [command, setCommand] = useState('');
    const [output, setOutput] = useState('');
    const [isRunning, setIsRunning] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const toggleModal = () => setIsModalOpen(!isModalOpen);

    const handleNodeSelect = (_e: unknown, value: string | number | undefined) => {
        const val = String(value);
        setSelectedNodes((prev) =>
            prev.includes(val) ? prev.filter((v) => v !== val) : [...prev, val]
        );
    };

    const handleSubmit = async () => {
        if (selectedNodes.length === 0 || command.trim() === '') return;

        setIsRunning(true);
        setOutput('');
        setError(null);

        const nodeList = selectedNodes.join(',');
        const cmd = ['pdsh', '-w', nodeList, command];

        try {
            const result = await cockpit.spawn(cmd, { err: 'out' });
            setOutput(result);
        } catch (err) {
            setError(err.toString());
        } finally {
            setIsRunning(false);
        }
    };

    return (
        <>
            <Button onClick={toggleModal}>{_("Run Pdsh Command")}</Button>

            <Modal
                title={_("Run Pdsh Command")}
                isOpen={isModalOpen}
                onClose={toggleModal}
            >
                <ModalBody>
                    <Form isHorizontal>
                        <FormGroup label={_("Select Nodes")} isRequired fieldId="nodes">
                            <Select
                                isOpen={isSelectOpen}
                                selected={selectedNodes}
                                onSelect={handleNodeSelect}
                                onOpenChange={setIsSelectOpen}
                                toggle={(ref) => (
                                    <MenuToggle
                                        ref={ref}
                                        onClick={() => setIsSelectOpen(!isSelectOpen)}
                                        isExpanded={isSelectOpen}
                                    >
                                        {_("Choose nodes")}{' '}
                                        {selectedNodes.length > 0 && (
                                            <Badge isRead>{selectedNodes.length}</Badge>
                                        )}
                                    </MenuToggle>
                                )}
                            >
                                <SelectList role="menu">
                                    {initialNodes.map((node) => (
                                        <SelectOption
                                            key={node}
                                            id={node}
                                            value={node}
                                            isSelected={selectedNodes.includes(node)}
                                            hasCheckbox
                                        >
                                            {node}
                                        </SelectOption>
                                    ))}
                                </SelectList>
                            </Select>
                        </FormGroup>

                        <FormGroup label={_("Command")} isRequired fieldId="command">
                            <TextArea
                                id="command"
                                name="command"
                                value={command}
                                onChange={(_event, value) => setCommand(value)}
                                isRequired
                                resizeOrientation="vertical"
                            />
                        </FormGroup>
                    </Form>

                    {isRunning && <Spinner size="lg" />}
                    {error && (
                        <Alert variant="danger" title={_("Command failed")} isInline>
                            {error}
                        </Alert>
                    )}
                    {output && (
                        <div
                            style={{
                                marginTop: '1rem',
                                whiteSpace: 'pre-wrap',
                                backgroundColor: 'black',
                                padding: '1rem',
                                borderRadius: '6px',
                                color: 'white'
                            }}
                        >
                            <div>{output}</div>
                        </div>
                    )}
                </ModalBody>
                <ModalFooter>
                    <Button
                        variant="primary"
                        onClick={handleSubmit}
                        isDisabled={isRunning}
                    >
                        {isRunning ? <Spinner size="sm" /> : _("Run")}
                    </Button>
                    <Button variant="link" onClick={toggleModal}>
                        {_("Cancel")}
                    </Button>
                </ModalFooter>
            </Modal>
        </>
    );
};

export default PdshCommandModal;
