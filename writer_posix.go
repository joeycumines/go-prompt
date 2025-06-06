//go:build !windows

package prompt

import (
	"syscall"
)

const flushMaxRetryCount = 3

// PosixWriter is a Writer implementation for POSIX environment.
// To control terminal emulator, this outputs VT100 escape sequences.
type PosixWriter struct {
	VT100Writer
	fd int
}

// Flush to flush buffer
func (w *PosixWriter) Flush() error {
	l := len(w.buffer)
	offset := 0
	retry := 0
	for {
		n, err := syscall.Write(w.fd, w.buffer[offset:])
		if err != nil {
			if retry < flushMaxRetryCount {
				retry++
				continue
			}
			return err
		}
		offset += n
		if offset == l {
			break
		}
	}
	w.buffer = []byte{}
	return nil
}

var _ Writer = &PosixWriter{}

var (
	// NewStandardOutputWriter returns Writer object to write to stdout.
	// This generates VT100 escape sequences because almost terminal emulators
	// in POSIX OS built on top of a VT100 specification.
	// Deprecated: Please use NewStdoutWriter
	NewStandardOutputWriter = NewStdoutWriter
)

// NewStdoutWriter returns Writer object to write to stdout.
// This generates VT100 escape sequences because almost terminal emulators
// in POSIX OS built on top of a VT100 specification.
func NewStdoutWriter() Writer {
	return &PosixWriter{
		fd: syscall.Stdout,
	}
}

// NewStderrWriter returns Writer object to write to stderr.
// This generates VT100 escape sequences because almost terminal emulators
// in POSIX OS built on top of a VT100 specification.
func NewStderrWriter() Writer {
	return &PosixWriter{
		fd: syscall.Stderr,
	}
}
