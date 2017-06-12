package main

import(
	"fmt"
	"os"
	"bytes"
)


func main()  {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "use main <seq> <string>...\n")
		os.Exit(1)
	}

	fmt.Printf("%s\n", join(os.Args[1], os.Args[2:]...))
}

func join(seq string, strs ...string) string  {
	// buf := new(bytes.Buffer)
	buf := bytes.Buffer{}
	for i, str := range strs {
		buf.WriteString(str)
		if i != len(strs) -1 {
			buf.WriteString(seq)
		}
	}
	return buf.String()

}
