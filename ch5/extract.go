package main

import(
	"fmt"
	"os"
	"log"
	"gopl.io/ch5/links"
)

func main()  {
	links, err := links.Extract(os.Args[1])
	if err != nil {
		log.Fatalf("extract: %v", err)
	}
	for _, link := range links {
		fmt.Println(link)
	}
}
