// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 359.

// Package equal provides a deep equivalence relation for arbitrary values.
package equal

import (
	"reflect"
	"unsafe"
)

//!+
func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if x.Type() != y.Type() {
		return false
	}

	// ...cycle check omitted (shown later)...

	//!-
	//!+cyclecheck
	// cycle check
	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true // identical references
		}
		c := comparison{xptr, yptr, x.Type()}
		if seen[c] {
			return true // already seen
		}
		seen[c] = true
	}
	//!-cyclecheck
	//!+
	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()

	case reflect.String:
		return x.String() == y.String()

	// ...numeric cases omitted for brevity...

	//!-
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return x.Int() == y.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return x.Uint() == y.Uint()

	case reflect.Float32, reflect.Float64:
		// return x.Float() == y.Float()
		return equalInErr(x.Float(), y.Float())

	case reflect.Complex64, reflect.Complex128:
		// return x.Complex() == y.Complex()
		return equalInErr(real(x.Complex()), real(y.Complex())) && equalInErr(imag(x.Complex()), imag(y.Complex()))
	//!+
	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return x.Pointer() == y.Pointer()

	case reflect.Ptr, reflect.Interface:
		return equal(x.Elem(), y.Elem(), seen)

	case reflect.Array, reflect.Slice:
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		return true

	// ...struct and map cases omitted for brevity...
	//!-
	case reflect.Struct:
		for i, n := 0, x.NumField(); i < n; i++ {
			if !equal(x.Field(i), y.Field(i), seen) {
				return false
			}
		}
		return true

	case reflect.Map:
		if x.Len() != y.Len() {
			return false
		}
		for _, k := range x.MapKeys() {

			// xd := x.MapIndex(k).Interface()
			// yd := y.MapIndex(k).Interface()
			// see := comparison{unsafe.Pointer(reflect.ValueOf(xd).InterfaceData()[1]),
			// 	unsafe.Pointer(reflect.ValueOf(yd).InterfaceData()[1]),
			// 	x.MapIndex(k).Type(),
			// }
			// if seen[see] {
			// 	return true
			// }
			// seen[see] = true
			// for x, deep := x.MapIndex(k), 0; x != nil && x.Kind() == reflect.Map; x = x.MapIndex() {
			// 	deep++
			// 	if deep > 10 {
			// 		return true
			// 	}
			// }
			//x.MapIndex is unaddressable, infinite loop for cycle map
			if !equal(x.MapIndex(k), y.MapIndex(k), seen) {
				return false
			}
		}
		return true
		//!+
	}
	panic("unreachable")
}

//!-

//!+comparison
// Equal reports whether x and y are deeply equal.
//!-comparison
//
// Map keys are always compared with ==, not deeply.
// (This matters for keys containing pointers or interfaces.)
//!+comparison
//can not deal with cycle map
func Equal(x, y interface{}) bool {
	seen := make(map[comparison]bool)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

//不同类型可以同一地址，如x and x[0]
//不用uintptr, 不安全
type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}

//!-comparison

//e13.1
func equalInErr(x, y float64) bool {
	min := x
	if x > y {
		min = y
	}
	return (y-x)*1e9 < min && (x-y)*1e9 < min
}
