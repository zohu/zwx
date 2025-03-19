package utils

import "sync"

type Buffer []byte

var bufPool = sync.Pool{
	New: func() any {
		b := make(Buffer, 0, 1024)
		return &b
	},
}

func NewBuffer() *Buffer {
	return bufPool.Get().(*Buffer)
}

func (b *Buffer) Free() {
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}
func (b *Buffer) Write(bytes []byte) (int, error) {
	*b = append(*b, bytes...)
	return len(bytes), nil
}

func (b *Buffer) WriteByte(char byte) error {
	*b = append(*b, char)
	return nil
}

func (b *Buffer) WriteString(str string) (int, error) {
	*b = append(*b, str...)
	return len(str), nil
}

func (b *Buffer) WriteStringIf(ok bool, str string) (int, error) {
	if !ok {
		return 0, nil
	}
	return b.WriteString(str)
}
func (b *Buffer) String() string {
	return string(*b)
}
