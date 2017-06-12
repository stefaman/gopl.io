package main

import(
	"fmt"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"os"
)

// pc[i] is population count of i
var pc [256]byte

func init() {
	for i := 0; i < 256; i++ {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func main()  {
	s384 := flag.Bool("3", false, "use sha384")
	s512 := flag.Bool("5", false, "use sha512")
	s256 := flag.Bool("2", false, "[default]use sha256")
	diff := flag.Bool("d", false, "count difference bit of two sha digests")

	flag.Parse()
	if *diff {
		if len(flag.Args()) != 2 {
			fmt.Fprintf(os.Stderr, "Use %s -d <string1> <string2>\n", os.Args[0])
			os.Exit(2)
		}
		var c [2][32]byte
		for i, v := range flag.Args() {
			c[i] = sha256.Sum256([]byte(v))
			fmt.Printf("sha256 of %s is %x\n", v, c[i])
		}
		fmt.Printf("Total %d bits are difference\n", countDiffbitsSlice(c[0][:], c[1][:]))

		return

	}

	for _, v := range flag.Args() {
		switch {
		case *s512:
			fmt.Printf("sha512 of %s is %x\n", v, sha512.Sum512([]byte(v)))
		case *s384:
			fmt.Printf("sha384 of %s is %x\n", v, sha512.Sum384([]byte(v)))
		case *s256:
			fallthrough
		default:
			fmt.Printf("sha256 of %s is %x\n", v, sha256.Sum256([]byte(v)))
		}
	}
}

func countDiffSha256(a, b *[32]uint8) int {
	count := 0
	for index:= 0; index < 32; index++ {
		count += popCountSlice([]byte{a[index] ^ b[index]})
	}
	return count
}

func countDiffbitsSlice(a, b []byte) int {
	count := 0
	for index := range a {
		// count += popCountSlice([]byte{a[index] ^ b[index]})
		count += int(pc[ a[index] ^ b[index] ])
	}
	return count
}

func popCountSlice(s []byte) int {
	n := 0
	for _, x := range s{
		// for x != 0 {
		// 	x &= x - 1
		// 	n++
		// }
		n += int(pc[x])
	}
	return n
}
