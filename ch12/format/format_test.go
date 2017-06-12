// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package format_test

import (
	"fmt"
	"testing"
	"time"
	"strings"
	"io"
	"os"
	"reflect"
	"math/rand"
	"gopl.io/ch12/format"
)

func Test(t *testing.T) {
	// The pointer values are just examples, and may vary from run to run.
	//!+time
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(format.Any(x))                  // "1"
	fmt.Println(format.Any(d))                  // "1"
	fmt.Println(format.Any([]int64{x}))         // "[]int64 0x8202b87b0"
	fmt.Println(format.Any([]time.Duration{d})) // "[]time.Duration 0x8202b87e0"
	//!-time
}

var(
	w io.Writer
	ptr *int
	st struct{
		v interface{};
		want string;
	}
	f func(string)

)
//stefaman
var tests = []struct{
	v interface{};
	want string;
}{
	{1, "1"},
	{1.2, "float64 value"},
	{[]byte("haha"), "[]uint8"},
	{[]rune{'蔡'}, "[]int32"},
	{"kk", `"kk"`},
	{[3]int{1,2,3}, "[3]int value"},
	{map[string]bool{"cai":true},""},
	{strings.Contains,".."},
	{struct{v interface{}; want string;}{1,""}, ""},
	{nil, "invalid"},
	{new(int), ""},
	{io.Writer(os.Stdout), ""},
	{rand.NewSource(1), "rngSource"},
	{reflect.ValueOf(1), ""},
	{reflect.TypeOf(1), ""},
	//zero value
	{w, "invalid"},
	{ptr, "*int"},
	{st, "struct { v interface {}; want string }"},
	{f, "func(string)"},
	// {},
	// {},



}
//add more tests
func TestTable(t *testing.T) {
	for _, test := range tests {
		s := format.Any(test.v)
		if !strings.Contains(s, test.want) || len(test.want) == 0 {
			t.Errorf("for format.Any(%v) get %q, want %q", test.v, s, test.want)
		}
	}
}
