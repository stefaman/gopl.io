package main

import(
	"fmt"
	"io"
	"bufio"
	"unicode"
	"unicode/utf8"
	"os"
)

func main()  {
	// charCount(os.Stdin)
	wordCount(os.Stdin)
}

func charCount(f *os.File) {
	countChar := make(map[rune]int)
	countLen := [utf8.UTFMax+1]int{}
	var invalidChar, totalChar int
	// countCategory := [10]int{}
	// const(
	// 	letter := iota
	// 	digit
	// 	other
	// )
	countCategory := make(map[string]int)

	input := bufio.NewReader(f)
	for {
		r, n, err := input.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read file %s:%v", f.Name(), err)
			os.Exit(1)
			// return false
		}
		totalChar++
		if r == unicode.ReplacementChar && n == 1 {
			invalidChar++
			continue
		}
		countChar[r]++
		countLen[n]++
		switch  {
		case unicode.IsLetter(r):
			countCategory["letter"]++
		case unicode.IsDigit(r):
			countCategory["digit"]++
		case unicode.IsPunct(r):
			countCategory["punct"]++
		case unicode.IsSpace(r):
			countCategory["space"]++
		default:
			countCategory["other"]++
			fmt.Printf("%q", r)
		}
	}
	fmt.Printf("total char are %d\n", totalChar)
	fmt.Printf("invalid char: %d\n", invalidChar)
	fmt.Println("\nchar\tcounts")
	for char, num := range countChar {
		fmt.Printf("%q\t%d\n", char, num)
	}
	fmt.Println("\nlen\tcounts")
	for len, num := range countLen {
		fmt.Printf("%d\t%d\n", len, num)
	}
	fmt.Println("\ncate\tcounts")
	for cate, num := range countCategory {
		fmt.Printf("%s\t%d\n", cate, num)
	}
}

func wordCount(f *os.File)  {
	var wordNum int
	words := make(map[string]int)

	input := bufio.NewScanner(f)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		words[input.Text()]++
		wordNum++
	}
	fmt.Printf("Total %d words\n", wordNum)
	fmt.Println("word\tcount\tpercent\n")
	for word, num := range words{
		fmt.Printf("%s\t%d\t%.2f\n", word, num, 100 * float64(num) /float64(wordNum))
	}
}
