// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 97.
//!+

// Charcount computes counts of Unicode characters.
package charcount

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

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
