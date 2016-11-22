# cloudStgBench
Benchmark tool for cloud storage

```go
type Interface interface {
	Put(r io.Reader, size int64) (fileName string, err error) // 单纯上传一个文件
	InitFileList(n int, size int64) (fileList []string)       // 上传n个文件，并且拿到文件列表， 内容无所谓，都为空白就行
	Get(fileName string) (io.Reader, error)                   // 对单个文件名进行下载
}
```
