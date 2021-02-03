package adapter

import "io"

// Buffer is a write buffer. It keeps the order of messages.
// TODO:
// - maybe use []byte instead or linked list
// - benchmark with a slice and linked list
// - shift first element once we reach the maximum capacity
// ISSUES:
// - we cannot use []byte since all the messages are protobuf encoded
// so it cannot be a sustained stream
type Buffer [][]byte

// NewBuffer creates a buffer with provided capacity.
func NewBuffer(capacity int) *Buffer {
	buff := make(Buffer, 0, capacity)

	return &buff
}

// Append new element to the buffer (chainable).
// TODO: use bytearray pool
func (buff *Buffer) Append(p []byte) *Buffer {
	if len(*buff) < cap(*buff) {
		*buff = append(*buff, p)

		return buff
	}

	if cap(*buff) > 0 {
		*buff = append((*buff)[1:], p)
	}

	return buff
}

// Len returns buffer length (number of lines).
func (buff Buffer) Len() int { return len(buff) }

// String representation of a buffer (for testing / debugging purposes only).
// NOTE: this is not the most efficient way to build a string so never use it in production.
func (buff Buffer) String() string {
	var output []byte

	for _, line := range buff {
		output = append(output, line...)
	}

	return string(output)
}

// Flush the buffer into provided io.Writer.
func (buff *Buffer) Flush(w io.Writer) (written int, err error) {
	var index int

	defer func() { *buff = (*buff)[index:] }()

	for ; index < len(*buff); index++ {
		n, err := w.Write((*buff)[index])
		if err != nil {
			return written, err
		}

		written += n
	}

	return
}
