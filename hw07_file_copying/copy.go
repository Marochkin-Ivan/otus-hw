package main

import (
	"errors"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrUnprocessableOffset   = errors.New("offset less than 0")
	ErrUnprocessableLimit    = errors.New("limit less than 0")
	ErrSameFiles             = errors.New("source and dest files are the same")
)

var bufferSize int64 = 10

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrUnprocessableOffset
	}

	if limit < 0 {
		return ErrUnprocessableLimit
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(toPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destFileInfo, err := destFile.Stat()
	if err != nil {
		return err
	}

	if os.SameFile(sourceFileInfo, destFileInfo) {
		return ErrSameFiles
	}

	if offset > sourceFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	copyingSize := sourceFileInfo.Size() - offset
	if limit != 0 && limit < copyingSize {
		copyingSize = limit
	}

	progressBar := progressbar.DefaultBytes(
		copyingSize,
		"copying...",
	)

	_, err = sourceFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	buf := make([]byte, bufferSize)
	for copyingSize > 0 {
		bytesToRead := bufferSize
		if copyingSize < bytesToRead {
			bytesToRead = copyingSize
		}

		read, err := sourceFile.Read(buf[:bytesToRead])
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		if read == 0 {
			break
		}

		_, err = destFile.Write(buf[:read])
		if err != nil {
			return err
		}

		progressBar.Add(read)

		copyingSize -= int64(read)
	}
	return nil
}
