package adapter

import (
	"bytes"
	"container/list"
	"io"
)

var defaultPool = NewSyncPool() // NoPool{} //

// Buffer is a write buffer. It keeps the order of messages.
// TODO:
// - maybe use []byte instead or linked list
// - benchmark with a slice and linked list
// - shift first element once we reach the maximum capacity
// ISSUES:
// - we cannot use []byte since all the messages are protobuf encoded
// so it cannot be a sustained stream
type Buffer struct {
	pool BufferPool
	list *list.List

	capacity int
}

// NewBuffer creates a buffer with provided capacity.
func NewBuffer(capacity int) *Buffer {
	return &Buffer{
		pool:     defaultPool,
		list:     list.New(),
		capacity: capacity,
	}
}

// Append new element to the buffer (chainable).
func (buff *Buffer) Append(p []byte) *Buffer {
	if buff.capacity == 0 {
		return buff
	}

	if buff.list.Len() == buff.capacity {
		buff.pool.Put(buff.list.Remove(buff.list.Front()).(*bytes.Buffer))
	}

	b := buff.pool.Get()

	b.Write(p)
	buff.list.PushBack(b)

	return buff
}

// Len returns buffer length (number of lines).
func (buff Buffer) Len() int { return buff.list.Len() }

// String representation of a buffer (for testing / debugging purposes only).
// NOTE: this is not the most efficient way to build a string so never use it in production.
func (buff Buffer) String() string {
	var output []byte

	for element := buff.list.Front(); element != nil; element = element.Next() {
		output = append(output, element.Value.(*bytes.Buffer).Bytes()...)
	}

	return string(output)
}

// Flush the buffer into provided io.Writer.
func (buff *Buffer) Flush(w io.Writer) (written int, err error) {
	for element := buff.list.Front(); element != nil; element = element.Next() {
		b := element.Value.(*bytes.Buffer)

		n, err := w.Write(b.Bytes())
		if err != nil {
			return written, err
		}

		defer buff.list.Remove(element)

		buff.pool.Put(b)

		written += n
	}

	return
}
