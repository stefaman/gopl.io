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

	"runtime/trace"
)

var done = make(chan struct{})
func main()  {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	trace.Start(f)
	verbose := flag.Bool("v", false, "print verbose")
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	//Bug: 不能在此实现时钟信号，不能广播，是先读者先得
	// var tick <-chan time.Time
	// if *verbose {
	// 	tick = time.Tick(200 * time.Millisecond)
	// }

	go func(){
		fmt.Println("Press 'Enter' key to cancel progress")
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()
	var wg sync.WaitGroup
	for _, root := range roots {
		time.Sleep(100 * time.Millisecond)
		root := root
		var tick <-chan time.Time
		if *verbose {
			tick = time.Tick(200 * time.Millisecond)
		}
		statistic := make(chan *info)
		go stat(root, statistic)
		wg.Add(1)
		go func(){
			defer wg.Done()
			var nbytes, nfiles, ndirs int64
			exit:
			for {
				//bug 阻塞在<-statistic的goroutine不能取消
				// if cancelled() {
				// 	// for range statistic{} //drain channel
				// 		return
				// }
				select {
					case count, ok := <-statistic:
							if !ok {
								break exit
							}
							if count.Type ==  FileTypeNorm {
								nbytes += count.Size
								nfiles++
							}else {
								ndirs++
							}
					case <-tick:
						printSize(root, nbytes, nfiles, ndirs)
					case <- done:
						for range statistic{} //drain channel,wai for existing goroutines
						return
				}
			}
			printSize(root, nbytes, nfiles, ndirs)
		}()
	}
	// defer panic("in main")//测试程序退出时goroutine是否清理干净
	wg.Wait()
	trace.Stop()
}

type info struct {
	Type int
	Size int64
}
const(
	FileTypeNorm = iota
	FileTypeDir
)

func printSize(root string, nbytes, nfiles, ndirs int64)  {
	fmt.Printf("%v %s: %d directories %d files, %.1f GB\n", time.Now(), root, ndirs, nfiles,  float64(nbytes)/1e9)
}

func cancelled() bool  {
	select {
		case <-done:
			return true
		default:
			return false
	}
}

func stat(root string, statistic chan<- *info)  {
	fmt.Println("gouroutine begin")

	handler := func(entry os.FileInfo){
		var count = new(info)
		if !entry.IsDir() {
			count.Type = FileTypeNorm
			count.Size = entry.Size()
		} else{
			count.Type = FileTypeDir
			// count.N <- len(dirents(entry.Name()))
		}
		statistic <- count
	}
	walkDir(root, handler)
	close(statistic)
	fmt.Println("gouroutine over")
}

func walkDir(dir string, handler func(os.FileInfo)) {
	var wg sync.WaitGroup
	wg.Add(1)
	go forEachEntry(dir, &wg, handler)
	wg.Wait()
}

func forEachEntry(dir string, wg *sync.WaitGroup, f func(os.FileInfo)) {
	defer wg.Done()
	if cancelled() {
		return
	}
	entries := dirents(dir)
	for _, entry := range entries {
		f(entry)
		if entry.IsDir() {
			wg.Add(1)
			d := filepath.Join(dir, entry.Name())
			go forEachEntry(d, wg, f)
		}
	}
}

//sema is a counting semaphore for limiting  concurrency in dirents
var sema = make(chan struct{}, 20)
func dirents(dir string) []os.FileInfo {
	//bug 阻塞在获取sema的goroutine不能取消了
	// if cancelled() {
	// 	return nil
	// }
	// sema <- struct{}{}
	select {
		case <- done:
			return nil
		default:
			sema <- struct{}{}
	}
	defer func(){<- sema}()
	// time.Sleep(100*time.Millisecond) //model heavy CPU time
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return nil
	}
	return entries
}
