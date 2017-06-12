// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch8/links"
	"flag"
	"sync"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(done <-chan struct{}, url string) []string {
	fmt.Println(url)
	select {
	case tokens <- struct{}{}: // acquire a token
	case <-done:
		return nil
	}
	list, err := links.Extract(done, url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

func cancelled(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

//!+
func main() {
	limit := flag.Int("depth", 0, "fetch depth limit")
	flag.Parse()
	done := make(chan struct{})
	go func(){
		fmt.Println("Press 'Enter' key to abort")
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	var list[][]string
	list = append(list, flag.Args())

	// Crawl the web concurrently.
	seen := make(map[string]bool)

	for ndepth := 0; ;ndepth++ {
		if (*limit > 0 && ndepth > *limit) || list == nil || cancelled(done) {
			break
		}
		fmt.Printf("----------level %d ---------------------------\n", ndepth)

		oneLevel := make(chan []string)
		var wg sync.WaitGroup
		for _, page := range list {
			for _, link := range page {
				if !seen[link] || cancelled(done) {
					wg.Add(1)
					seen[link] = true
					go func(link string) {
						defer wg.Done()
						select {
						case oneLevel <- crawl(done, link):
						case <-done:
							return
						}
					}(link)
				}
			}
		}
		go func(){
			wg.Wait()
			close(oneLevel)
		}()
		list = nil
		for page := range oneLevel {
			list = append(list, page)
		}
	}

}

//!-
