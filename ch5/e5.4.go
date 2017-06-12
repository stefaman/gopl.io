// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 123.

	// Outline prints the outline of an HTML document tree.
package main

import (
	"fmt"
	"os"

	// "golang.org/x/net/html"
	// "bytes"
	// "net/http"
	"gopl.io/ch5/exercise5"
	"log"
	"time"
)

//!+
func main() {

	for _, url := range os.Args[1:] {
		if err := exercise5.WaitForServer(url, 30 *time.Second); err != nil {
			log.Printf("site is shut down. %v", err)
			continue
		}
		doc, err := exercise5.ParseUrl(url)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "%s\n", err)
			log.SetPrefix("main: ")
			log.SetFlags(0)
			log.Printf("parseUrl: %v", err)
			continue
		}

		//e5.4
		for _, line := range exercise5.PrintHTML(doc, "", "    ") {
			fmt.Printf(line)
		}

	}


}
