package main

import(
	"fmt"

	"gopl.io/ch7/exercise7"

)


func main()  {
	var t exercsie7.Tree
	t.Add(5,4,3,2,1)
	fmt.Println(&t)
}
