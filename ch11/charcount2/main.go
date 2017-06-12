// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 97.
//!+

// Charcount computes counts of Unicode characters.
package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
	"io"
	"unicode"
)

func main() {
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings

	in := bufio.NewReader(os.Stdin)
	counts, invalid, err := Count(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for rune, n := range counts {
		utflen[utf8.RuneLen(rune)]  += n
	}
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

func Count(input io.Reader) (map[rune]int, int, error){
	counts := make(map[rune]int)    // counts of Unicode characters
	invalid := 0                    // count of invalid UTF-8 characters

	in := bufio.NewReader(input)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("Count: %v\n", err)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
	}
	return counts, invalid, nil
}
//!-
