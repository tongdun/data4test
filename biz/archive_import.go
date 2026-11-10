package biz

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const archiveMaxFiles = 2000
const archiveMaxBytes int64 = 500 << 20

// archiveTarget rejects traversal and platform-dependent separators before any write.
func archiveTarget(root, name string) (string, error) {
	if strings.Contains(name, "\\") || filepath.IsAbs(name) {
		return "", fmt.Errorf("unsafe archive path: %s", name)
	}
	for _, part := range strings.Split(name, "/") {
		if part == ".." {
			return "", fmt.Errorf("unsafe archive path: %s", name)
		}
	}
	target := filepath.Join(root, name)
	if target != filepath.Clean(root) && !strings.HasPrefix(target, filepath.Clean(root)+string(os.PathSeparator)) {
		return "", fmt.Errorf("unsafe archive path: %s", name)
	}
	return target, nil
}

func writeArchiveEntry(root, name string, r io.Reader, remaining *int64) error {
	target, err := archiveTarget(root, name)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	n, copyErr := io.Copy(f, io.LimitReader(r, *remaining+1))
	closeErr := f.Close()
	*remaining -= n
	if *remaining < 0 {
		return fmt.Errorf("archive exceeds uncompressed size limit")
	}
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func extractBoundedZip(data []byte, root string, maxFiles int, maxBytes int64) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	if len(zr.File) > maxFiles {
		return fmt.Errorf("archive exceeds file count limit")
	}
	remaining := maxBytes
	for _, entry := range zr.File {
		target, err := archiveTarget(root, entry.Name)
		if err != nil {
			return err
		}
		if entry.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("archive links are not allowed")
		}
		if entry.FileInfo().IsDir() {
			if err = os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}
		if entry.UncompressedSize64 > uint64(remaining) {
			return fmt.Errorf("archive exceeds uncompressed size limit")
		}
		if strings.HasPrefix(filepath.Base(entry.Name), ".") {
			remaining -= int64(entry.UncompressedSize64)
			continue
		}
		r, err := entry.Open()
		if err != nil {
			return err
		}
		err = writeArchiveEntry(root, entry.Name, r, &remaining)
		r.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func extractBoundedTar(r io.Reader, root string, maxFiles int, maxBytes int64) error {
	tr := tar.NewReader(r)
	remaining := maxBytes
	for count := 0; ; count++ {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if count >= maxFiles {
			return fmt.Errorf("archive exceeds file count limit")
		}
		target, err := archiveTarget(root, header.Name)
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if header.Size > remaining {
				return fmt.Errorf("archive exceeds uncompressed size limit")
			}
			if strings.HasPrefix(filepath.Base(header.Name), ".") {
				remaining -= header.Size
				continue
			}
			if err = writeArchiveEntry(root, header.Name, tr, &remaining); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported archive entry: %s", header.Name)
		}
	}
}
