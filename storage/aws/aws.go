package aws

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/SunRunAway/cloudStgBench/lib"
	"github.com/SunRunAway/cloudStgBench/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSStg struct {
	S3 *s3.S3
}

var (
	bucketName = ""
	fnameNum   = uint32(0)
)

func (self AWSStg) Put(r io.Reader, size int64) (fileName string, err error) {
	fname := atomic.AddUint32(&fnameNum, 1)
	fileName, err = self.upload(fmt.Sprint(fname), r, size)
	return
}

func (self AWSStg) InitFileList(n int, size int64) (fileList []string, err error) {
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

func (self AWSStg) Get(fileName string) (io.ReadCloser, error) {

	input := s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(fileName),
	}

	resp, err := self.S3.GetObject(&input)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (self AWSStg) upload(filename string, r io.Reader, size int64) (ret string, err error) {

	bucket := bucketName
	key := filename
	input := s3.PutObjectInput{
		ACL:    aws.String("private"),
		Body:   r.(io.ReadSeeker),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if _, err = self.S3.PutObject(&input); err != nil {
		return "", err
	}
	return filename, nil
}

func init() {
	sess, err := session.NewSession(&aws.Config{
		DisableSSL:              aws.Bool(true),
		DisableComputeChecksums: aws.Bool(true),
		Endpoint:                aws.String(os.Getenv("AWS_CONFIG_ENDPOINT")),
		S3ForcePathStyle:        aws.Bool(true),
		HTTPClient:              lib.NewClient(),
	})
	if err != nil {
		fmt.Printf("Register aws storage failed: %#v", err)
		return
	}

	bucketName = os.Getenv("AWS_BUCKET")

	S3 := s3.New(sess)

	storage.RegisterStorage("aws", AWSStg{S3: S3})
}
