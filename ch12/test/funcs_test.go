package test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"

	"gopl.io/ch12/display"
)

func TestFuncs(t *testing.T) {
	s := []string{"a", "b", "c"}
	sw := reflect.Swapper(s)
	// var f func(int, int)
	f := sw
	TimeFunc(&f, sw)
	f(1, 2)
	t.Log(s)
	// t.Log(Map(s, func(s string) string { return strings.ToUpper(s) }))
	// t.Log(Map("abc", func(r byte) byte { return byte(unicode.ToUpper(rune(r))) }))

	fc := sw
	pc := runtime.FuncForPC(reflect.ValueOf(fc).Pointer())
	// pc = runtime.FuncForPC(uintptr(unsafe.Pointer(&display.Display)))
	display.Display("fc", pc)
	t.Log(pc.Name(), pc.Entry())
	t.Log(pc.FileLine(reflect.ValueOf(fc).Pointer() + 3))
}

//泛型函数修饰器decorator

//TimeFunc decorator inputf function fs, make a function and store in ft
func TimeFunc(decoPtr, f interface{}) {
	target := reflect.ValueOf(f)
	decorator := reflect.ValueOf(decoPtr).Elem()
	fn := func(in []reflect.Value) (out []reflect.Value) {
		defer func(now time.Time) {
			fmt.Printf("---%s elapse %v ---\n", getFuncName(f), time.Since(now))
		}(time.Now())
		return target.Call(in)
	}

	decoFunc := reflect.MakeFunc(target.Type(), fn)
	decorator.Set(decoFunc)
}

func getFuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

//泛型交换函数

//泛型map函数
// s is string, slice, array; f is func(v T) T, T is elemnt type of s
func Map(s, f interface{}) interface{} {
	sv := reflect.ValueOf(s)
	fv := reflect.ValueOf(f)
	size := sv.Len()
	var etyp reflect.Type
	if sv.Kind() == reflect.String {
		etyp = reflect.TypeOf(byte(0))
	} else {
		etyp = sv.Type().Elem()
	}
	typ := reflect.ArrayOf(size, etyp)
	buf := reflect.New(typ).Elem()
	for i := 0; i < size; i++ {
		buf.Index(i).Set(fv.Call([]reflect.Value{sv.Index(i)})[0])
	}
	switch sv.Kind() {
	case reflect.String:
		// buf := reflect.MakeSlice(reflect.TypeOf(byte(0)), 0, size)
		return buf.Slice(0, size).Convert(sv.Type()).Interface()
	case reflect.Slice:
		return buf.Slice(0, size).Interface()
	case reflect.Array:
		return buf.Interface()
	}
	return nil
}

//泛型

//泛型与接口????
