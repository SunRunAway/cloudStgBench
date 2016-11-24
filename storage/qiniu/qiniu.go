package qiniu

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/SunRunAway/cloudStgBench/lib"
	"github.com/SunRunAway/cloudStgBench/storage"

	"golang.org/x/net/context"

	"qiniupkg.com/api.v7/kodo"
)

var (
	fnameNum     = uint32(0)
	bucketDomain string
	endPoint     string
)

type QStg struct {
	bucket *kodo.Bucket
	ctx    context.Context
	client *http.Client
}

func (self QStg) Put(r io.Reader, size int64) (fileName string, err error) {
	fname := atomic.AddUint32(&fnameNum, 1)
	fileName, err = self.upload(fmt.Sprint(fname), r, size)
	return
}

func (self QStg) InitFileList(n int, size int64) (fileList []string, err error) {
	fileList = make([]string, n)
	pool := make(chan struct{}, 50)

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := range fileList {
		pool <- struct{}{}
		go func(i int) {
			defer func() {
				<-pool
				wg.Done()
			}()

			r := lib.NewReader(size)
			fileList[i], err = self.upload(strconv.FormatInt(int64(i), 10), r, size)
			if err != nil {
				fmt.Printf("upload with err %#v\n", err)
				return
			}
		}(i)
	}
	wg.Wait()
	return
}

func (self QStg) Get(fileName string) (io.ReadCloser, error) {
	downloadURL := endPoint + "/" + fileName

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return nil, err
	}

	req.Host = bucketDomain

	resp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (self QStg) upload(filename string, r io.Reader, size int64) (ret string, err error) {
	err = self.bucket.Put(self.ctx, nil, filename, r, size, nil)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func init() {
	kodo.SetMac(os.Getenv("QINIU_ACCESS_KEY_ID"), os.Getenv("QINIU_SECRET_KEY"))
	zoneEnv := os.Getenv("QINIU_ZONE")
	zone, err := strconv.ParseInt(zoneEnv, 10, 64)
	if err != nil {
		zone = 0
	}

	bucketDomain = os.Getenv("QINIU_BUCKET_DOMAIN")
	endPoint = os.Getenv("QINIU_ENDPOINT")

	c := kodo.New(int(zone), nil)
	bucket := c.Bucket(os.Getenv("QINIU_BUCKET"))
	ctx := context.Background()

	client := lib.NewClient()
	storage.RegisterStorage("qiniu", QStg{&bucket, ctx, client})
}
