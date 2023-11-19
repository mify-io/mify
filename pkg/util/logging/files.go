package logging

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	LOG_FILE_MAX_SIZE = int64(1024 * 1024) // 1MB
)

func shrinkFileSize(filePath string, maxSize int64) error {
	// Open the file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get initial file size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Check if file size is already under the limit
	if fileInfo.Size() <= maxSize {
		return nil // No need to remove lines
	}

	var offset, writeOffset int64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		offset += int64(len(scanner.Bytes())) + 1 // +1 for newline character

		// Check file size after removing each line
		if fileInfo.Size()-offset <= maxSize {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Set the file pointer to the offset
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	bufSize := 1024 // 1KB buffer
	buf := make([]byte, bufSize)

	for {
		bytesRead, readErr := file.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		if bytesRead == 0 {
			break
		}

		_, writeErr := file.WriteAt(buf[:bytesRead], writeOffset)
		if writeErr != nil {
			return writeErr
		}

		writeOffset += int64(bytesRead)

		if readErr == io.EOF {
			break
		}
	}

	// Truncate the file to remove the leftover content
	err = file.Truncate(writeOffset)
	if err != nil {
		return err
	}

	return nil
}

func NewLogFile(filePath string) (*os.File, error) {
	err := shrinkFileSize(filePath, LOG_FILE_MAX_SIZE)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare log file: %s: %w", filePath, err)
	}

	logFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %s: %w", filePath, err)
	}
	return logFile, nil
}
