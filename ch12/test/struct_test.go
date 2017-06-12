package test

import (
	"os"
	"reflect"
	"testing"
)

func TestSturct(t *testing.T) {
	type Em1 struct {
		Emb int
		One int
	}
	type Em2 struct {
		Emb  int
		Two  int
		file int
	}
	type St struct {
		Em1
		Em2
		// pack.Test
		os.File
		// file     int
		Export   int
		unexport int `tag:"test" display:"kk"`
	}
	var s = reflect.ValueOf(St{})
	stp := s.Type()
	for i := 0; i < stp.NumField(); i++ {

		t.Log(stp.Field(i))
	}
	t.Log(stp.FieldByName("file"))

	for i := 0; i < s.NumField(); i++ {

		t.Log(s.Field(i))
	}
	t.Log(s.FieldByName("b"))

}
