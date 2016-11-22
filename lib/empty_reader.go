package lib

import "io"

type emptyReader struct {
}

func (r *emptyReader) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func NewReader(size int64) io.Reader {
	return io.LimitReader(&emptyReader{}, size)
}
