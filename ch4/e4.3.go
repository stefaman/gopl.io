package main

import(
	"fmt"
	"testing"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func main()  {

	//e4.3
	fmt.Println("\nTest e4.3")
	s := []int{1,2,3,4,5,6,7,8,9}
	// s = []int{}
	reverseSliceInt(s)
	a := [...]int{1,2,3,4,5}
	reverseArrayInt(&a)
	fmt.Println(s, a)
	s=reverseSliceInt(s)

	//e4.4
	fmt.Println("\nTest e4.4")
	step, _ := strconv.Atoi(os.Args[1])
	fmt.Println("first", s)
	fmt.Println(step, rotateSliceInt(s, step),)
	fmt.Println(-step, rotateSliceInt(s, -step))

	//e4.5
	fmt.Println("\nTest e4.5")
	// fmt.Println(`"" == "" is:`, ""=="")
	// fmt.Println(`"" == "a" is:`, ""=="a")
	strArray := []string{"abc", "abc", "abc", "abd", "", "", "abc", "abc", ""}
	fmt.Printf("%q\n", strArray)
	fmt.Printf("%q\n", eliminateDuplicates(strArray))

	//e4.6
	fmt.Println("\nTest e4.6")
	byteSlice := []byte("a  \u8521 \t\n\f\v\r \U0001f34c c\u0085\u00a0\xe8\x94\u00a0")
	fmt.Printf("%s:%[1]q[% [1]X]\n", byteSlice)
	byteSlice = squashSpaces(byteSlice)
	fmt.Printf("%s:%[1]q[% [1]X]\n", byteSlice)

	//e4.7
	fmt.Println("\nTest e4.7")
	strUTF8 := "a \u8521 b \U0001f34c \xe8\x94\u8521 c \u00A3"
	str := []byte(strUTF8)
	fmt.Printf("%s:%[1]q[% [1]X]\n", str)
	reverseSliceUTF8(str)
	str = []byte(strUTF8)
	fmt.Printf("Use reverseSliceUTF8,get:%s:%[1]q[% [1]X]\n",str)
	reverseSliceUTF8V2(str)
	fmt.Printf("Use reverseSliceUTF8V2,get:%s:%[1]q[% [1]X]\n", str)
	// reverseSliceUTF8(str)


}



func reverseSliceInt(s []int) []int {
	for i, j:= 0, len(s)-1; i < j; i, j= i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//e4.3
func reverseArrayInt(s *[5]int) *[5]int {
	for i, j:= 0, len(s)-1; i < j; i, j= i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//e4.4
// step < 0, rotate to right/high index
func rotateSliceInt(s []int, step int) []int {
	lenS := len(s)
	if lenS == 0 || step == 0 {
		 return s
	}
	step = step % lenS
	lenT := step
	if lenT < 0 {
		lenT = -lenT
	}
	tmp := make([]int, lenT)
	if step > 0 {
		copy(tmp, s[:lenT]) //;fmt.Println("s1", s, tmp)
		copy(s, s[lenT:]) //;fmt.Println("s2",s, tmp)
		copy(s[lenS-lenT:], tmp) //;fmt.Println("s3",s, tmp)
	}else {
		reverseSliceInt(s)
		reverseSliceInt(s[:lenT])
		reverseSliceInt(s[lenT:])
	}
	return s
}

//4.5
func eliminateDuplicates( s []string) []string {
	j := 0
	lenS := len(s)
	for i := 0; i < lenS -1; i++ {
		if s[i] == s[i+1] {
			// i++ //error
		}else {
			s[j] = s[i]
			j++
		}
		// fmt.Printf("%q\ni is %d, j is %d\n", s, i, j)
	}
	s[j] = s[lenS-1];j++
	// fmt.Println(s)
	return s[:j]
}

//e4.6
func squashSpaces(s []byte) []byte  {
	j := 0
	preIsSpace := false
	lenS := len(s)
	for i := 0; i < lenS; {
		r, size := utf8.DecodeRune(s[i:])
		if unicode.IsSpace(r){
			if !preIsSpace {
				s[j] = 0x20
				j++
				preIsSpace = true
			}
		}else{
			preIsSpace = false
			// utf8.EncodeRune(s[j:], r)
			copy(s[j:], s[i:i+size])
			j += size
		}
		i += size
	}
	return s[:j]
}

//e4.7
// asume utf-8 encoding is not error
func reverseSliceUTF8(s []byte) ([]byte, bool){

	for i :=0; i < len(s); {
		switch {
		case s[i] & 0x80 == 0:
			i++
		case s[i] & 0xe0 == 0xc0:
			s[i], s[i+1] = s[i+1], s[i]
			i +=2
		case s[i] & 0xf0 == 0xe0:
			s[i], s[i+2] = s[i+2], s[i]
			i +=3
		case s[i] & 0xf8 == 0xf0:
			s[i], s[i+1], s[i+2], s[i+3] =
				s[i+3], s[i+2], s[i+1], s[i]
			i += 4
		default:
			return s, false
		}
	}
	for i, j:= 0, len(s)-1; i < j; i, j= i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s, true
}

// think of utf-8 encoding errors
func reverseSliceUTF8V2 (s []byte) []byte  {
	for index := 0; index < len(s);  {
		_, size := utf8.DecodeRune(s[index:])
		for i, j:= index, index+size-1; i < j; i, j= i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
		index += size
	}
	for i, j:= 0, len(s)-1; i < j; i, j= i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}


func TestReverseSliceInt(t testing.T) {
	s := []int{1,2,3,4,5}
	// sr := []int{5,4,3,2,1}

	s = reverseSliceInt(s)
}
