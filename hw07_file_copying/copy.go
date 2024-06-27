package main

import (
	"bytes"
	"errors"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	info, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	if offset > info.Size() {
		return ErrOffsetExceedsFileSize
	}

	buf := make([]byte, info.Size()-offset)
	_, err = sourceFile.ReadAt(buf, offset)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	destFile, err := os.OpenFile(toPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buffer := bytes.NewBuffer(buf)

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
