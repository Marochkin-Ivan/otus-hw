package main

import (
	"bytes"
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

	if limit == 0 || limit > int64(buffer.Len()) {
		bar := progressbar.DefaultBytes(
			int64(buffer.Len()),
			"copying...",
		)

		_, err = io.Copy(io.MultiWriter(destFile, bar), buffer)
		if err != nil {
			return err
		}
	} else {
		bar := progressbar.DefaultBytes(
			limit,
			"copying...",
		)
		_, err = io.CopyN(io.MultiWriter(destFile, bar), buffer, limit)
		if err != nil {
			return err
		}
	}

	return nil
}
