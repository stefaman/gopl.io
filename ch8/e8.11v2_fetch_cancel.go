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
	door := make(chan struct{})
	errC := make(chan error, len(urls))

	go func(){door<- struct{}{}}()
	var wg sync.WaitGroup
	wg.Add(len(urls))
	get := func(url string) {
		defer wg.Done()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errC <- err
			return
		}
		req.Cancel = done
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errC <- err
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errC <- fmt.Errorf("get %s: %s", url, resp.Status)
			return
		}
		select {
		case <-door:
			close(done)
		case <-done:
			return
		}
		// io.Copy(w, resp.Body)
		fmt.Fprintln(w, resp.Request.URL)
	}

	for _, url := range urls {
		go get(url)
	}
	wg.Wait()
	select {
	case <-done:
		return nil
	default:
		//do nothing
	}
	close(errC)
	var errs []error
	for err := range errC {
		errs = append(errs, err)
	}
	return fmt.Errorf("All urls are errors: %v", errs)

}

func main()  {
	fmt.Println(fetch(os.Stdout, os.Args[1:]...))
}
