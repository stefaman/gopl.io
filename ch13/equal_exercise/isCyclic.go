package equal

import (
	"reflect"
	"unsafe"
)

//IsCycli check variable if cycle, can not deal with cycle map
func IsCyclic(x interface{}) bool {
	seen := make(map[see]bool)
	return cyclic(reflect.ValueOf(x), seen)
}

type see struct {
	ptr unsafe.Pointer
	typ reflect.Type
}

func cyclic(v reflect.Value, seen map[see]bool) bool {
	if v.Kind() == reflect.Invalid {
		return false
	}
	if v.CanAddr() {
		s := see{ptr: unsafe.Pointer(v.UnsafeAddr()), typ: v.Type()}
		if seen[s] {
			return true
		}
		seen[s] = true
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if cyclic(v.Index(i), seen) {
				return true
			}
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if cyclic(v.Field(i), seen) {
				return true
			}
		}
	case reflect.Map: //can not deal cycle map
		for _, k := range v.MapKeys() {
			if cyclic(v.MapIndex(k), seen) {
				return true
			}
		}
	case reflect.Ptr, reflect.Interface:
		return cyclic(v.Elem(), seen)
	}
	return false
}
