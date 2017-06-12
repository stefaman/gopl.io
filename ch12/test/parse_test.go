package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"text/scanner"
)

var tokenStyle = map[rune]string{
	scanner.EOF:       "eof",
	scanner.Char:      "char",
	scanner.Comment:   "comment",
	scanner.Float:     "float",
	scanner.Int:       "int",
	scanner.Ident:     "indent",
	scanner.String:    "string",
	scanner.RawString: "rawString",
}

//搞半天，已经有库函数scanner.TokenString()
func tokenType(t rune) string {
	if s, found := tokenStyle[t]; found {
		return s
	}
	return fmt.Sprintf("%q", string(t))
}

func TestParser(t *testing.T) {
	var src = `
// This is scanned code.
var(
	someParsable = true
)
if a > 10 {
		str := "string\""
		3.14E12
		'蔡' kk/*comment*/dd
}` + "`a\"/*comment*/`"
	var s scanner.Scanner
	s.Filename = "example"
	s.Init(strings.NewReader(src))
	s.Mode = scanner.GoTokens | scanner.ScanInts & ^scanner.SkipComments
	// ^scanner.ScanFloats & //识别不了“10”， “3.14E12”
	// ^scanner.ScanRawStrings
	// ^scanner.ScanStrings & //将".."作为三个token
	// ^scanner.ScanInts &
	s.Mode = scanner.ScanFloats | scanner.ScanIdents | scanner.ScanStrings | scanner.ScanChars | scanner.ScanRawStrings |
		scanner.ScanComments |
		scanner.SkipComments
	var tok rune
	for tok != scanner.EOF {
		tok = s.Scan()
		str, _ := strconv.Unquote(s.TokenText())
		fmt.Println("At position", s.Pos(), ":", str, scanner.TokenString(tok))
	}
}
