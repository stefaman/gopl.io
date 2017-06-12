package exercise7

import(
	"io"
	"bytes"
)

type T struct {
	w io.Writer
	count *int64
	Words *int64
	Lines *int64
}


func (t T) Write(p []byte) (int, error) {
	*t.count += int64(len(p))
	*t.Words += countWords(p)
	*t.Lines += countLines(p)

	return t.w.Write(p)

}

func CountingWrite(w io.Writer) (T, *int64) {
	var cw T
	cw.w = w
	cw.count = new(int64)
	cw.Words = new(int64)
	cw.Lines = new(int64)
	return cw, cw.count
}

func countWords(p []byte) int64 {
	words := bytes.Split(p, []byte(" "))
	return int64(len(words))

}

func countLines(p []byte) int64  {
	lines := bytes.Split(p, []byte("\n"))
	return int64(len(lines))
}
