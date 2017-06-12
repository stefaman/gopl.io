package exercise7

import(
	"io"
)


type reader struct{
	buf []byte
}

func (r *reader) Read(p []byte) (n int, err error)  {
	// lenRead := len(p)
	// lenBuf := len(buf)
	buf := (*r).buf
	if len(p) >= len(buf) {
		err = io.EOF
	}
	n = copy(p, buf)

	(*r).buf = buf[n:]
	return
}

func NewReader(s string) io.Reader {
	var r reader
	r.buf = make([]byte, len(s))
	copy(r.buf, s)
	return &r

}
