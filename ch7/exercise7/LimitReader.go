package exercise7

import(
	"io"
)


type readerT struct{
	remainder int64
	io.Reader
}



func (r *readerT) Read(p []byte) (n int, err error)  {
	if int64(len(p)) >= r.remainder {
		p = p[:r.remainder]
		n, err = r.Reader.Read(p)
		if err == nil {
			err = io.EOF
		}
	}else{
		n, err = r.Reader.Read(p)
	}

	(*r).remainder -= int64(n)
	return
}

func LimitReader(r io.Reader, n int64) io.Reader {
	var read readerT
	read.remainder = n
	read.Reader = r
	return &read

}
