
package main



import (
	"fmt"
	"os"
	"runtime"
)

func main()  {
/*
	for i := 0; i < 5; i++ {
		defer func(i int){ test(i)}(i)
	}

	for i := 0; i < 5; i++ {
		defer 	func(){ test(i)}()
	}
*/
	test :=func(i int) (ret int){
		defer func(i int){ test(i)}(i)
		defer func() {
			fmt.Println("first defer\n")
			// j := 0
			// fmt.Println(i/j)
			if v := recover(); v == i {
				ret = i
				printStack()
			}else{
				panic(v)
			}
		}()
		panic(i)
	}

	fmt.Println(test(33))


}

func test(i int) {
	fmt.Println("in test: ", i)
}

func printStack(){
	buf := [4096]byte{}
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}
