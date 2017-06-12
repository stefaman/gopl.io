// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 333.

// Package display provides a means to displayE structured data.
package display

import (
	"fmt"
	"reflect"
	"strconv"
)

//Limit limits print depth
const Limit = 30

var depth = 0

//DisplayE display any variable, including circle infinite type
//e12.1/2
func DisplayE(name string, x interface{}) {
	depth = 0
	fmt.Printf("Display %s (%T):\n", name, x)
	displayE(name, reflect.ValueOf(x))
}
func fmtKey(key reflect.Value) string {
	switch key.Kind() {
	case reflect.Array:
		return fmt.Sprintf("%v", key)
	case reflect.Struct:
		return fmt.Sprintf("%v", key)
	default:
		return formatAtom(key)
	}
}

//Display 递归的显示内部，因为参数的名字会丢失，所以需要单独传入变量的identifier
//另一方面，单独传入name也是递归的path
func displayE(path string, v reflect.Value) {
	if depth > Limit { //avoid infinite circle reference
		fmt.Println(path, ": depth excced limit", Limit)
		depth = 0
		return
	}
	depth++
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		}
		fallthrough
	case reflect.Array: //if nil or empty
		if v.Len() == 0 {
			fmt.Printf("%s = []\n", path)
		}
		for i := 0; i < v.Len(); i++ {
			displayE(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		if v.NumField() == 0 { //显示empty struct
			fmt.Printf("%s = {}\n", path)
			return
		}
		for i := 0; i < v.NumField(); i++ { //内嵌成员没有区别对待
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			displayE(fieldPath, v.Field(i))
		}
	case reflect.Map:
		if v.IsNil() { // if nil
			fmt.Printf("%s = nil\n", path)
			return
		}
		if v.Len() == 0 { //if empty
			fmt.Printf("%s = [][]\n", path)
			return
		}
		for _, key := range v.MapKeys() {
			displayE(fmt.Sprintf("%s[%s]", path,
				fmtKey(key)), v.MapIndex(key))
		}

	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			displayE(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			displayE(path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

//!-displayE

// formatAtom formats a value without inspecting its internal structure.
// It is a copy of the the function in gopl.io/ch11/format.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr,
		reflect.Slice, reflect.Map, reflect.UnsafePointer:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}
