package eval

import(
	"fmt"
	"bytes"
)

func (v Var) String() string {
	return string(v)
}

func (f literal) String() string {
	return fmt.Sprintf("%.2g", float64(f))
}

func (u unary) String() string {
	buf := new(bytes.Buffer)
	switch u.op {
	// case '+', '-':
	default:
		buf.WriteString("(")
		buf.WriteRune(u.op)
		buf.WriteString(u.x.String())
		buf.WriteString(")")
		return buf.String()
		// return fmt.Sprintf("unknow operator %q", u.op)
	}
}

func (f call) String() string {
		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "%s(", f.fn)
		for i, arg := range f.args {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(arg.String())
		}
		buf.WriteString(")")
		return buf.String()
}

func (b binary) String() string {
		buf := new(bytes.Buffer)
		buf.WriteString("(")
		buf.WriteString(b.x.String())
		fmt.Fprintf(buf, " %c ", b.op)
		buf.WriteString(b.y.String())
		buf.WriteString(")")
		return buf.String()
}
