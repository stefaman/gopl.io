// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 357.

// Package unsafeptr demonstrates basic use of unsafe.Pointer.
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	//!+main
	var x struct {
		a bool
		b int16
		c []int
	}

	// equivalent to pb := &x.b
	// pb := (*int16)(unsafe.Pointer(
	// 	uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	// *pb = 42
	//!+wrong
	//tmp只是一个number,x没有被引用了或者内存位置（地址）因为GC发生了变动，pb有可能指向一个没有意义错误的旧地址
	// NOTE: subtly incorrect!
	tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	pb := (*int16)(unsafe.Pointer(tmp))
	*pb = 43
	//!-wrong

	fmt.Println(x.b) // "42"
	//!-main
	/*
	 */
}
