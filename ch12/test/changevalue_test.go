package test

import (
	"io"
	"os"
	"reflect"
	"testing"

	"gopl.io/ch12/display"
)

type W int

func (w W) Write(p []byte) (int, error) {
	return int(w), nil
}

func TestChangeValue(t *testing.T) {
	x := 2
	xp := reflect.ValueOf(&x).Elem()
	t.Log(xp.CanAddr(), xp.CanSet())
	xd := xp.Addr().Interface().(*int)
	*xd = 3
	t.Log(x)

	xp.Set(reflect.ValueOf(4))
	t.Log(x)

	xp.SetInt(5)
	t.Log(x)

	e := []int{1, 2, 3}
	ev := reflect.ValueOf(&e).Elem().Index(2)
	//ev := reflect.ValueOf(e).Index(2) //is ok
	t.Log(ev.CanAddr(), ev.CanSet()) //true true
	ev.SetInt(55)
	t.Log(e) //[1 2 55]
	ep := reflect.ValueOf(&e).Elem()
	ep.Set(reflect.ValueOf([]int{2, 3, 4}))
	ep.Set(reflect.Append(ep, reflect.ValueOf(5)))
	e2 := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(1)), 0, 2)
	ep.Set(reflect.AppendSlice(e2, reflect.ValueOf([]int{7, 8})))
	t.Log(e) // [2 3 4 5 7 8]

	// xp.Set(reflect.ValueOf(int64(4)))//panic: value of type int64 is not assignable to type int

	// x = 2
	// b := reflect.ValueOf(x)
	// b.Set(reflect.ValueOf(3)) // panic: Set using unaddressable value

	//SetInt 兼容interger, 而 Set 要求类型一致
	type I int
	var v I
	t.Log(reflect.TypeOf(v)) //test.I
	sv := reflect.ValueOf(&v).Elem()
	// sv.Set(reflect.ValueOf(3))        // panic:value of type int is not assignable to type test.I
	sv.SetInt(2)                //ok
	t.Log(reflect.TypeOf(v), v) //test.I, 2

	var y interface{}
	t.Log(reflect.TypeOf(y), reflect.ValueOf(y).Kind() == reflect.Invalid) //<nil> true
	ry := reflect.ValueOf(&y).Elem()
	// ry.SetInt(2)                     // panic: SetInt called on interface Value
	ry.Set(reflect.ValueOf(3))  // OK, y = int(3)
	t.Log(reflect.TypeOf(y), y) //int 3
	// ry.SetString("hello")            // panic: SetString called on interface Value
	ry.Set(reflect.ValueOf("hello")) // OK, y = "hello"
	t.Log(reflect.TypeOf(y), y)      //string hello

	//传入接口指针，保留接口类型信息
	var w io.Writer
	t.Log(reflect.TypeOf(w), reflect.ValueOf(w).Kind() == reflect.Invalid) //<nil> true
	// w = os.Stdin
	t.Log(reflect.TypeOf(w), reflect.ValueOf(w).Kind() == reflect.Ptr)                      //*os.File true
	t.Log(reflect.TypeOf(&w), reflect.ValueOf(&w).Elem().Type(), reflect.TypeOf(&w).Elem()) //*io.Writer io.Writer io.Writer
	rw := reflect.ValueOf(&w).Elem()
	var ws W = 33
	rw.Set(reflect.ValueOf(ws))
	t.Log(reflect.TypeOf(w), w) //test.W 33

	//unexported value can not be changed
	stdout := reflect.ValueOf(os.Stdout).Elem() // *os.Stdout, an os.File var
	t.Log(stdout.Type())                        // "os.File"
	t.Log(stdout, stdout.CanAddr(), stdout.CanSet(), stdout.CanInterface())
	stdout.Set(reflect.ValueOf(*os.Stderr))
	// stdout.Set(reflect.ValueOf(*os.Stdin)) //设置为*os.Stdin test无法输出???
	fd := stdout.FieldByName("fd")
	t.Log(fd.CanAddr(), fd.CanSet(), fd.CanInterface()) //true false false
	t.Log(fd.Kind() == reflect.Uintptr, fd.Uint())      //windows:ture, XXX
	// t.Log(fd.Kind() == reflect.Int, fd.Int()) //linux:true, 1
	// fd.SetInt(2) // panic: using value obtained using unexported field
	// t.Log(fd.Interface()) //cannot return value obtained from unexported field or method

	//同一个包的unexpected field 也不能Set
	type St struct {
		Export   int
		unexport int `tag:"test" display:"kk"`
	}
	vt := St{33, 44}
	// st := reflect.ValueOf(&vt).Elem()
	// unEx := st.FieldByName("unexport")
	// Ex := st.FieldByName("Export")
	// t.Log(st.CanAddr(), st.CanSet(), st.CanInterface())       //true true true
	// t.Log(unEx.CanAddr(), unEx.CanSet(), unEx.CanInterface()) //true false false
	// t.Log(Ex.CanAddr(), Ex.CanSet(), Ex.CanInterface())       //true true true
	// unEx.SetInt(44)//panic: reflect: reflect.Value.SetInt using value obtained using unexported field
	// Ex.SetInt(55)
	// t.Log(vt)
	// st.Set(reflect.ValueOf(St{44, 33}))
	// t.Log(vt)

	fieldInfo, _ := reflect.TypeOf(vt).FieldByName("unexport")
	display.Display("f", fieldInfo)
	t.Log(fieldInfo.Type) //int

	mp := make(map[int]bool)
	mp[33] = true
	vmp := reflect.ValueOf(&mp).Elem()
	kmp := vmp.MapIndex(vmp.MapKeys()[0])
	t.Log(vmp.CanAddr(), vmp.CanSet(), vmp.CanInterface())
	t.Log(kmp.CanAddr(), kmp.CanSet(), kmp.CanInterface())

	// t.Log(reflect.ValueOf(nil).Interface())//panic
	var i io.Writer = os.Stdout
	var c io.Closer
	id := reflect.ValueOf(&i).Elem()                                //id is io.Writer, dynamic type is *os.File
	concrete := id.Elem()                                           // *os.File
	assertType := reflect.TypeOf(&c).Elem()                         // io.Closer type
	cd := concrete.Convert(assertType)                              //io.Closer,  dynamic type is *os.File
	t.Log(cd.Elem() == concrete)                                    //true
	t.Log(reflect.DeepEqual(id.Interface(), id.Elem().Interface())) //true
	t.Log(reflect.ValueOf(nil), reflect.TypeOf(nil))

	nvp := reflect.Value{}
	t.Log(nvp, reflect.TypeOf(nvp), reflect.ValueOf(reflect.ValueOf(nil)))
	t.Log(reflect.TypeOf(reflect.TypeOf(1)))
	t.Log(reflect.TypeOf((*io.Writer)(&i)).Elem())

}
