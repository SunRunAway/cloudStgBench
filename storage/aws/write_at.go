package aws

import (
	"io"

	"github.com/SunRunAway/cloudStgBench/lib"
)

type writerAt struct {
	//buf []byte
	size int64
}

func (this *writerAt) WriteAt(p []byte, off int64) (n int, err error) {
	//intOff := int(off)
	// limit := intOff + len(p)
	this.size += int64(len(p))
	return len(p), nil
}

func (this *writerAt) Reader() io.Reader {
	return lib.NewReader(this.size)
	// return bytes.NewReader(this.buf)
}

func NewWriterAt() *writerAt {
	return &writerAt{}
}
