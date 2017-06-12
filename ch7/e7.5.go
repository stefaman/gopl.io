package main

import(
	"fmt"
	"os"
	"net/http"
	"bufio"
	"io"

	"gopl.io/ch7/exercise7"

	"golang.org/x/net/html"
	"gopl.io/ch5/exercise5"
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
	reader = exercise7.LimitReader(reader, 1230)

	printReader(reader)
	printHTMLReader(reader)

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

func printReader(reader io.Reader) {
	input := bufio.NewScanner(reader)
	// input.Split(bufio.ScanWords)
	for input.Scan() {
		if err := input.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "printReader: %s", err)
		}
		fmt.Println(input.Text())
	}
}

func printHTMLReader(reader io.Reader) {
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "printHTMLReader: %v",  err)
		os.Exit(1)
	}
	exercise5.Pretty(doc)
}
