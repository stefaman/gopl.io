// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 344.

// Package sexpr provides a means for converting Go objects to and
// from S-expressions.
package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

//!+Unmarshal

// Unmarshal parses S-expression data and populates the variable
// whose address is in the non-nil pointer out.
func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	lex.next() // get the first token
	defer func() {
		// NOTE: this is not an example of ideal error handling.
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}

//!-Unmarshal

//!+lexer
type lexer struct {
	scan  scanner.Scanner
	token rune // the current token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
	if lex.token != want { // NOTE: Not an example of good error handling.
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

//!-lexer

// The read function is a decoder for a small subset of well-formed
// S-expressions.  For brevity of our example, it takes many dubious
// shortcuts.
//
// The parser assumes
// - that the S-expression input is well-formed; it does no error checking.
// - that the S-expression input corresponds to the type of the variable.
// - that all numbers in the input are non-negative decimal integers.
// - that all keys in ((key value) ...) struct syntax are unquoted symbols.
// - that the input does not contain dotted lists such as (1 2 . 3).
// - that the input does not contain Lisp reader macros such 'x and #'x.
//
// The reflection logic assumes
// - that v is always a variable of the appropriate type for the
//   S-expression value.  For example, v must not be a boolean,
//   interface, channel, or function, and if v is an array, the input
//   must have the correct number of elements.
// - that v in the top-level call to read has the zero value of its
//   type and doesn't need clearing.
// - that if v is a numeric variable, it is a signed integer.

//!+read
func read(lex *lexer, v reflect.Value) {
	//!+stef
	//支持指针类型
	if v.Kind() == reflect.Ptr { //Marshal 是显示*p
		if lex.text() == "nil" { //避免(*T)(nil)指针处理成*p = T(zero vale)
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
		// fmt.Println(v, v.CanSet(), v.Elem().CanSet()) //<nil> true false
		value := reflect.New(v.Type().Elem()) //value type is the same as v
		// fmt.Println(value, value.CanSet(), value.Elem().CanSet())//0xc0420500c8 false true
		read(lex, value.Elem())
		v.Set(value)
		return
	} //!-stef
	switch lex.token {
	case scanner.Ident:
		// The only valid identifiers are
		// "nil" and struct field names.
		if lex.text() == "nil" { //include "false"
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
		if lex.text() == "t" { //t is "true", nil is "false"
			v.SetBool(true)
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text()) // NOTE: ignoring errors
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text()) // NOTE: ignoring errors
		v.SetInt(int64(i))
		lex.next()
		return
	case scanner.Float:
		f, err := strconv.ParseFloat(lex.text(), 64)
		if err != nil {
			panic(fmt.Sprintf("can't parse %q as float number", lex.text()))
		}
		v.SetFloat(f)
		lex.next()
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.next() // consume ')'
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

//!-read

//!+readlist
func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array: // (item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}

	case reflect.Slice: // (item ...)
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}

	case reflect.Struct: // ((name value) ...)
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}

	case reflect.Map: // ((key value) ...)
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}
	case reflect.Interface: //传入Unmarshal的容器参数，接口类型field是nil，解析时需要获得具体类型的信息；
		// 要么是查表，但这种实现，表的容量是有限的
		// 或者要么规定传入容器参数的接口类型field必须包括想要解析的具体类型
		tName := ""
		for lex.token == scanner.Ident || strings.ContainsRune(".*", lex.token) {
			tName += lex.text()
			lex.next()
		}
		//查表的实现
		// typ, err := gotType(tName)
		// if err != nil {
		// 	panic(fmt.Sprintf("cannot decode %s type of interface %s: %v", tName, v.Type(), err))
		// }
		// value := reflect.New(typ).Elem()
		// read(lex, value)
		// v.Set(value)

		//传入类型的实现
		// fmt.Println(v.Elem(), v.Elem().CanSet()) // xx, false
		value := reflect.New(v.Elem().Type()).Elem() //if v is nil, then painic: call of reflect.Value.Type on zero Value
		read(lex, value)
		v.Set(value)

	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

func gotType(typ string) (reflect.Type, error) {
	return reflect.TypeOf(I(33)), nil
}

type I int

func (i I) Write([]byte) (int, error) {
	return int(i), nil
}

//!-readlist
