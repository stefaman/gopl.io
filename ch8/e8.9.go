package main

import(
	"fmt"
	"flag"
	"path/filepath"
	"time"
	"io/ioutil"
	"os"
	"log"
	"sync"
)

func main()  {
	verbose := flag.Bool("v", false, "print verbose")
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(200 * time.Millisecond)
	}
	done := make(chan struct{})
	for _, root := range roots {
		go func(root string){
			sizeChan := sizeChan(root)
			var size, n int64
			exit:
			for {
				select{
				case s, ok := <-sizeChan:
						if !ok {
							break exit
						}
						size += s
						n++
					case <-tick:
						printSize(root, size, n)
				}
			}
			printSize(root, size, n)
			done <- struct{}{}
		}(root)
	}
	for range roots {
		<-done
	}
}

func printSize(root string, size, n int64)  {
	fmt.Printf("%s: %d files, %.1f GB\n", root, n, float64(size)/1e9)
}

func sizeChan(root string) <-chan int64 {
	size := make(chan int64)
	go getSize(root, size) //忘记加go, 整个程序deadlock
	return size
}

func getSize(root string, size chan<- int64)  {
	deal := func(entry os.FileInfo){
		if !entry.IsDir() {
			size <- entry.Size()
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go forEachEntry(root, &wg, deal)
	wg.Wait()
	close(size)
}

func forEachEntry(dir string, wg *sync.WaitGroup, f func(os.FileInfo)) {
	entries := getEntries(dir)
	for _, entry := range entries {
		f(entry)
		if entry.IsDir() {
			wg.Add(1)
			d := filepath.Join(dir, entry.Name())
			go forEachEntry(d, wg, f)
		}
	}
	wg.Done()
}

var tokens = make(chan struct{}, 20)
func getEntries(dir string) []os.FileInfo {
	tokens <- struct{}{}
	entries, err := ioutil.ReadDir(dir)
	<- tokens
	if err != nil {
		log.Print(err)
		return nil
	}
	return entries
}
