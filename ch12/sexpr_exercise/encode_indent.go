// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 339.

package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
)

//这种实现结构上没有做到模块化，Marshal and MarshalIndent需要两份代码，需改进

//MarshalIndent encodes a Go value in indented S-expression form.
//optional: ignore zero value of field？ display unexported field?。
//如果显示了unexported field，而Unmarshal 不能处理。但是即使不显示，如果传入Unmarshal的结构含有unexported filed, 将维持原有的值（通常是zerovalue)，这样v != Unmarshal(Marshal(v))
func MarshalIndent(v interface{}) ([]byte, error) {
	var buf = new(bytes.Buffer)
	const pre = ' '
	var nPre = 0
	onlyExported := true //if true, display unexported field
	opt := false         ////if true, ignore zero value of field
	indent := func() {
		for i := 0; i < nPre; i++ {
			buf.WriteRune(pre)
		}
	}
	var encodeIndent func(v reflect.Value) error
	encodeIndent = func(v reflect.Value) error {
		switch v.Kind() {
		case reflect.Invalid:
			buf.WriteString("nil")

		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			fmt.Fprintf(buf, "%d", v.Int())

		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			fmt.Fprintf(buf, "%d", v.Uint())

		case reflect.Float32, reflect.Float64:
			fmt.Fprintf(buf, "%f", v.Float()) //如果%g, 0.0显示为0，Unmarshal 成了int,got error

		case reflect.Complex128, reflect.Complex64:
			fmt.Fprintf(buf, "#C(%f %f)", real(v.Complex()), imag(v.Complex()))

		case reflect.Bool:
			if v.Bool() {
				fmt.Fprint(buf, "t")
			} else {
				fmt.Fprint(buf, "nil")
			}

		case reflect.Interface:
			if v.IsNil() {
				fmt.Fprint(buf, "nil")
			} else {
				n, _ := fmt.Fprintf(buf, "(%s ", v.Elem().Type())
				nPre += n
				if err := encodeIndent(v.Elem()); err != nil {
					return err
				}
				fmt.Fprint(buf, ")")
				nPre -= n
			}

		case reflect.String:
			fmt.Fprintf(buf, "%q", v.String())

		case reflect.Ptr:
			return encodeIndent(v.Elem())

		case reflect.Array, reflect.Slice: // (value ...)
			buf.WriteByte('(')
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					buf.WriteByte('\n')
					indent()
				}
				if err := encodeIndent(v.Index(i)); err != nil {
					return err
				}
			}
			buf.WriteByte(')')

		case reflect.Struct: // ((name value) ...)
			buf.WriteByte('(')
			nPre++
			for i := 0; i < v.NumField(); i++ {
				if opt && isZero(v.Field(i)) || (onlyExported && v.Type().Field(i).PkgPath != "") {
					continue
				}
				if i > 0 {
					buf.WriteByte('\n')
					indent()
				}
				n, _ := fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
				nPre += n
				if err := encodeIndent(v.Field(i)); err != nil {
					return err
				}
				buf.WriteByte(')')
				nPre -= n
			}
			buf.WriteByte(')')
			nPre--

		case reflect.Map: // ((key value) ...)
			buf.WriteByte('(')
			nPre++
			for i, key := range v.MapKeys() {
				if opt && isZero(v.MapIndex(key)) {
					continue
				}
				if i > 0 {
					buf.WriteByte('\n')
					indent()
				}
				buf.WriteByte('(')
				if err := encodeIndent(key); err != nil {
					return err
				}
				buf.WriteByte(' ')
				if err := encodeIndent(v.MapIndex(key)); err != nil {
					return err
				}
				buf.WriteByte(')')
			}
			buf.WriteByte(')')
			nPre--
		default: // chan, func
			return fmt.Errorf("unsupported type: %s", v.Type())
		}
		return nil
	}

	if err := encodeIndent(reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0.0
	case reflect.Bool:
		return v.Bool() == false
	case reflect.String:
		return v.String() == ""
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i)) {
				return false
			}
		}
		return true //[0]int{} is zero value
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true //struct{}{} is zero value
	default: //func, map, slice, chan
		return v.IsNil()
	}
}

//only safe in sigle goroutine，多协程版本需要将nPre 传入pretty,所以干脆实现成lambda, 省掉了函数参数传递，但是同时结构上太紧凑，没有模块化
var nPre int
var pre = ' '

func marshalIndent(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	nPre = 0
	err := pretty(buf, reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func indent(buf *bytes.Buffer, nPre int) {
	for i := 0; i < nPre; i++ {
		buf.WriteRune(pre)
	}
}

func pretty(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", v.Float())

	case reflect.Complex128, reflect.Complex64:
		fmt.Fprintf(buf, "#C(%g %g)", real(v.Complex()), imag(v.Complex()))

	case reflect.Bool:
		if v.Bool() {
			fmt.Fprint(buf, "t")
		} else {
			fmt.Fprint(buf, "nil")
		}

	case reflect.Interface:
		if v.IsNil() {
			fmt.Fprint(buf, "nil")
		} else {
			n, _ := fmt.Fprintf(buf, "(%s ", v.Elem().Type())
			nPre += n
			if err := pretty(buf, v.Elem()); err != nil {
				return err
			}
			fmt.Fprint(buf, ")")
			nPre -= n
		}

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return pretty(buf, v.Elem())

	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := pretty(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')

	case reflect.Struct: // ((name value) ...)
		buf.WriteByte('(')
		nPre++
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				indent(buf, nPre)
			}
			n, _ := fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			nPre += n
			switch v.Field(i).Kind() {
			case reflect.Map, reflect.Struct:
			}
			if err := pretty(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
			buf.WriteByte('\n')
			nPre -= n
		}
		buf.Truncate(buf.Len() - 1) //delete lasat '\n'
		buf.WriteByte(')')
		nPre--

	case reflect.Map: // ((key value) ...)
		buf.WriteByte('(')
		nPre++
		for i, key := range v.MapKeys() {
			if i > 0 {
				indent(buf, nPre)
			}
			buf.WriteByte('(')
			if err := pretty(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := pretty(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
			buf.WriteByte('\n')
		}
		buf.Truncate(buf.Len() - 1) //delete lasat '\n'
		buf.WriteByte(')')
		nPre--
	default: // chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
