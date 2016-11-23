package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync/atomic"
	"time"

	"github.com/SunRunAway/cloudStgBench/lib"
	"github.com/SunRunAway/cloudStgBench/storage"

	_ "github.com/SunRunAway/cloudStgBench/storage/aws"
	_ "github.com/SunRunAway/cloudStgBench/storage/mock"
)

var smallSize = flag.Int64("ssize", 4*1024, "small size to test iops")
var concurrent = flag.Int("c", 20, "concurrent to test iops")
var duration = flag.Int("s", 10, "iops duration in second")
var largeSize = flag.Int64("lsize", 100*1024*1024, "large size to test speed")

func main() {
	flag.Parse()
	for name, stg := range storage.StgMap {
		fmt.Printf("==================> start testing %v\n", name)

		fmt.Println("-----put iops-----")
		testIops(func() (n int64, err error) {
			_, err = stg.Put(lib.NewReader(*smallSize), *smallSize)
			if err != nil {
				return 0, err
			}
			return *smallSize, nil
		})
		fmt.Println("-----done-----")

		fileList, err := stg.InitFileList(*concurrent*10, *smallSize)
		if err != nil {
			log.Fatalln("InitFileList failed:", err)
		}
		index := uint32(0)
		fmt.Println("-----get iops-----")
		testIops(func() (n int64, err error) {
			fileName := fileList[atomic.AddUint32(&index, 1)%uint32(len(fileList))]
			rc, err := stg.Get(fileName)
			if err != nil {
				return 0, err
			}
			defer rc.Close()
			return io.Copy(ioutil.Discard, rc)
		})
		fmt.Println("-----done-----")

		var fileName string
		fmt.Println("-----put speed-----")
		testSpeed(func() (n int64, err error) {
			fileName, err = stg.Put(lib.NewReader(*largeSize), *largeSize)
			if err != nil {
				return 0, err
			}
			return *largeSize, err
		})
		fmt.Println("-----done-----")
		if err != nil {
			log.Fatalln("stg.Put error", err)
		}

		fmt.Println("-----get speed-----")
		testSpeed(func() (n int64, err error) {
			rc, err := stg.Get(fileName)
			if err != nil {
				return 0, err
			}
			defer rc.Close()
			return io.Copy(ioutil.Discard, rc)
		})
		fmt.Println("-----done-----")

		fmt.Printf("==================> end testing %v\n", name)
	}
}

func testIops(invoke func() (n int64, err error)) {
	totalSize := int64(0)
	totalCount := int32(0)
	totalError := int32(0)
	oldSize := int64(0)
	oldCount := int32(0)
	oldError := int32(0)

	stop := make(chan struct{})

	go func() {
		for i := 0; i < *concurrent; i++ {
			go func() {
				for {
					select {
					case <-stop:
						return
					default:
						break
					}
					n, err := invoke()
					if err == nil {
						atomic.AddInt64(&totalSize, n)
						atomic.AddInt32(&totalCount, 1)
					} else {
						atomic.AddInt32(&totalError, 1)
					}
				}
			}()
		}
	}()

	t := time.NewTicker(time.Second)
	i := 0
	for {
		select {
		case <-t.C:
			i += 1
			size := atomic.LoadInt64(&totalSize)
			count := atomic.LoadInt32(&totalCount)
			errCnt := atomic.LoadInt32(&totalError)
			fmt.Printf("%v, iops: %v, throughput: %v, error: %v\n", i, count-oldCount, size-oldSize, errCnt-oldError)
			oldSize = size
			oldCount = count
			oldError = errCnt
			if i >= *duration {
				t.Stop()
				close(stop)
				return
			}
		}
	}
}

func testSpeed(invoke func() (n int64, err error)) {
	now := time.Now()
	n, err := invoke()
	dur := time.Since(now)
	if err != nil {
		fmt.Printf("testSpeed error: %v\n", err)
		return
	}
	fmt.Printf("testSpeed length: %v, cost: %v, speed: %v MB/s", n, dur, float64(n/1024/1024)/(float64(dur)/float64(time.Second)))
}
