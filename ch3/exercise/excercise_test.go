package exercise

import(
	"bytes"
	"fmt"
	"os"
	"testing"
)

const(
	point = '.'
	seqChar = ','
	width = 3
)

func main()  {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Use %s <number>...\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Println(comma(os.Args[1]))
}

func comma(s string) string  {
	var buf  bytes.Buffer

	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		buf.WriteByte(s[0])
		s = s[1:]
	}

	lenInt, lenStr := -1, len(s)
	for index := 0; index < len(s); index++ {
		if s[index] == point {
			lenInt = index
			break;
		}
	}
	if lenInt == -1 {
		lenInt = lenStr
	}

	remainder := lenInt % width
	if lenInt > 0 && remainder == 0 {
		remainder = width
	}
	buf.WriteString(s[:remainder])
	for index := remainder; index < lenInt; index += width {
		buf.WriteByte(seqChar)
		buf.WriteString(s[index : index + width])
	}

	buf.WriteString(s[lenInt:lenStr])
	return buf.String()
	}

	func isAnagrams(s, t string) bool  {
		if len(s) == 0 || len(t) == 0 {
			return false
		}
		for i, j := 0, len(t)-1; i < len(s) && j >= 0; i++{
			if s[i] != t[j] {
				return false
			}
			j--
		}
		return true
	}

func TestIsAnagrams(t *testing.T)  {
	testTrue :=
		isAnagrams("abc", "cba") &&
		isAnagrams("a", "a") &&
		isAnagrams("aaa", "aaa")

	testFalse :=
		isAnagrams("", "") ||
		isAnagrams("", "x") ||
		isAnagrams("x", "") ||
		isAnagrams("abc", "abc") ||
		isAnagrams("abc", "cab")

	if !testTrue || testFalse {
		t.Errorf("isAnagrams failed.")
	}

}
