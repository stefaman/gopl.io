package thumbnail_test

import(

	"testing"
	"strconv"
	"log"
	// "gopl.io/ch8/thumbnail"

)
/*bash
x=1000; while [[ $x > 0 ]]; do { declare x=$(($x - 1)); cp ~/picture/fish.jpg ~/picture/$x.jpg; } done

rm ~/picture/*.thumb.jpg;go test -v -bench=BenchmarkMakeThumbnails2 -cpu=4 gopl.io/ch8/thumbnail;
ls ~/picture/*.thumb.jpg;


*/
const N = 1000
var files = make([]string, N)
func init(){

	for i := 0; i < N; i++ {
		files[i] = "/home/stefaman/picture/" + strconv.Itoa(i) + ".jpg"
	}
}

func BenchmarkMakeThumbnails1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeThumbnails(files)
	}
}

func BenchmarkMakeThumbnails2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeThumbnails2(files)
	}
}

func BenchmarkMakeThumbnails3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeThumbnails3(files)
	}
}

//bash
/*
rm ~/picture/500.jpg
rm ~/picture/*.thumb.jpg;
prlimit --nofile=1024 go test -v -bench=BenchmarkMakeThumbnails4 -cpu=4 -count=3 gopl.io/ch8/thumbnail;
*/

//goroutines leak, cause too many open files
func BenchmarkMakeThumbnails4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		log.Println(makeThumbnails4(files))
	}
}

//bash
/*
rm ~/picture/500.jpg
rm ~/picture/*.thumb.jpg;
prlimit --nofile=6556:6556 go test -v -bench=BenchmarkMakeThumbnails5 -cpu=4 -count=3 gopl.io/ch8/thumbnail;
*/

//avoid warning: open too many open files
func BenchmarkMakeThumbnails5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		log.Println(makeThumbnails5(files))
	}
}


//bash
/*
rm ~/picture/500.jpg
rm ~/picture/*.thumb.jpg;
prlimit --nofile=6556:6556 go test -v -bench=BenchmarkMakeThumbnails6 -cpu=4 -count=3 gopl.io/ch8/thumbnail;
*/

//avoid warning: open too many open files
func BenchmarkMakeThumbnails6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fs := make(chan string)
		go func(){
			for _, file := range files {
				fs <- file
			}
			close(fs)
			}()
		makeThumbnails6(fs)
	}
}
