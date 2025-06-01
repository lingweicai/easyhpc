import React, { useState } from 'react';
import {
    Button,
    Form,
    FormGroup,
    Modal,
    ModalBody,
    ModalFooter,
    ModalHeader,
    ModalVariant,
    TextInput,
    Alert
} from '@patternfly/react-core';
import cockpit from 'cockpit';
import { useSinfoContext } from '../SinfoContext';
import PdshCommandModal from './PdshCommandModal';
import { TerminalIcon } from '@patternfly/react-icons';

const _ = cockpit.gettext;

const SallocLauncher: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [nodes, setNodes] = useState('1');
    const [partition, setPartition] = useState('');
    const [output, setOutput] = useState('');
    const [jobId, setJobId] = useState('');
    const [state, setState] = useState('');
    const [batchHost, setBatchHost] = useState('');
    const [nodeList, setNodeList] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);

    const { refresh } = useSinfoContext();

    const toggleModal = () => {
        setIsModalOpen(!isModalOpen);
        setOutput('');
        setJobId('');
        setState('');
        setNodeList('');
        setBatchHost('');
    };

    const parseOutput = (text: string) => {
        const jobMatch = text.match(/Job ID:\s*(\d+)/);
        const stateMatch = text.match(/State:\s*(\w+)/);
        const nodeListMatch = text.match(/Node List:\s*(.*)/);
        const batchHostMatch = text.match(/BatchHost:\s*(.*)/);

        setJobId(jobMatch?.[1] || '');
        const jobState = stateMatch?.[1] || '';
        setState(jobState);
        setNodeList(nodeListMatch?.[1] || '');
        setBatchHost(batchHostMatch?.[1] || '');

        if (jobState === 'RUNNING') {
            refresh();
        }
    };

    const handleSubmit = () => {
        setIsSubmitting(true);
        const args = ['/home/dev/easyhpc/src/salloc.sh'];
        if (partition) args.push('-p', partition);
        if (nodes) args.push('-N', nodes);

        cockpit.spawn(args, { superuser: false })
                .then((result: string) => {
                    setOutput(result);
                    parseOutput(result);
                })
                .catch((err: any) => {
                    setOutput(typeof err === 'string' ? err : JSON.stringify(err));
                    setState('ERROR');
                })
                .finally(() => setIsSubmitting(false));
    };

    const launchTerminal = () => {
        if (batchHost) {
            const url = '/cockpit/@localhost/system/terminal.html';
            const newTab = window.open(url, '_blank');
            setTimeout(() => {
                newTab?.alert(_("Please run: ssh") + ` ${batchHost}`);
            }, 1000);
        } else {
            window.open('/cockpit/@localhost/system/terminal.html', '_blank');
        }
    };

    return (
        <>
            <Button variant="primary" onClick={toggleModal} icon={<TerminalIcon />}>
                {_("salloc")}
            </Button>
            <Modal
                variant={ModalVariant.large}
                isOpen={isModalOpen}
                onClose={toggleModal}
                aria-labelledby="salloc-modal-title"
            >
                <ModalHeader title={_("Submit salloc.sh")} labelId="salloc-modal-title" />
                <ModalBody>
                    <Form isHorizontal>
                        <FormGroup label={_("Partition (optional)")} fieldId="partition">
                            <TextInput
                                id="partition"
                                value={partition}
                                onChange={(_e, val) => setPartition(val)}
                            />
                        </FormGroup>
                        <FormGroup label={_("Number of Nodes")} fieldId="nodes">
                            <TextInput
                                type="number"
                                id="nodes"
                                value={nodes}
                                onChange={(_e, val) => setNodes(val)}
                            />
                        </FormGroup>
                    </Form>

                    {output && (
                        <div className="mt-sm">
                            <Alert variant="info" isInline title={_("salloc.sh Output")}>
                                <pre style={{ whiteSpace: 'pre-wrap' }}>{output}</pre>
                                <div className="mt-sm">
                                    <strong>{_("Job ID")}:</strong> {jobId}<br />
                                    <strong>{_("State")}:</strong> {state}<br />
                                    <strong>{_("BatchHost")}:</strong> {batchHost}<br />
                                    <strong>{_("Node List")}:</strong> {nodeList}
                                </div>
                            </Alert>
                        </div>
                    )}

                    {state === 'RUNNING' && (
                        <>
                            <PdshCommandModal nodeList={nodeList} />
                            <Button variant="secondary" className="mt-sm" onClick={launchTerminal}>
                                {_("Open Terminal")}
                            </Button>
                        </>
                    )}

                    {state === 'CANCELLED' && (
                        <Alert
                            variant="warning"
                            isInline
                            title={_("Job was cancelled. No resources allocated.")}
                            className="pf-v5-u-mt-md"
                        />
                    )}

                    {state === 'ERROR' && (
                        <Alert
                            variant="danger"
                            isInline
                            title={_("Error running salloc.sh")}
                            className="pf-v5-u-mt-md"
                        />
                    )}
                </ModalBody>
                <ModalFooter>
                    <Button
                        variant="primary"
                        onClick={handleSubmit}
                        isDisabled={isSubmitting}
                    >
                        {isSubmitting ? _("Submitting...") : _("Submit")}
                    </Button>
                    <Button variant="link" onClick={toggleModal}>
                        {_("Cancel")}
                    </Button>
                </ModalFooter>
            </Modal>
        </>
    );
};

export default SallocLauncher;
