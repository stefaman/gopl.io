package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
)

func main()  {
	counts := make(map[string]string)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "use %s <filename...>\n", "dup")
		return
	}

	for _, arg := range os.Args[1:] {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s error: %v\n", "dup", err)
			continue
		}
		countLines(f, counts, arg)
	}

	for line, names := range counts {
		fmt.Printf("%s: %s\n", line, names)
	}
}

func countLines(f *os.File, counts map[string]string, name string){
	input := bufio.NewScanner(f)
	for input.Scan() {
		if strings.Index(counts[input.Text()], name) == -1 {
			counts[input.Text()] += " " + name
		}
	}
}

//test
//test
