

// e5.7
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
	"gopl.io/ch5/exercise5"

)

//!+
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "use main <url> <tagName>...\n")
		os.Exit(1)
	}
	url := os.Args[1]
	doc, err := exercise5.ParseUrl(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	tags := os.Args[2:]
	nodes := elementByTagName(doc, tags...)
	// nodes = append(nodes, nil) // error, append "nil" node
	// ok, nothing done
	// nodes = append(nodes, nil...)
	// nodes = append(nodes,)
	// nodes = append(nodes)



	for _, node := range nodes {
		exercise5.Pretty(node)
		// fmt.Println(node.Data)
	}
}

func elementByTagName(node *html.Node, tags ...string) []*html.Node {
	nodes := []*html.Node{}
	for _, tag := range tags {
		elements := exercise5.ElementByTag(node, tag)
		nodes = append(nodes, elements...)
	}
	return nodes
}
