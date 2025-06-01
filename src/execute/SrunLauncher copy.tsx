import React, { useState, useEffect } from "react";
import {
    Button,
    Modal,
    ModalBody,
    ModalFooter,
    Form,
    FormGroup,
    TextInput,
    Spinner,
    TreeView,
    TreeViewDataItem,
    InputGroup,
    InputGroupItem,
    Tooltip,
} from "@patternfly/react-core";
import { PlayIcon, FolderIcon, FileIcon, FolderOpenIcon } from "@patternfly/react-icons";
import cockpit from "cockpit";
import { useSinfoContext } from '../SinfoContext';

const _ = cockpit.gettext;

interface MyTreeViewDataItem extends TreeViewDataItem {
  isExpanded?: boolean;
  children?: MyTreeViewDataItem[];
}

export default function SrunLauncher() {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [appPath, setAppPath] = useState("");
    const [numNodes, setNumNodes] = useState("1");
    const [numProcs, setNumProcs] = useState("1");
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [isLoadingFiles, setIsLoadingFiles] = useState(false);
    const [fileTree, setFileTree] = useState<MyTreeViewDataItem[]>([]);
    const [isFileTreeOpen, setIsFileTreeOpen] = useState(false);
    const [currentDir, setCurrentDir] = useState<string>("");
    const [srunResults, setSrunResults] = useState("");
    const [isRunning, setIsRunning] = useState(false);

    const { refresh } = useSinfoContext();

    useEffect(() => {
        if (cockpit.info?.user?.home) {
            setCurrentDir(cockpit.info.user.home);
        } else {
            cockpit.user().then(user => {
                setCurrentDir(user.home);
            });
        }
    }, []);

    const fetchDirectory = async (path: string): Promise<MyTreeViewDataItem[]> => {
        try {
            const result = await cockpit.spawn(["ls", "-l", path]);
            const lines = result.trim().split("\n")
                    .slice(1); // Skip total line
            const items: MyTreeViewDataItem[] = lines.map((line) => {
                const parts = line.trim().split(/\s+/);
                const name = parts.slice(8).join(" ");
                const isDir = line[0] === "d";
                const isExecutable = line.includes("x");
                const fullPath = `${path}/${name}`;
                if (!isDir && !isExecutable && !name.endsWith(".sh")) return null;
                return {
                    name,
                    id: fullPath,
                    icon: isDir ? <FolderIcon /> : <FileIcon />,
                    isExpanded: false,
                    children: isDir ? [] : undefined
                };
            }).filter(Boolean) as MyTreeViewDataItem[];
            return items;
        } catch (err) {
            console.error("Error reading directory:", err);
            return [];
        }
    };

    const handleSelect = async (_event: React.MouseEvent, item: MyTreeViewDataItem) => {
        if (!item.id) return;

        if (item.children !== undefined) {
            // it's a folder
            if (item.children.length === 0) {
                const children = await fetchDirectory(item.id);
                item.children = children;
            }
            item.isExpanded = !item.isExpanded;
            setFileTree([...fileTree]); // trigger re-render
        } else {
            // it's a file
            setAppPath(item.id);
            setIsFileTreeOpen(false);
        }
    };

    const handleBrowseExecutables = async () => {
        setIsLoadingFiles(true);
        try {
            const items = await fetchDirectory(currentDir);
            setFileTree(items);
        } finally {
            setIsLoadingFiles(false);
        }
    };

    const handleSubmit = async () => {
        setIsSubmitting(true);
        setSrunResults('');
        const bashCommand = `module load ucx && srun -N ${numNodes} -n ${numProcs} -l ${appPath}`;
        const args = ["bash", "-l", "-c", bashCommand];
        setIsRunning(true);
        try {
            const result = await cockpit.spawn(args, { superuser: false });
            console.log("srun result:", result);
            setSrunResults(result);
            refresh();
            // Don't close the modal here
        } catch (err) {
            console.error("srun error:", err);
            alert(`Failed to launch job: ${err.message || err}`);
            setIsModalOpen(false); // Still close on failure
        } finally {
            setIsSubmitting(false);
            setIsRunning(false); // Also reset this flag
        }
    };

    return (
        <>
            <Button
                variant="primary"
                icon={<PlayIcon />}
                onClick={() => {
                    setIsModalOpen(true);
                    handleBrowseExecutables();
                }}
            >
                {_("srun")}
            </Button>

            <Modal
                title={_("Launch SLURM Job")}
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
            >
                <ModalBody>
                    <Form isHorizontal>
                        <FormGroup label="Application file (.sh or executable)" isRequired fieldId="app-path">
                            <InputGroup>
                                <TextInput
                                id="app-path"
                                value={appPath}
                                onChange={(_evt, val) => setAppPath(val)}
                                isRequired
                                />
                                <InputGroupItem isPlain>
                                    <Tooltip content="Browse Files">
                                        <Button
                                        variant="control"
                                        icon={<FolderOpenIcon />}
                                        onClick={() => setIsFileTreeOpen(!isFileTreeOpen)}
                                        />
                                    </Tooltip>
                                </InputGroupItem>
                            </InputGroup>
                            {isFileTreeOpen && (
                                isLoadingFiles
                                    ? (
                                        <Spinner size="md" />
                                    )
                                    : (
                                        <TreeView
                                            data={fileTree}
                                            hasGuides
                                            onSelect={handleSelect}
                                            activeItems={appPath ? [{ id: appPath, name: appPath }] : []}
                                        />

                                    )
                            )}
                            {appPath && <div style={{ marginTop: "0.5rem" }}>
                                {_("Selected")}: {appPath}
                                        </div>}
                        </FormGroup>

                        <FormGroup label="Number of nodes" isRequired fieldId="num-nodes">
                            <TextInput
                                id="num-nodes"
                                type="number"
                                value={numNodes}
                                onChange={(_, val) => setNumNodes(val)}
                                isRequired
                                placeholder="e.g. 2"
                                min={1}
                            />
                        </FormGroup>
                        <FormGroup label="Number of processes" isRequired fieldId="num-procs">
                            <TextInput
                                id="num-procs"
                                type="number"
                                value={numProcs}
                                onChange={(_, val) => setNumProcs(val)}
                                isRequired
                                placeholder="e.g. 4"
                                min={1}
                            />
                        </FormGroup>
                    </Form>
                    {isRunning && <Spinner size="lg" />}
                    {srunResults &&
                        <div style={{ marginTop: "0.5rem", backgroundColor: "black", color:"white" }}>
                            <pre>{srunResults}</pre>
                        </div>}
                </ModalBody>
                <ModalFooter>
                    <Button
                        key="submit"
                        variant="primary"
                        onClick={handleSubmit}
                        isDisabled={!appPath || !numNodes || !numProcs}
                    >
                        Submit
                    </Button>
                    <Button key="cancel" variant="link" onClick={() => { setIsModalOpen(false); setSrunResults('') }}>
                        Cancel
                    </Button>
                </ModalFooter>
            </Modal>
        </>
    );
}
