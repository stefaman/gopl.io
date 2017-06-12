package main

import(
	"fmt"
	"os"
	"net/http"
	"bufio"

	"gopl.io/ch7/exercise7"
	// "golang.org/x/net/html"
	// "gopl.io/ch5/exercise5"
	// "strings"

)

func main()  {


	url := os.Args[1]
	text, err := readUrl(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	// fmt.Println(string(text))
	reader := exercise7.NewReader(string(text))

	input := bufio.NewScanner(reader)
	// input.Split(bufio.ScanWords)
	for input.Scan() {
		if err := input.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading %s: %s", url, err)
		}
		fmt.Println(input.Text())
	}

	// doc, err := html.Parse(reader)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "parsing %s as html %v", url, err)
	// 	os.Exit(1)
	// }
	//
	// exercise5.Pretty(doc)

}

func readUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	input := bufio.NewScanner(resp.Body)
	var bs []byte
	for input.Scan() {
		if err := input.Err(); err != nil {
			return bs, fmt.Errorf("reading %s: %s", url, err)
		}
		bs = append(bs, []byte(input.Text() + "\n")... )
	}
	return bs, nil

}
