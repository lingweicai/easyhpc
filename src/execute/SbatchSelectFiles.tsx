import React, { useEffect, useState } from 'react';
import cockpit from 'cockpit';
import {
    Button,
    Modal,
    ModalVariant,
    TextArea,
    TextInput,
    Form,
    FormGroup,
    Spinner,
    Divider,
    Alert,
    ModalBody,
    ModalFooter,
    ModalHeader
} from '@patternfly/react-core';
import {
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td
} from '@patternfly/react-table';

const _ = cockpit.gettext;

type SbatchSelectFilesProps = {
  onValueChange: (value: string) => void;
  onCloseModal: () => void;
};

const SbatchSelectFiles: React.FC<SbatchSelectFilesProps> = ({ onValueChange, onCloseModal }) => {
    const [files, setFiles] = useState<{ name: string; size: string; date: string }[]>([]);
    const [directory, setDirectory] = useState('');
    const [selectedFile, setSelectedFile] = useState<string | null>(null);
    const [fileContent, setFileContent] = useState('');
    const [editedContent, setEditedContent] = useState('');
    const [newFilename, setNewFilename] = useState('');
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [loading, setLoading] = useState(false);
    const [errorMsg, setErrorMsg] = useState('');

    useEffect(() => {
        cockpit.spawn(['sh', '-c', 'echo $HOME'])
                .then(output => {
                    const dir = `${output.trim()}/slurm_working`;
                    setDirectory(dir);
                    return cockpit.spawn(['ls', '-l', dir]);
                })
                .then(parseFileList)
                .catch(err => setErrorMsg(_("Error accessing slurm_working: ") + err));
    }, []);

    const parseFileList = (output: string) => {
        const lines = output.split('\n').slice(1);
        const fileList = lines.map(line => {
            const parts = line.trim().split(/\s+/);
            if (parts.length < 9) return null;
            const name = parts.slice(8).join(' ');
            const size = parts[4];
            const date = `${parts[5]} ${parts[6]} ${parts[7]}`;
            return { name, size, date };
        }).filter(
            (f): f is { name: string; size: string; date: string } =>
                f !== null && /\.(sh|script)$/.test(f.name)
        );
        setFiles(fileList);
    };

    const openFile = (filename: string) => {
        const fullPath = `${directory}/${filename}`;
        setLoading(true);
        cockpit.file(fullPath)
                .read()
                .then(content => {
                    setSelectedFile(fullPath);
                    setFileContent(content);
                    setEditedContent(content);
                    setIsModalOpen(true);
                    setIsEditing(false);
                    setNewFilename('');
                })
                .catch(err => setErrorMsg(_("Error reading file: ") + err))
                .finally(() => setLoading(false));
    };

    const saveFile = async (filename: string, content: string) => {
        const fullPath = `${directory}/${filename}`;
        try {
            await cockpit.file(fullPath, { superuser: false }).replace(content);
            setFiles(prev => {
                if (prev.find(f => f.name === filename)) return prev;
                return [...prev, {
                    name: filename,
                    size: `${content.length}`,
                    date: new Date().toDateString()
                }];
            });
            setSelectedFile(fullPath);
            setIsModalOpen(false);
            onValueChange(fullPath);
        } catch (err) {
            setErrorMsg(_("Error saving file: ") + err);
        }
    };

    const handleConfirmSelection = () => {
        if (selectedFile) {
            onValueChange(selectedFile);
            onCloseModal(); // Close parent modal
        }
    };

    return (
        <>
            <ModalHeader>{_("Select a script file")}</ModalHeader>
            <ModalBody>
                {errorMsg && <Alert variant="danger" title={errorMsg} />}
                <Table variant="compact" aria-label={_("Script Files Table")}>
                    <Thead>
                        <Tr>
                            <Th>{_("Select")}</Th>
                            <Th>{_("File Name")}</Th>
                            <Th>{_("Size")}</Th>
                            <Th>{_("Date")}</Th>
                            <Th />
                        </Tr>
                    </Thead>
                    <Tbody>
                        {files.map(file => (
                            <Tr key={file.name}>
                                <Td>
                                    <input
                                        type="radio"
                                        name="fileSelect"
                                        checked={selectedFile === `${directory}/${file.name}`}
                                        onChange={() => setSelectedFile(`${directory}/${file.name}`)}
                                    />
                                </Td>
                                <Td>{file.name}</Td>
                                <Td>{file.size}</Td>
                                <Td>{file.date}</Td>
                                <Td>
                                    <Button variant="link" onClick={() => openFile(file.name)}>
                                        {_("View")}
                                    </Button>
                                </Td>
                            </Tr>
                        ))}
                    </Tbody>
                </Table>
                <Divider />
                {selectedFile && (
                    <div style={{ marginTop: '1rem' }}>
                        <strong>{_("Selected:")}</strong> {selectedFile}
                    </div>
                )}

                <Modal
                    id="filecontent"
                    isOpen={isModalOpen}
                    variant={ModalVariant.large}
                    title={selectedFile?.split('/').pop() || _("View File")}
                    onClose={() => setIsModalOpen(false)}
                >
                    <ModalBody>
                        {loading
                            ? <Spinner isSVG />
                            : (
                                <Form isHorizontal>
                                    {isEditing && (
                                        <FormGroup label={_("Save As")} fieldId="new-filename">
                                            <TextInput
                                                id="new-filename"
                                                value={newFilename}
                                                onChange={(_, val) => setNewFilename(val)}
                                                placeholder={_("Leave blank to overwrite")}
                                            />
                                        </FormGroup>
                                    )}
                                    <FormGroup fieldId="file-content">
                                        <TextArea
                                            value={isEditing ? editedContent : fileContent}
                                            onChange={(_, val) => setEditedContent(val)}
                                            readOnly={!isEditing}
                                            rows={20}
                                        />
                                    </FormGroup>
                                </Form>
                            )}
                    </ModalBody>
                    <ModalFooter>
                        {isEditing
                            ? (
                                <>
                                    <Button
                                        variant="primary"
                                        onClick={() => saveFile(newFilename || selectedFile!.split('/').pop()!, editedContent)}
                                        isDisabled={loading}
                                    >
                                        {_("Save")}
                                    </Button>
                                    <Button variant="link" onClick={() => setIsModalOpen(false)}>
                                        {_("Cancel")}
                                    </Button>
                                </>
                            )
                            : (
                                <>
                                    <Button variant="secondary" onClick={() => setIsEditing(true)}>
                                        {_("Edit")}
                                    </Button>
                                    <Button variant="link" onClick={() => setIsModalOpen(false)}>
                                        {_("Close")}
                                    </Button>
                                </>
                            )}
                    </ModalFooter>
                </Modal>
            </ModalBody>
            <ModalFooter>
                <Button
                    variant="primary"
                    onClick={handleConfirmSelection}
                    isDisabled={!selectedFile}
                >
                    {_("Confirm")}
                </Button>
                <Button variant="link" onClick={onCloseModal}>
                    {_("Cancel")}
                </Button>
            </ModalFooter>
        </>
    );
};

export default SbatchSelectFiles;
