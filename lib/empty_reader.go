package lib

import (
	"io"
)

type emptyReader struct {
}

func (r *emptyReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func (r *emptyReader) ReadAt(p []byte, off int64) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func NewReader(size int64) io.Reader {
	return io.NewSectionReader(&emptyReader{}, 0, size)
}
