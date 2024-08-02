package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	inputPathTest = "./testdata/input.txt"
	outPathTest   = "./testdata/out.txt"

	outOffset0Limit0       = "./testdata/out_offset0_limit0.txt"
	outOffset0Limit10      = "./testdata/out_offset0_limit10.txt"
	outOffset0Limit1000    = "./testdata/out_offset0_limit1000.txt"
	outOffset0Limit10000   = "./testdata/out_offset0_limit10000.txt"
	outOffset100Limit1000  = "./testdata/out_offset100_limit1000.txt"
	outOffset6000Limit1000 = "./testdata/out_offset6000_limit1000.txt"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name        string
		offset      int64
		limit       int64
		compareFile string
	}{
		{
			name:        "Offset = 0, Limit = 0",
			offset:      0,
			limit:       0,
			compareFile: outOffset0Limit0,
		},
		{
			name:        "Offset = 0, Limit = 10",
			offset:      0,
			limit:       10,
			compareFile: outOffset0Limit10,
		},
		{
			name:        "Offset = 0, Limit = 1000",
			offset:      0,
			limit:       1000,
			compareFile: outOffset0Limit1000,
		},
		{
			name:        "Offset = 0, Limit = 10000",
			offset:      0,
			limit:       10000,
			compareFile: outOffset0Limit10000,
		},
		{
			name:        "Offset = 100, Limit = 1000",
			offset:      100,
			limit:       1000,
			compareFile: outOffset100Limit1000,
		},
		{
			name:        "Offset = 6000, Limit = 1000",
			offset:      6000,
			limit:       1000,
			compareFile: outOffset6000Limit1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(inputPathTest, outPathTest, tt.offset, tt.limit)
			if err != nil {
				require.NoError(t, err, "Copy() error: %v")
				return
			}

			fileInfo, err := os.Stat(tt.compareFile)
			if err != nil {
				require.NoError(t, err, "Error getting file information: %v")
			}

			dstFileInfo, err := os.Stat(outPathTest)
			if err != nil {
				require.NoError(t, err, "Error getting file information: %v")
			}
			if dstFileInfo.Size() != fileInfo.Size() {
				require.Equalf(t, fileInfo.Size(), dstFileInfo.Size(),
					"The copied file size does not match the expected size: %d != %d", dstFileInfo.Size(), fileInfo.Size())
			}

			os.Remove(outPathTest)
		})
	}
}

func TestCopyN(t *testing.T) {
	tempFile, err := os.CreateTemp("", "copyN_test")
	if err != nil {
		require.NoError(t, err)
	}
	defer os.Remove(tempFile.Name())

	srcFile, err := os.Open(inputPathTest)
	if err != nil {
		require.NoError(t, err)
	}
	defer srcFile.Close()

	copiedBytes, err := copyN(tempFile, srcFile, 10)
	if err != nil {
		require.NoError(t, err)
	}

	if copiedBytes != 10 {
		require.Equal(t, 10, copiedBytes, "Copied %d bytes, expected 10", copiedBytes)
	}

	tempFile.Seek(0, io.SeekStart)
	data, err := io.ReadAll(tempFile)
	if err != nil {
		require.NoError(t, err)
	}

	expectedData := []byte("Go\nDocumen")
	if !bytes.Equal(data, expectedData) {
		require.Equal(t, expectedData, data, "Invalid data in temporary file: %s", data)
	}
}
