import React, { useState } from 'react';
import {
    Button,
    Modal,
    ModalVariant,
    Alert,
    ModalHeader,
    ModalBody,
    ModalFooter,
    Spinner,
    Flex,
    FlexItem
} from '@patternfly/react-core';
import cockpit from 'cockpit';
import SbatchSelectFiles from './SbatchSelectFiles';
import SbatchUploadFile from './SbatchUploadFile';
import { useSinfoContext } from '../SinfoContext';
import { FileCodeIcon } from '@patternfly/react-icons';

const _ = cockpit.gettext;

const SbatchLauncher: React.FC = () => {
    const [isLauncherModalOpen, setLauncherModalOpen] = useState(false);
    const [isFileSelectModalOpen, setFileSelectModalOpen] = useState(false);
    const [isUploadModalOpen, setUploadModalOpen] = useState(false);
    const [selectedScript, setSelectedScript] = useState<string | null>(null);
    const [submitting, setSubmitting] = useState(false);
    const [submissionOutput, setSubmissionOutput] = useState<string | null>(null);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);

    const { refresh } = useSinfoContext();

    const handleSubmit = () => {
        if (!selectedScript) {
            setErrorMsg(_("Please select a script file first."));
            return;
        }

        setSubmitting(true);
        setSubmissionOutput(null);
        setErrorMsg(null);

        cockpit.spawn(['/usr/bin/sbatch', selectedScript])
                .then(output => {
                    setSubmissionOutput(output);
                    refresh();
                })
                .catch(error => {
                    setErrorMsg(`${_("Submission failed")}: ${error}`);
                })
                .finally(() => {
                    setSubmitting(false);
                });
    };

    return (
        <>
            <Button onClick={() => setLauncherModalOpen(true)} variant="primary" icon={<FileCodeIcon />}>
                {_("sbatch")}
            </Button>

            <Modal
                variant={ModalVariant.medium}
                title={_("Launch SLURM Job")}
                isOpen={isLauncherModalOpen}
                onClose={() => setLauncherModalOpen(false)}
            >
                <ModalHeader>{_("Select and Submit Job Script")}</ModalHeader>
                <ModalBody>
                    {errorMsg && <Alert variant="danger" title={errorMsg} />}
                    {submissionOutput && <Alert variant="success" title={_("Submission Output")} isInline>{submissionOutput}</Alert>}

                    <div style={{ marginBottom: '1rem' }}>
                        <strong>{_("Selected Script")}:</strong>{' '}
                        {selectedScript ? <code>{selectedScript}</code> : _("None selected")}
                    </div>

                    <Flex spaceItems={{ default: 'spaceItemsMd' }}>
                        <FlexItem>
                            <Button variant="secondary" onClick={() => setFileSelectModalOpen(true)}>
                                {_("Choose or Edit Script File")}
                            </Button>
                        </FlexItem>
                        <FlexItem>
                            <Button variant="secondary" onClick={() => setUploadModalOpen(true)}>
                                {_("Upload Script")}
                            </Button>
                        </FlexItem>
                    </Flex>
                </ModalBody>
                <ModalFooter>
                    <Button
                        variant="primary"
                        onClick={handleSubmit}
                        isDisabled={!selectedScript || submitting}
                    >
                        {submitting ? <Spinner size="sm" /> : _("Submit")}
                    </Button>
                    <Button variant="link" onClick={() => setLauncherModalOpen(false)}>{_("Close")}</Button>
                </ModalFooter>
            </Modal>

            <Modal
                variant={ModalVariant.large}
                isOpen={isFileSelectModalOpen}
                onClose={() => setFileSelectModalOpen(false)}
                title={_("Select a Script File")}
            >
                <SbatchSelectFiles
                    onValueChange={(value) => {
                        setSelectedScript(value);
                        setFileSelectModalOpen(false);
                    }}
                    onCloseModal={() => setFileSelectModalOpen(false)}
                />
            </Modal>

            <Modal
                variant={ModalVariant.medium}
                isOpen={isUploadModalOpen}
                onClose={() => setUploadModalOpen(false)}
                title={_("Upload Slurm Script")}
            >
                <SbatchUploadFile
                    onValueChange={(value: React.SetStateAction<string | null>) => {
                        setSelectedScript(value);
                        setUploadModalOpen(false);
                    }}
                    onCloseModal={() => setUploadModalOpen(false)}
                />
            </Modal>
        </>
    );
};

export default SbatchLauncher;
