package adapter

import (
	"bytes"
	"io"
)

var defaultPool = NewSyncPool()

// Buffer is a write buffer. It keeps the order of messages.
// ISSUES:
// - we cannot use []byte since all the messages are protobuf encoded
// so it cannot be a sustained stream
type Buffer struct {
	pool     BufferPool
	buff     []*bytes.Buffer
	capacity int
}

// NewBuffer creates a buffer with provided capacity.
func NewBuffer(capacity int) *Buffer {
	return &Buffer{
		pool:     defaultPool,
		buff:     make([]*bytes.Buffer, 0, capacity),
		capacity: capacity,
	}
}

// Append new element to the buffer (chainable).
func (buff *Buffer) Append(p []byte) *Buffer {
	if buff.capacity == 0 {
		return buff
	}

	// the buffer should be able to keep at least one message
	if len(buff.buff) == buff.capacity {
		buff.pool.Put(buff.buff[0])

		buff.buff = buff.buff[1:]
	}

	b := buff.pool.Get()

	b.Write(p)

	buff.buff = append(buff.buff, b)

	return buff
}

// Len returns buffer length (number of lines).
func (buff Buffer) Len() int { return len(buff.buff) }

// String representation of a buffer (for testing / debugging purposes only).
// NOTE: this is not the most efficient way to build a string so never use it in production.
func (buff Buffer) String() string {
	var output []byte

	for _, b := range buff.buff {
		output = append(output, b.Bytes()...)
	}

	return string(output)
}

// Flush the buffer into provided io.Writer.
func (buff *Buffer) Flush(w io.Writer) (written int, err error) {
	var index int

	defer func() {
		buff.buff = buff.buff[index:]

		// if the buffer was fully flushed
		if len(buff.buff) == 0 {
			// micro optimization to prevent a couple of slice allocations
			// after the fuffer is flushed
			buff.buff = make([]*bytes.Buffer, 0, buff.capacity)
		}
	}()

	for ; index < len(buff.buff); index++ {
		b := buff.buff[index]

		n, err := w.Write(b.Bytes())
		if err != nil {
			return written, err
		}

		buff.pool.Put(b)

		written += n
	}

	return
}
