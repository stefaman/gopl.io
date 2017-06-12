package main

import(
	"fmt"
	"flag"
	"gopl.io/ch7/exercise7"

)


func main()  {
	var temp = exercise7.TempFlag("temp", 0, "xxC xxK xxF")
	flag.Parse()
	fmt.Println(*temp)
}
