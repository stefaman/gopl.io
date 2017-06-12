package test

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"unsafe"

	"gopl.io/ch12/display"
)

type E struct {
	eb bool
	ec complex128
	ei int16
}

func TestAlign(t *testing.T) {
	type st struct {
		b    bool
		i    int16
		c    complex128
		slia []int
		slib []complex128
		s    struct{ bb bool }
		E
	}
	var s st
	t.Log(unsafe.Sizeof(s), unsafe.Alignof(s.slia), unsafe.Alignof(s.slib))
	t.Log(unsafe.Offsetof(s.ec), unsafe.Offsetof(s.E.ec))
	printSt(s)
	// printSt(s.E)
	t.Log(unsafe.Sizeof([10]int{}), unsafe.Sizeof(1+2i), unsafe.Sizeof("[10]int{}"))
	sp := reflect.ValueOf(&s).Elem()
	t.Log(sp.UnsafeAddr() == sp.Addr().Pointer())

	var i = W("string...")
	var w io.Writer = i
	wv := reflect.ValueOf(&w).Elem() //interface kind Value
	// t.Logf("%x", wv.InterfaceData())
	cValPtr := unsafe.Pointer(wv.InterfaceData()[1])
	cTypPtr := unsafe.Pointer(wv.InterfaceData()[0])
	cTypValue := reflect.NewAt(reflect.TypeOf(reflect.TypeOf(0)), cTypPtr).Elem() //Type Value
	cTyp := cTypValue.Interface().(reflect.Type)                                  //Type
	cVal := reflect.NewAt(wv.Elem().Type(), cValPtr).Elem()
	p := *(*W)(cValPtr)
	t.Log(p == i)                     //true
	t.Log(cVal.Interface().(W), cTyp) //string... io.Writer
	display.Display("cTyp", cTyp)
	display.Display("wv", wv)
	display.Display("wvE", wv.Elem())
	display.Display("cVal", cVal)
}

type W string

func (w W) Write([]byte) (int, error) { return 1, nil }

func printSt(s interface{}) {
	sv := reflect.ValueOf(s)
	sTyp := sv.Type()
	fmt.Println(sTyp, "s, a af: ", sTyp.Size(), sTyp.Align(), sTyp.FieldAlign())
	for i := 0; i < sv.NumField(); i++ {
		fmt.Println(sTyp.Field(i))
		fieldTyp := sTyp.Field(i).Type
		fmt.Println(fieldTyp.Size(), fieldTyp.Align(), fieldTyp.FieldAlign())
		if sTyp.Field(i).Anonymous {
			printSt(sv.Field(i).Interface())
		}
	}
}
