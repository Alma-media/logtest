package adapter

import (
	"errors"
	"io"
	"sync"
)

var errWriteFailure = errors.New("failed to write")

// BufferingWriter is an adapter (implements io.Wrapper interface) using a buffer
// to keep messages in case of failure. Every time when `Write()` is called it tries
// to send the data starting with buffered (failed) messages to keep the natural order
// of messages.
type BufferingWriter struct {
	*Buffer
	io.Writer
}

// NewBufferingWriter creates a new buffering adapter.
func NewBufferingWriter(writer io.Writer, buffer *Buffer) *BufferingWriter {
	return &BufferingWriter{
		Buffer: buffer,
		Writer: writer,
	}
}

// Write writes len(p) bytes from p to the the data stream.
func (writer *BufferingWriter) Write(p []byte) (int, error) {
	var (
		msg = make([]byte, len(p))
		n   = copy(msg, p)
	)

	// it should not return an error if something goes wrong!
	// othervise we will get an internal error messages in stdout:
	// "Failed to write to log, ...(reason)..."
	writer.Buffer.Append(msg).Flush(writer.Writer)

	return n, nil
}

// FailingWriter is guaranteed to return an error on a certain call of Write() method.
type FailingWriter struct {
	mu sync.Mutex

	w              io.Writer
	called, failOn int
}

// NewFailingWriter wraps provided io.Writer to fail on `failOn` Write() call.
func NewFailingWriter(w io.Writer, failOn int) *FailingWriter {
	return &FailingWriter{
		w:      w,
		failOn: failOn,
	}
}

// Write writes len(p) bytes from p to the the data stream.
func (w *FailingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	defer func() { w.called++ }()

	if w.called >= w.failOn {
		return 0, errWriteFailure
	}

	return w.w.Write(p)
}
