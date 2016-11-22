package mock

import (
	"io"
	"strconv"
	"sync"
	"time"
)

type MockStg struct{}

func (self MockStg) Put(r io.Reader, size int64) (fileName string, err error) {
	time.Sleep(0.1e9)
	return strconv.FormatInt(size, 10), nil
}

func (self MockStg) InitFileList(n int, size int64) (fileList []string, err error) {
	fileList = make([]string, n)
	for i := range fileList {
		fileList[i] = strconv.FormatInt(size, 10)
	}
	return
}

func (self MockStg) Get(fileName string) (io.Reader, error) {
	size, err := strconv.ParseInt(fileName, 10, 64)
	if err != nil {
		return nil, err
	}
	return io.LimitReader(&timeReader{dur: 0.1e9}, size), nil
}

type timeReader struct {
	dur  time.Duration
	once sync.Once
}

func (tr *timeReader) Read(p []byte) (int, error) {
	tr.once.Do(func() {
		time.Sleep(0.1e9)
	})
	return len(p), nil
}

// func init() {
// 	storage.RegisterStorage("mock", MockStg{})
// }
