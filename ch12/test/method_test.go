package test

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"gopl.io/ch12/pack"
)

//
func TestMethod(t *testing.T) {

	// m's type is reflect.Method struct
	m, _ := reflect.TypeOf(time.Hour).MethodByName("Hours")
	// display.DisplayE("m", m)
	// t.Log(m.Type.String())        // func Hours(time.Duration) float64
	// t.Log(m.Func.Type().String()) // 同上
	//m.Func is a reflect.Value, is a method expression, should pass reciever as first parameter
	retM := m.Func.Call([]reflect.Value{reflect.ValueOf(time.Hour)})[0].Interface().(float64)
	t.Log(retM) //1

	//不能用于nil interface value
	//f relfect.Value, is a method value, a closure that has a reciver
	f := reflect.ValueOf(time.Hour).MethodByName("Hours")
	// t.Log(f.Type().String()) // func() float64
	retF := f.Call(nil)[0].Interface().(float64)
	t.Log(retF) //1

	// 接口类型方法获取
	// var w io.Writer
	// ftyp := reflect.TypeOf(&w).Elem()
	// fm, _ := ftyp.MethodByName("Write")
	// t.Log(ftyp, fm.Type, fm.Func) //io.Writer func([]uint8) (int, error) <invalid reflect.Value>

	PrintInterface((*sort.Interface)(nil))
	//Output
	// type sort.Interface
	// func() int
	// func(int, int) bool
	// func(int, int)

	// w = os.Stdout
	// tv := reflect.ValueOf(&w).Elem().Elem()
	// fmt.Printf("tve %s\n", tv.Type())
	// for i := 0; i < tv.NumMethod(); i++ {
	// 	fmt.Printf("%s\n", tv.Method(i).Type())
	// }
	// PrintMethods(w)

	type st struct {
		W *io.Writer
		// os.File
		// *pack.Test
		pack.Test
	}
	// s := st{&pack.Test{}}
	s := st{}
	t.Log(reflect.TypeOf(s).MethodByName("Efunc"))
	PrintMethods(reflect.ValueOf(&s).Elem().Addr().Interface())
	PrintMethods(&s)
	// s.Ufunc()
	// s.Tfunc()
	// (&s).Tfunc()
	s.Ufunc()
	(&s).Ufunc()
	// s.Close()
	// (&s).Close()
	// s.Write([]byte{}) //ambiguous selector s.Write
	// ityp := (reflect.TypeOf((*io.Writer)(nil)).Elem())
	// t.Log(reflect.TypeOf(s).Implements(ityp))
	// t.Log(reflect.TypeOf(s).Field(0).Type.Implements(ityp))
	// t.Log(reflect.TypeOf(s).MethodByName("Write"))         // {<nil> <invalid Value> 0} false
	// t.Log(reflect.ValueOf(s).MethodByName("Close").Type()) //func() error

}

func PrintMethods(i interface{}) { //i is not nil
	// v := reflect.Indirect(reflect.ValueOf(i))
	v := reflect.ValueOf(i)
	typ := v.Type()
	fmt.Printf("value type is %s\n", typ)
	for i := 0; i < typ.NumMethod(); i++ {
		fmt.Printf("func (%s) %s%s\n", typ, typ.Method(i).Name, strings.TrimPrefix(v.Method(i).Type().String(), "func"))
	}
}

func PrintInterface(i interface{}) { //i is *interface
	typ := reflect.TypeOf(i).Elem()
	fmt.Printf("interface type %s\n", typ.String())
	for i := 0; i < typ.NumMethod(); i++ {
		fmt.Printf("func %s%s\n", typ.Method(i).Name, strings.TrimPrefix(typ.Method(i).Type.String(), "func"))
	}
}
