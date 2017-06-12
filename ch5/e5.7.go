

// e5.7
package main

import (
	// "fmt"
	"os"

	// "golang.org/x/net/html"
	// "bytes"
	// "net/http"
	"gopl.io/ch5/exercise5"
	"log"
	// "time"
)

//!+
func main() {

	for _, url := range os.Args[1:] {

		doc, err := exercise5.ParseUrl(url)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "%s\n", err)
			log.SetPrefix("main: ")
			log.SetFlags(0)
			log.Printf("parseUrl: %v", err)
			continue
		}

		//e5.7
		exercise5.Pretty(doc)

	}


}
