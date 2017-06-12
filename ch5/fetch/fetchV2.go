// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 148.

// Fetch saves the contents of a URL into a local file.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"encoding/json"
)

//!+
// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	local := path.Base(resp.Request.URL.Path)
	switch local {
		case "/":
		local = "index.html"
		case ".":
	}
	fmt.Println(json.MarshalIndent(resp.Request.Header, "", " "))
	fmt.Printf("%#v\n", resp)
	// fmt.Printf("%#v\n", resp.Request.Header)
	// for k, v := range resp.Header {
	// 	fmt.Println(k, v)
	// }
	fmt.Println(resp.Request.URL, resp.Request.URL.Path, local,)

	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}

//e15.18
	close := func(){
		// Close file, but prefer error from Copy, if any.
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}
	n, err = io.Copy(f, resp.Body)
	defer close()
	return local, n, err
}

//!-

func main() {
	for _, url := range os.Args[1:] {
		local, n, err := fetch(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
	}
}
