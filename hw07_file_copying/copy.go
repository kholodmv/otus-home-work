package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenSourceFile        = errors.New("failed to open source file")
	ErrFileInfo              = errors.New("failed to get file information")
	ErrSeekFile              = errors.New("could not seek in source file")
	ErrCreateFile            = errors.New("could not create destination file")
	ErrCopyN                 = errors.New("failed to copy file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		return ErrOpenSourceFile
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return ErrFileInfo
	}
	fileSize := fileInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return ErrSeekFile
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return ErrCreateFile
	}
	defer dstFile.Close()

	if limit <= 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	bar := pb.New(int(limit))
	bar.Start()

	n, err := copyN(dstFile, file, limit)
	if err != nil {
		return ErrCopyN
	}

	bar.Finish()

	fmt.Printf("Copied %d bytes\n", n)

	return nil
}

func copyN(dst io.Writer, src io.Reader, n int64) (int64, error) {
	copiedBytes := int64(0)
	buf := make([]byte, n)
	for copiedBytes < n {
		readBytes, err := src.Read(buf)
		if err != nil {
			return copiedBytes, err
		}

		if readBytes == 0 {
			break
		}

		writtenBytes, err := dst.Write(buf[:readBytes])
		if err != nil {
			return copiedBytes, err
		}

		copiedBytes += int64(writtenBytes)
	}
	return copiedBytes, nil
}
