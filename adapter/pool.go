package adapter

import (
	"bytes"
	"sync"
)

// BufferPool describes the pull of buffers.
type BufferPool interface {
	Put(*bytes.Buffer)
	Get() *bytes.Buffer
}

// NoPool implements BufferPool interface but does not provide any optimization.
// It returns new bytes.Buffer every time it is requested by `Get()`.
type NoPool struct{}

// Put does nothing (required to implement BufferPool).
func (p NoPool) Put(buf *bytes.Buffer) {}

// Get returns a new *bytes.Buffer every time it is called.
func (p NoPool) Get() *bytes.Buffer { return new(bytes.Buffer) }

// SyncPool reduces memory allocations reusing existing byte buffers for log entries.
type SyncPool struct{ pool *sync.Pool }

// NewSyncPool creates new empty SyncPool.
func NewSyncPool() *SyncPool {
	return &SyncPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Put a buffer back to the pool.
func (p *SyncPool) Put(buf *bytes.Buffer) { buf.Reset(); p.pool.Put(buf) }

// Get a buffer from the pool.
func (p *SyncPool) Get() *bytes.Buffer { return p.pool.Get().(*bytes.Buffer) }
