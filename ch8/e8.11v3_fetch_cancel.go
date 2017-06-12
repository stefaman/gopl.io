package main

import(
	"os"
	"fmt"
	"net/http"
	"io"
	"sync"
	"bytes"

	// "errors"
)
//假设不知派生协程的数量，只有全部失败才最终返回错误，返回第一个成功的url
func fetch(urls... string) ([]byte, error) {
	done := make(chan struct{})
	over := make(chan struct{})
	door := make(chan struct{})
	errC := make(chan error)
	var wg sync.WaitGroup
	buf := new(bytes.Buffer)

	go func(){ door <- struct{}{} }()
	get := func(url string) {
		defer wg.Done()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			select {
			case <-done://do nothing
			case errC <- err:
			}
			return
		}
		req.Cancel = done
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			select {
			case <-done:
				//do nothing
			case errC <- err:
			}
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			select {
			case <-done:
				//do nothing
			case errC <- fmt.Errorf("get %s: %s", url, resp.Status):
			}
			return
		}
		select {
		case <-door://第一个进入
			close(done)//阻挡后面的进入
			// close(errC)//错误回收协程的循环提前退出//bug, 多次close
		case <-done:
			return
		}
		io.Copy(buf, resp.Body)
		// io.Copy(os.Stdout, resp.Body)
		fmt.Fprintln(os.Stdout, "test ", resp.Request.URL)
	}

	for _, url := range urls {
		wg.Add(1)
		go get(url)
	}

//回收error
	var errs []error
	var finalErr error
	go func(){
		wg.Wait()
		close(errC)
	}()
	go func(){
		for err := range errC {
			errs = append(errs, err)
		}
		select {
		case <- done:
			finalErr = nil
		default:
			//如果close(done)房子此处，bug, finalErr 写与最后的renturn存在race
			finalErr = fmt.Errorf("All urls are errors: %v", errs)
		}
		over<- struct{}{}
	}()

	// wg.Wait()//bug, 竞争，多处wg.Wait()
	//<- done //第一个成功返回，或者全部错误 //bug, 竞争，get()函数最后的处理
	<-over
	return buf.Bytes(), finalErr
}

func main()  {
	// bytes, err := fetch(os.Args[1:]...)
	// fmt.Println(string(bytes), err)
	_, err := fetch(os.Args[1:]...)
	fmt.Println( err)
}
