package main

import(
	"fmt"
	"os"

	"gopl.io/ch7/exercise7"

)


func main()  {
	cw, p := exercise7.CountingWrite(os.Stdout)
	fmt.Printf("type is %T\n", cw)
	
	fmt.Fprintf(cw, "abc\n kk bc\ndd kk\n")
	fmt.Printf("write %d bytes\n", *p)
	c := exercise7.T(cw)
	fmt.Printf("write %d words, %d lines\n", *c.Words, *c.Lines)

	// cw.Write([]byte("kkk"))

}
