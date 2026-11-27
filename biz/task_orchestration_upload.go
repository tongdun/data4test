package biz

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CheckTaskOrchestrationUpload never uses the client filename or writes into the
// attachment directory. The raw upload is removed on success and failure alike.
func CheckTaskOrchestrationUpload(source io.Reader, userName string) (map[string]interface{}, error) {
	dir, err := os.MkdirTemp(DownloadBasePath, "orch_upload_")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "archive.zip")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, err
	}
	size, copyErr := io.Copy(f, io.LimitReader(source, archiveMaxBytes+1))
	closeErr := f.Close()
	if copyErr != nil {
		return nil, copyErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if size > archiveMaxBytes {
		return nil, fmt.Errorf(T("error.orchestration_archive_too_large"))
	}
	return ImportTaskOrchestrationCheck(path, userName)
}
