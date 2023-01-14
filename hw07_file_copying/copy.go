package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWrongOffsetParam      = errors.New("wrong offset")
	ErrWrongLimitParam       = errors.New("wrong limit")
	ErrWrongFromParam        = errors.New("empty from file path")
	ErrWrongToParam          = errors.New("empty to file path")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrWrongFromParam
	}
	if toPath == "" {
		return ErrWrongToParam
	}
	if offset < 0 {
		return ErrWrongOffsetParam
	}
	if limit < 0 {
		return ErrWrongLimitParam
	}

	var err error
	var fileFrom *os.File
	if fileFrom, err = os.Open(fromPath); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	defer fileFrom.Close()

	var fi fs.FileInfo
	if fi, err = fileFrom.Stat(); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fromFileSize := fi.Size()
	remainSize := fromFileSize - offset
	if remainSize <= 0 {
		return ErrOffsetExceedsFileSize
	}

	needSize := offset + limit
	if limit == 0 || needSize > fromFileSize {
		limit = remainSize
	}

	if _, err = fileFrom.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	var fileTo *os.File
	if fileTo, err = os.Create(toPath); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	defer fileTo.Close()

	bar := pb.Full.Start64(limit)
	defer bar.Finish()

	bar.Set(pb.Bytes, true)

	barReader := bar.NewProxyReader(fileFrom)
	if _, err = io.CopyN(fileTo, barReader, limit); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}
