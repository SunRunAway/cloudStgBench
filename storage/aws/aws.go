package aws

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/SunRunAway/cloudStgBench/lib"
	"github.com/SunRunAway/cloudStgBench/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSStg struct {
	sess       *session.Session
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

var (
	bucketName      = ""
	defaultFileName = "file_1"
	fnameNum        = uint32(0)
)

func (self AWSStg) Put(r io.Reader, size int64) (fileName string, err error) {
	fname := atomic.AddUint32(&fnameNum, 1)
	return self.upload(fmt.Sprint(fname), r, size)
}

func (self AWSStg) InitFileList(n int, size int64) (fileList []string, err error) {
	fileList = make([]string, n)

	for i := range fileList {
		r := lib.NewReader(size)
		fileList[i], err = self.upload(strconv.FormatInt(int64(i), 10), r, size)
		if err != nil {
			fmt.Printf("upload with err %#v\n", err)
			return
		}
	}

	return
}

func (self AWSStg) Get(fileName string) (io.Reader, error) {
	w := NewWriterAt()

	input := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(fileName),
	}

	_, err := self.downloader.Download(w, input)
	if err != nil {
		return nil, err
	}

	return w.Reader(), nil
}

func (self AWSStg) upload(filename string, r io.Reader, size int64) (ret string, err error) {
	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &filename,
		Body:   r,
	}

	_, err = self.uploader.Upload(upParams)
	if err != nil {
		return
	}

	return filename, nil
}

func init() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Printf("Register aws storage failed: %#v", err)
		return
	}

	bucketName = os.Getenv("AWS_BUCKET")

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.Concurrency = 1
	})

	downloader := s3manager.NewDownloader(sess, func(u *s3manager.Downloader) {
		u.Concurrency = 1
	})

	storage.RegisterStorage("aws", AWSStg{sess, uploader, downloader})
}
