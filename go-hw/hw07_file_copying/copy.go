package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fFrom.Close()
	stat, err := fFrom.Stat()
	if err != nil {
		return err
	}
	if !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	fTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fTo.Close()
	if limit == 0 {
		limit = stat.Size() - offset
	}
	// start new bar
	bar := pb.Full.Start64(limit)
	// create proxy reader
	barReader := bar.NewProxyReader(fFrom)
	// finish bar
	defer bar.Finish()

	_, err = io.CopyN(fTo, barReader, limit)
	if err != nil {
		return err
	}

	return nil
}
