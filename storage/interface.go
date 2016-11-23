package storage

import "io"

type Interface interface {
	Put(r io.Reader, size int64) (fileName string, err error)      // 单纯上传一个文件
	InitFileList(n int, size int64) (fileList []string, err error) // 上传n个文件，并且拿到文件列表， 内容无所谓，都为空白就行
	Get(fileName string) (io.ReadCloser, error)                    // 对单个文件名进行下载
}

var StgMap = make(map[string]Interface)

func RegisterStorage(name string, stg Interface) {
	if _, ok := StgMap[name]; ok {
		panic(name + " is registered")
	}
	StgMap[name] = stg
}
