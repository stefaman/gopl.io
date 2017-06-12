package exercise7

import(
	"fmt"
	"bytes"

)


type Tree struct{
	root *tNode
}



type tNode struct {
	value int;
	left, right *tNode
}


func (t *tNode) add(v int) *tNode {
	if t == nil {
		// return &tNode{v, nil, nil}
		t = new(tNode)
		t.value = v
		return t
	}
	if v < t.value {
		t.left = t.left.add(v)
		return t
	}
	if v > t.value {
		t.right = t.right.add(v)
		return t
	}
	return t

}

func (t *Tree) Add(values ...int)  {
	for _, v := range values {
		t.root = t.root.add(v)
	}
}

func (t *Tree) String() string {
	var buf bytes.Buffer
	seq := ", "
	t.root.string(&buf, seq)
	buf.Truncate(buf.Len() - len(seq))
	buf.WriteString("}")
	return "{" + buf.String()
}

func (t *tNode) string(buf *bytes.Buffer, seq string) {
	if t == nil {
		return
	}
	t.left.string(buf, seq)
	fmt.Fprintf(buf, "%d%s", t.value, seq)
	t.right.string(buf, seq)
	return
}
