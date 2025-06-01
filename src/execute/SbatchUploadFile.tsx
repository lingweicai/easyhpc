import React, { useState } from 'react';
import {
    Form,
    FormGroup,
    FileUpload,
    Button,
    Alert,
    Spinner,
    ModalHeader,
    ModalBody,
} from '@patternfly/react-core';
import cockpit from 'cockpit';

const _ = cockpit.gettext;

interface SbatchUploadFileProps {
  onValueChange: (path: string) => void;
  onCloseModal: () => void;
}

const SbatchUploadFile: React.FC<SbatchUploadFileProps> = ({
    onValueChange,
    onCloseModal,
}) => {
    const [value, setValue] = useState('');
    const [filename, setFilename] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);

    const handleFileInputChange = (_event: unknown, file: File) => {
        setFilename(file.name);
    };

    const handleTextChange = (_event: React.ChangeEvent<HTMLTextAreaElement>, text: string) => {
        setValue(text);
    };

    const handleDataChange = (_event: any, data: string) => {
        setValue(data);
    };

    const handleClear = () => {
        setFilename('');
        setValue('');
        setErrorMsg(null);
    };

    const handleReadStarted = () => {
        setIsLoading(true);
    };

    const handleReadFinished = () => {
        setIsLoading(false);
    };

    const handleUpload = () => {
        if (!filename || !value) {
            setErrorMsg(_("Please select and edit a script before uploading."));
            return;
        }

        setIsLoading(true);
        setErrorMsg(null);

        cockpit.user().then(user => {
            const uploadDir = `${user.home}/slurm_working`;
            const path = `${uploadDir}/${filename}`;

            cockpit.spawn(['mkdir', '-p', uploadDir], { superuser: false })
                    .then(() => {
                        const unixFormatted = value.replace(/\r\n/g, '\n');
                        return cockpit.file(path, { superuser: false }).replace(unixFormatted);
                    })
                    .then(() => {
                        onValueChange(path);
                        onCloseModal();
                    })
                    .catch(err => {
                        setErrorMsg(_("Upload failed: ") + (err.problem || err.message || err));
                    })
                    .finally(() => {
                        setIsLoading(false);
                    });
        })
                .catch(err => {
                    setErrorMsg(_("Failed to get user home directory: ") + (err.problem || err.message || err));
                    setIsLoading(false);
                });
    };

    return (
        <ModalBody>
            <Form isHorizontal>
                <FormGroup label={_("Upload Script File")} fieldId="upload-script-file">
                    <FileUpload
          id="upload-script-file"
          type="text"
          value={value}
          filename={filename}
          filenamePlaceholder={_("Drag and drop or upload a file")}
          onFileInputChange={handleFileInputChange}
          onDataChange={handleDataChange}
          onTextChange={handleTextChange}
          onClearClick={handleClear}
          onReadStarted={handleReadStarted}
          onReadFinished={handleReadFinished}
          isLoading={isLoading}
          allowEditingUploadedText
          browseButtonText={_("Browse")}
                    />
                </FormGroup>

                {errorMsg && (
                    <Alert variant="danger" title={_("Upload Error")} isInline>
                        {errorMsg}
                    </Alert>
                )}

                <Button
        variant="primary"
        onClick={handleUpload}
        isDisabled={!filename || !value || isLoading}
                >
                    {isLoading ? <Spinner size="sm" /> : _("Upload")}
                </Button>
                <Button variant="link" onClick={onCloseModal}>{_("Cancel")}</Button>
            </Form>
        </ModalBody>
    );
};

export default SbatchUploadFile;
