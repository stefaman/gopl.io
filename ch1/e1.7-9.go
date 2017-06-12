package main

import(
	"fmt"
	"os"
	"strings"
	"net/http"
	"io"
)


func main()  {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "use fetch <url>...\n")
		os.Exit(1)
	}
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		fetch(url)
	}
}

func fetch(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v", err)
		os.Exit(1)
	}


	_, errCopy := io.Copy(os.Stdout, resp.Body)
	resp.Body.Close()
	if errCopy != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", errCopy)
		os.Exit(1)
	}
	fmt.Printf("get http response status: %v\n", resp.Status)
}
