package main

import(
	"os"
	"fmt"
	"net/http"
	"io"
	// "io/ioutil"
	"sync"
)

func fetch(w io.Writer, urls... string) error {
	done := make(chan struct{})
	value := make(chan struct{})
	errC := make(chan error, len(urls))
	var wg sync.WaitGroup
	wg.Add(len(urls))
	get := func(url string) {
		defer wg.Done()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			select {
			case <-done:
				//do nothing
			case errC <- err:
			}
			// errC <- err
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
			// errC <- err
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			select {
			case <-done:
				//do nothing
			case errC <- fmt.Errorf("get %s: %s", url, resp.Status):
			}
			// errC <- fmt.Errorf("get %s: %s", url, resp.Status)
			return
		}
		select {
		case value<- struct{}{}:
		case <-done:
			return
		// default:
		// 	//do nothing
		}
		// done<- struct{}{}
		// close(done)
		io.Copy(w, resp.Body)
	}

	for _, url := range urls {
		go get(url)
	}
	<- value
	close(done)
	select {
	case <- value:
		close(done)
	case err:= <-errC:
		return err
	// default:
	// 	return nil
	}
	wg.Wait()
	return nil

}

func main()  {
	fmt.Println(fetch(os.Stdout, os.Args[1:]...))
}
