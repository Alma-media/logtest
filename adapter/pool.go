package adapter

import (
	"bytes"
	"sync"
)

type BufferPool interface {
	Put(*bytes.Buffer)
	Get() *bytes.Buffer
}

type DefaultPool struct{ pool *sync.Pool }

func NewBufferPool() *DefaultPool {
	return &DefaultPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (p *DefaultPool) Put(buf *bytes.Buffer) { buf.Reset(); p.pool.Put(buf) }

func (p *DefaultPool) Get() *bytes.Buffer { return p.pool.Get().(*bytes.Buffer) }
